// Copyright (c) 2023-2025 RapidaAI
// SIP Test Client - Tests SIP integration locally before connecting to production providers
// This simulates a SIP endpoint calling into the assistant-api

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
)

// Config holds test configuration
type Config struct {
	LocalIP      string
	LocalPort    int
	RemoteHost   string
	RemotePort   int
	CallerID     string
	CalleeID     string
	Transport    string
	CallDuration time.Duration
	Username     string
	Password     string
}

func main() {
	cfg := parseFlags()

	log.Printf("SIP Test Client Starting")
	log.Printf("Local:  %s:%d", cfg.LocalIP, cfg.LocalPort)
	log.Printf("Remote: %s:%d", cfg.RemoteHost, cfg.RemotePort)
	log.Printf("Caller: %s -> Callee: %s", cfg.CallerID, cfg.CalleeID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	if err := runSIPClient(ctx, cfg); err != nil {
		log.Fatalf("SIP client error: %v", err)
	}
}

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.LocalIP, "local-ip", "127.0.0.1", "Local IP address")
	flag.IntVar(&cfg.LocalPort, "local-port", 5061, "Local SIP port")
	flag.StringVar(&cfg.RemoteHost, "remote-host", "localhost", "Remote SIP server host")
	flag.IntVar(&cfg.RemotePort, "remote-port", 5060, "Remote SIP server port")
	flag.StringVar(&cfg.CallerID, "caller", "sip:testcaller@localhost", "Caller SIP URI")
	flag.StringVar(&cfg.CalleeID, "callee", "sip:assistant@localhost", "Callee SIP URI (assistant)")
	flag.StringVar(&cfg.Transport, "transport", "udp", "Transport protocol (udp/tcp/tls)")
	flag.DurationVar(&cfg.CallDuration, "duration", 30*time.Second, "Call duration for test")
	flag.StringVar(&cfg.Username, "username", "testuser", "SIP username for auth")
	flag.StringVar(&cfg.Password, "password", "testpass", "SIP password for auth")

	flag.Parse()
	return cfg
}

func runSIPClient(ctx context.Context, cfg *Config) error {
	// Create SIP User Agent
	ua, err := sipgo.NewUA(
		sipgo.WithUserAgent("RapidaSIPTestClient/1.0"),
	)
	if err != nil {
		return fmt.Errorf("failed to create SIP UA: %w", err)
	}

	// Create SIP client
	client, err := sipgo.NewClient(ua,
		sipgo.WithClientHostname(cfg.LocalIP),
		sipgo.WithClientPort(cfg.LocalPort),
	)
	if err != nil {
		return fmt.Errorf("failed to create SIP client: %w", err)
	}

	// Also create a server to receive responses
	server, err := sipgo.NewServer(ua)
	if err != nil {
		return fmt.Errorf("failed to create SIP server: %w", err)
	}

	// Register response handlers
	var callID string
	responseChan := make(chan *sip.Response, 10)

	server.OnBye(func(req *sip.Request, tx sip.ServerTransaction) {
		log.Println("Received BYE from remote")
		resp := sip.NewResponseFromRequest(req, 200, "OK", nil)
		tx.Respond(resp)
	})

	// Start listening for responses
	listenAddr := fmt.Sprintf("%s:%d", cfg.LocalIP, cfg.LocalPort)
	go func() {
		if err := server.ListenAndServe(ctx, cfg.Transport, listenAddr); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	log.Printf("SIP client listening on %s/%s", listenAddr, cfg.Transport)
	time.Sleep(500 * time.Millisecond) // Wait for server to start

	// Build INVITE request
	remoteAddr := fmt.Sprintf("%s:%d", cfg.RemoteHost, cfg.RemotePort)
	toURI := sip.Uri{
		User: "assistant",
		Host: cfg.RemoteHost,
		Port: cfg.RemotePort,
	}

	req := sip.NewRequest(sip.INVITE, toURI)
	req.SetDestination(remoteAddr)

	// Set From header
	fromURI := sip.Uri{
		User: "testcaller",
		Host: cfg.LocalIP,
		Port: cfg.LocalPort,
	}
	from := sip.FromHeader{
		Address: sip.Address{Uri: fromURI},
		Params:  sip.NewParams(),
	}
	from.Params.Add("tag", sip.GenerateTagN(8))
	req.AppendHeader(&from)

	// Set To header
	to := sip.ToHeader{
		Address: sip.Address{Uri: toURI},
	}
	req.AppendHeader(&to)

	// Set Contact header
	contact := sip.ContactHeader{
		Address: sip.Address{Uri: fromURI},
	}
	req.AppendHeader(&contact)

	// Set Call-ID
	callID = sip.GenerateTagN(16)
	callIDHeader := sip.CallIDHeader(callID)
	req.AppendHeader(&callIDHeader)

	// Set CSeq
	cseq := sip.CSeqHeader{
		SeqNo:      1,
		MethodName: sip.INVITE,
	}
	req.AppendHeader(&cseq)

	// Set Max-Forwards
	maxFwd := sip.MaxForwardsHeader(70)
	req.AppendHeader(&maxFwd)

	// Add SDP body (minimal audio offer)
	sdpBody := buildSDP(cfg.LocalIP, 10000)
	req.SetBody([]byte(sdpBody))

	contentType := sip.ContentTypeHeader("application/sdp")
	req.AppendHeader(&contentType)

	contentLength := sip.ContentLengthHeader(len(sdpBody))
	req.AppendHeader(&contentLength)

	log.Printf("Sending INVITE to %s", remoteAddr)
	log.Printf("Call-ID: %s", callID)

	// Send INVITE
	tx, err := client.TransactionRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send INVITE: %w", err)
	}

	// Wait for responses
	callConnected := false
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled")
			if callConnected {
				sendBye(client, ctx, cfg, callID)
			}
			return nil

		case resp := <-responseChan:
			handleResponse(resp)

		case resp, ok := <-tx.Responses():
			if !ok {
				log.Println("Transaction closed")
				return nil
			}

			statusCode := resp.StatusCode
			log.Printf("Received response: %d %s", statusCode, resp.Reason)

			switch {
			case statusCode == 100:
				log.Println("Call is trying...")
			case statusCode == 180:
				log.Println("Call is ringing...")
			case statusCode == 183:
				log.Println("Session progress...")
			case statusCode == 200:
				log.Println("Call connected! Sending ACK...")
				callConnected = true

				// Send ACK
				ack := sip.NewAckRequest(req, resp, nil)
				if err := client.WriteRequest(ack); err != nil {
					log.Printf("Failed to send ACK: %v", err)
				}

				// Simulate call duration
				log.Printf("Call active for %v...", cfg.CallDuration)
				time.Sleep(cfg.CallDuration)

				// Send BYE
				sendBye(client, ctx, cfg, callID)
				log.Println("Call ended normally")
				return nil

			case statusCode >= 400:
				log.Printf("Call failed: %d %s", statusCode, resp.Reason)
				return fmt.Errorf("call failed with status %d", statusCode)
			}
		}
	}
}

