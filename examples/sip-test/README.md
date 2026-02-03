# SIP Test Client

A local SIP test client for testing the Rapida Voice AI SIP integration before connecting to production providers like Twilio, Vonage, or Exotel.

## Prerequisites

1. **Start the services:**
   ```bash
   # From project root
   make up-all
   ```

2. **Ensure assistant-api is configured for SIP:**
   - SIP server should be listening on port 5060 (or configured port)
   - Assistant deployment should have SIP enabled

## Usage

### Basic Test Call

```bash
cd examples/sip-test
go mod tidy
go run main.go
```

### Custom Configuration

```bash
go run main.go \
  -local-ip 127.0.0.1 \
  -local-port 5061 \
  -remote-host localhost \
  -remote-port 5060 \
  -caller "sip:testcaller@localhost" \
  -callee "sip:assistant@localhost" \
  -transport udp \
  -duration 30s
```

### Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-local-ip` | `127.0.0.1` | Local IP address |
| `-local-port` | `5061` | Local SIP port |
| `-remote-host` | `localhost` | Remote SIP server host (assistant-api) |
| `-remote-port` | `5060` | Remote SIP server port |
| `-caller` | `sip:testcaller@localhost` | Caller SIP URI |
| `-callee` | `sip:assistant@localhost` | Callee SIP URI |
| `-transport` | `udp` | Transport protocol (udp/tcp/tls) |
| `-duration` | `30s` | Call duration for test |
| `-username` | `testuser` | SIP username |
| `-password` | `testpass` | SIP password |

## Testing Scenarios

### 1. Basic Connectivity Test
```bash
go run main.go -duration 5s
```
Verifies that the SIP server accepts INVITE and responds with 100/180/200.

### 2. Extended Call Test
```bash
go run main.go -duration 60s
```
Tests a longer call to verify session stability.

### 3. Different Transport Test
```bash
go run main.go -transport tcp
go run main.go -transport tls
```

## Alternative Testing Methods

### Using SIPp (Load Testing)
```bash
# Install SIPp
brew install sipp

# Basic INVITE test
sipp -sn uac localhost:5060 -m 1 -d 5000

# Load test with 10 concurrent calls
sipp -sn uac localhost:5060 -m 100 -r 10 -d 5000
```

### Using Linphone (GUI Softphone)
```bash
brew install --cask linphone
```
Configure Linphone with:
- Account: testuser@localhost:5060
- Transport: UDP

### Using pjsua (CLI Softphone)
```bash
# Install PJSIP tools
brew install pjsip

# Make a test call
pjsua --local-port=5061 sip:assistant@localhost:5060
```

## Debugging

### View SIP Traffic
```bash
# Using tcpdump
sudo tcpdump -i lo0 -n port 5060 -A

# Using tshark (Wireshark CLI)
tshark -i lo0 -f "port 5060" -V
```

### Check assistant-api logs
```bash
# From project root
make logs-assistant
```

## Expected Flow

1. **INVITE sent** → Client initiates call
2. **100 Trying** ← Server acknowledges
3. **180 Ringing** ← Server ringing
4. **200 OK** ← Call connected (with SDP answer)
5. **ACK sent** → Client confirms
6. *(Call active for configured duration)*
7. **BYE sent** → Client ends call
8. **200 OK** ← Server confirms hangup

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Connection refused | Ensure assistant-api is running with SIP enabled |
| 403 Forbidden | Check authentication credentials |
| 404 Not Found | Verify the assistant/callee URI is correct |
| 408 Timeout | Network issue or server not responding |
| 503 Service Unavailable | Server overloaded or not configured |

## Next Steps

After local testing succeeds:
1. Configure Twilio SIP Trunk to point to your public assistant-api
2. Set up proper TLS certificates for production
3. Configure authentication in vault credentials
4. Test with Twilio test numbers before going live