func handleResponse(resp *sip.Response) {
	log.Printf("Async response: %d %s", resp.StatusCode, resp.Reason)
}

func sendBye(client *sipgo.Client, ctx context.Context, cfg *Config, callID string) {
	remoteAddr := fmt.Sprintf("%s:%d", cfg.RemoteHost, cfg.RemotePort)
	toURI := sip.Uri{
		User: "assistant",
		Host: cfg.RemoteHost,
		Port: cfg.RemotePort,
	}

	bye := sip.NewRequest(sip.BYE, toURI)
	bye.SetDestination(remoteAddr)

	// Set headers
	fromURI := sip.Uri{
		User: "testcaller",
		Host: cfg.LocalIP,
		Port: cfg.LocalPort,
	}
	from := sip.FromHeader{
		Address: sip.Address{Uri: fromURI},
		Params:  sip.NewParams(),
	}
	from.Params.Add("tag", sip.GenerateTagN(8))
	bye.AppendHeader(&from)

	to := sip.ToHeader{
		Address: sip.Address{Uri: toURI},
	}
	bye.AppendHeader(&to)

	callIDHeader := sip.CallIDHeader(callID)
	bye.AppendHeader(&callIDHeader)

	cseq := sip.CSeqHeader{
		SeqNo:      2,
		MethodName: sip.BYE,
	}
	bye.AppendHeader(&cseq)

	maxFwd := sip.MaxForwardsHeader(70)
	bye.AppendHeader(&maxFwd)

	log.Println("Sending BYE...")
	if _, err := client.TransactionRequest(ctx, bye); err != nil {
		log.Printf("Failed to send BYE: %v", err)
	}
}

func buildSDP(localIP string, rtpPort int) string {
	return fmt.Sprintf(`v=0
o=- %d %d IN IP4 %s
s=RapidaSIPTest
c=IN IP4 %s
t=0 0
m=audio %d RTP/AVP 0 8 101
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000
a=rtpmap:101 telephone-event/8000
a=fmtp:101 0-16
a=sendrecv
a=ptime:20
`, time.Now().Unix(), time.Now().Unix(), localIP, localIP, rtpPort)
}
