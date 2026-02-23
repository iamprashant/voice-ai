import {
  AssistantPhoneDeployment,
  DeploymentAudioProvider,
} from '@rapidaai/react';
import { Tab } from '@/app/components/tab';
import { cn } from '@/utils';
import { ChevronRight } from 'lucide-react';
import { ModalProps } from '@/app/components/base/modal';
import { RightSideModal } from '@/app/components/base/modal/right-side-modal';
import { FieldSet } from '@/app/components/form/fieldset';
import { FormLabel } from '@/app/components/form-label';
import { CONFIG } from '@/configs';
import { CopyButton } from '@/app/components/form/button/copy-button';
import { InputHelper } from '@/app/components/input-helper';
import { YellowNoticeBlock } from '@/app/components/container/message/notice-block';
import { ProviderPill } from '@/app/components/pill/provider-model-pill';
import { FC, useMemo } from 'react';

interface AssistantPhoneCallDeploymentDialogProps extends ModalProps {
  deployment: AssistantPhoneDeployment;
}
/**
 *
 * @param props
 * @returns
 */
export function AssistantPhoneCallDeploymentDialog(
  props: AssistantPhoneCallDeploymentDialogProps,
) {
  const providerName = props.deployment?.getPhoneprovidername()?.toLowerCase();
  const assistantId = props.deployment?.getAssistantid();
  const mediaHost = CONFIG.connection.media;
  const sipHost = CONFIG.connection.sip;
  const socketHost = CONFIG.connection.socket;

  const webhookUrl = `${mediaHost}/v1/talk/${providerName}/call/${assistantId}?x-api-key={{PROJECT_CRDENTIAL_KEY}}`;
  const eventUrl = `${mediaHost}/v1/talk/${providerName}/event/${assistantId}?x-api-key={{PROJECT_CRDENTIAL_KEY}}`;

  return (
    <RightSideModal
      modalOpen={props.modalOpen}
      setModalOpen={props.setModalOpen}
      className="w-2/3 xl:w-1/3 flex-1"
    >
      <div className="flex items-center p-4 border-b text-sm/6 font-medium">
        <div className="font-medium">Assistant</div>
        <ChevronRight size={18} className="mx-2" />
        <div className="font-medium">Deployment</div>
        <ChevronRight size={18} className="mx-2" />
        <div className="font-medium">vrsn_dpl_{props.deployment.getId()}</div>
      </div>
      <div className="relative overflow-auto h-[calc(100vh-50px)] flex flex-col flex-1">
        <Tab
          active="Integration"
          className={cn(
            'text-sm',
            'bg-gray-50 border-b dark:bg-gray-900 dark:border-gray-800 sticky top-0 z-1',
          )}
          tabs={[
            {
              label: 'Integration',
              element: (
                <div className="flex-1 px-4 space-y-8">
                  {providerName === 'sip' && (
                    <SipIntegrationInstructions
                      sipHost={sipHost}
                      assistantId={assistantId}
                    />
                  )}

                  {providerName === 'asterisk' && (
                    <AsteriskIntegrationInstructions
                      mediaHost={mediaHost}
                      audioSocketHost={socketHost}
                      assistantId={assistantId}
                    />
                  )}

                  {providerName !== 'sip' && providerName !== 'asterisk' && (
                    <>
                      <FieldSet className="col-span-2">
                        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
                          Inbound webhook url
                        </div>
                        <div className="flex items-center gap-2">
                          <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
                            {webhookUrl}
                          </code>
                          <div className="flex shrink-0 border divide-x">
                            <CopyButton className="h-7 w-7">
                              {webhookUrl}
                            </CopyButton>
                          </div>
                        </div>
                        <InputHelper>
                          You can add all the additional agent arguments in
                          query parameters for example if you are expecting
                          argument
                          <code className="text-red-600">`name`</code>
                          add{' '}
                          <code className="text-red-600">
                            `?name=your-name`
                          </code>
                        </InputHelper>
                      </FieldSet>
                      <FieldSet className="col-span-2">
                        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
                          Call status changes / Event callback webhook
                        </div>
                        <div className="flex items-center gap-2">
                          <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
                            {eventUrl}
                          </code>
                          <div className="flex shrink-0 border divide-x">
                            <CopyButton className="h-7 w-7">
                              {eventUrl}
                            </CopyButton>
                          </div>
                        </div>
                      </FieldSet>
                    </>
                  )}
                </div>
              ),
            },
            {
              label: 'Audio',
              element: (
                <div className="flex-1 space-y-8">
                  <VoiceInput deployment={props.deployment?.getInputaudio()} />

                  <VoiceOutput
                    deployment={props.deployment?.getOutputaudio()}
                  />
                </div>
              ),
            },
          ]}
        />
      </div>
    </RightSideModal>
  );
}

/* -------------------------------------------------------------------------- */
/*  SIP Provider Integration Instructions                                     */
/* -------------------------------------------------------------------------- */

const SipIntegrationInstructions: FC<{
  sipHost?: string;
  assistantId: string;
}> = ({ sipHost, assistantId }) => {
  const sipEndpoint = `sip:${sipHost}`;
  const sipUri = `sip:${assistantId || '{ASSISTANT_ID}'}:{{PROJECT_CREDENTIAL_KEY}}@${sipHost}`;

  return (
    <>
      <FieldSet className="col-span-2">
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          SIP Server Endpoint
        </div>
        <div className="flex items-center gap-2">
          <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
            {sipEndpoint}
          </code>
          <div className="flex shrink-0 border divide-x">
            <CopyButton className="h-7 w-7">{sipEndpoint}</CopyButton>
          </div>
        </div>
        <InputHelper>
          Point your SIP trunk / PBX outbound proxy to this address. Rapida
          accepts SIP INVITE and establishes an RTP media session directly.
        </InputHelper>
      </FieldSet>

      <FieldSet className="col-span-2">
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          SIP URI (Authentication)
        </div>
        <div className="flex items-center gap-2">
          <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
            {sipUri}
          </code>
          <div className="flex shrink-0 border divide-x">
            <CopyButton className="h-7 w-7">{sipUri}</CopyButton>
          </div>
        </div>
        <InputHelper>
          Authentication is embedded in the SIP URI:{' '}
          <code className="text-red-600">
            sip:{'{'}
            assistantID{'}'}:{'{'}
            apiKey{'}'}@host
          </code>
          . Replace{' '}
          <code className="text-red-600">
            {'{{'}PROJECT_CREDENTIAL_KEY{'}}'}
          </code>{' '}
          with your project API key.
        </InputHelper>
      </FieldSet>

      <FieldSet className="col-span-2">
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          SIP Configuration Details
        </div>
        <div className="text-xs text-gray-500 dark:text-gray-400 space-y-3 pt-2">
          <div className="grid grid-cols-2 gap-2">
            <div className="font-medium text-gray-700 dark:text-gray-300">
              Transport
            </div>
            <div>UDP, TCP, or TLS</div>
            <div className="font-medium text-gray-700 dark:text-gray-300">
              Port
            </div>
            <div>5060</div>
            <div className="font-medium text-gray-700 dark:text-gray-300">
              Codec
            </div>
            <div>
              G.711 μ-law (PCMU), G.711 A-law (PCMA) + telephone-event (DTMF)
            </div>
            <div className="font-medium text-gray-700 dark:text-gray-300">
              Authentication
            </div>
            <div>
              URI-based — credentials in SIP URI userinfo (assistantID:apiKey)
            </div>
            <div className="font-medium text-gray-700 dark:text-gray-300">
              Media
            </div>
            <div>RTP (direct media, no WebSocket)</div>
          </div>
        </div>
      </FieldSet>

      <FieldSet className="col-span-2">
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          PBX Dial Plan Example
        </div>
        <div>
          <div className="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
            FreeSWITCH
          </div>
          <div className="relative">
            <pre className="dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs overflow-auto rounded">
              {`<extension name="rapida-ai">
  <condition field="destination_number" expression="^(\\d+)$">
    <action application="bridge"
            data="sofia/external/${sipUri}"/>
  </condition>
</extension>`}
            </pre>
            <div className="absolute top-1 right-1">
              <CopyButton className="h-6 w-6 bg-gray-200 dark:bg-gray-800 rounded">
                {`<extension name="rapida-ai">
  <condition field="destination_number" expression="^(\\d+)$">
    <action application="bridge"
            data="sofia/external/${sipUri}"/>
  </condition>
</extension>`}
              </CopyButton>
            </div>
          </div>
        </div>
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          Asterisk (pjsip.conf + extensions.conf)
        </div>
        <div className="relative">
          <pre className="dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs overflow-auto rounded">
            {`; pjsip.conf — register a trunk to Rapida SIP
[rapida-trunk]
type = endpoint
transport = transport-udp
context = from-rapida
aors = rapida-trunk
outbound_auth = rapida-trunk-auth

[rapida-trunk-auth]
type = auth
auth_type = userpass
username = ${assistantId || '{ASSISTANT_ID}'}
password = <YOUR_API_KEY>

[rapida-trunk]
type = aor
contact = sip:${sipHost}:5060

; extensions.conf — route calls to Rapida
[rapida-outbound]
exten => _X.,1,Dial(PJSIP/\${EXTEN}@rapida-trunk)`}
          </pre>
          <div className="absolute top-1 right-1">
            <CopyButton className="h-6 w-6 bg-gray-200 dark:bg-gray-800 rounded">
              {`; pjsip.conf — register a trunk to Rapida SIP
[rapida-trunk]
type = endpoint
transport = transport-udp
context = from-rapida
aors = rapida-trunk
outbound_auth = rapida-trunk-auth

[rapida-trunk-auth]
type = auth
auth_type = userpass
username = ${assistantId || '{ASSISTANT_ID}'}
password = <YOUR_API_KEY>

[rapida-trunk]
type = aor
contact = sip:${sipHost}:5060

; extensions.conf — route calls to Rapida
[rapida-outbound]
exten => _X.,1,Dial(PJSIP/\${EXTEN}@rapida-trunk)`}
            </CopyButton>
          </div>
        </div>
      </FieldSet>
    </>
  );
};

/* -------------------------------------------------------------------------- */
/*  Asterisk Provider Integration Instructions                                */
/* -------------------------------------------------------------------------- */

const AsteriskIntegrationInstructions: FC<{
  mediaHost: string;
  audioSocketHost?: string;
  assistantId: string;
}> = ({ mediaHost, audioSocketHost, assistantId }) => {
  const rapidaHostname = useMemo(() => {
    try {
      return new URL(mediaHost).hostname;
    } catch {
      return '<your-rapida-host>';
    }
  }, [mediaHost]);

  const audioSocketHostPart = useMemo(() => {
    if (!audioSocketHost) return '<your-rapida-host>';
    return audioSocketHost.split(':')[0] || '<your-rapida-host>';
  }, [audioSocketHost]);

  const audioSocketPort = useMemo(() => {
    if (!audioSocketHost) return '4573';
    const parts = audioSocketHost.split(':');
    return parts.length > 1 ? parts[1] : '4573';
  }, [audioSocketHost]);

  const webhookUrl = `https://${rapidaHostname}/v1/talk/asterisk/call/${assistantId || '{ASSISTANT_ID}'}?from=\${CALLERID(num)}&x-api-key={{PROJECT_CREDENTIAL_KEY}}`;

  return (
    <>
      {/* Integrate with WebSocket */}
      <FieldSet className="col-span-2">
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          Integrate with WebSocket
        </div>
        <div className="space-y-3 pt-2">
          <div>
            <FormLabel>Endpoint</FormLabel>
            <div className="flex items-center gap-2">
              <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
                wss://{rapidaHostname}/v1/talk/asterisk/ctx/{'{contextId}'}
              </code>
              <div className="flex shrink-0 border divide-x">
                <CopyButton className="h-7 w-7">
                  {`wss://${rapidaHostname}/v1/talk/asterisk/ctx/{contextId}`}
                </CopyButton>
              </div>
            </div>
          </div>
          <div>
            <FormLabel>Dialplan (extensions.conf)</FormLabel>
            <div className="relative">
              <pre className="dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs overflow-auto rounded">
                {`[rapida-inbound-ws]
exten => _X.,1,Answer()
 same => n,Set(CTX=\${CURL(${webhookUrl})})
 same => n,GotoIf($["\${CTX}" = ""]?error)
 same => n,WebSocket(wss://${rapidaHostname}/v1/talk/asterisk/ctx/\${CTX})
 same => n,Hangup()
 same => n(error),Playback(an-error-has-occurred)
 same => n,Hangup()`}
              </pre>
              <div className="absolute top-1 right-1">
                <CopyButton className="h-6 w-6 bg-gray-200 dark:bg-gray-800 rounded">
                  {`[rapida-inbound-ws]
exten => _X.,1,Answer()
 same => n,Set(CTX=\${CURL(${webhookUrl})})
 same => n,GotoIf($["\${CTX}" = ""]?error)
 same => n,WebSocket(wss://${rapidaHostname}/v1/talk/asterisk/ctx/\${CTX})
 same => n,Hangup()
 same => n(error),Playback(an-error-has-occurred)
 same => n,Hangup()`}
                </CopyButton>
              </div>
            </div>
          </div>
          <InputHelper>
            Requires <code>chan_websocket.so</code> (Asterisk 20+). Uses WSS
            port 443 — ideal for cloud / NAT traversal.
          </InputHelper>
        </div>
      </FieldSet>

      {/* Integrate with AudioSocket */}
      <FieldSet className="col-span-2">
        <div className="font-medium border-b -mx-4 px-4 py-2 text-sm/6">
          Integrate with AudioSocket
        </div>
        <div className="space-y-3 pt-2">
          <div>
            <FormLabel>Endpoint</FormLabel>
            <div className="flex items-center gap-2">
              <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
                {audioSocketHost}
              </code>
              <div className="flex shrink-0 border divide-x">
                <CopyButton className="h-7 w-7">{audioSocketHost}</CopyButton>
              </div>
            </div>
          </div>
          <div>
            <FormLabel>Dialplan (extensions.conf)</FormLabel>
            <div className="relative">
              <pre className="dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs overflow-auto rounded">
                {`[rapida-inbound]
exten => _X.,1,Answer()
 same => n,Set(CHANNEL(audioreadformat)=slin)
 same => n,Set(CHANNEL(audiowriteformat)=slin)
 same => n,Set(CTX=\${CURL(${webhookUrl})})
 same => n,GotoIf($["\${CTX}" = ""]?error)
 same => n,AudioSocket(\${CTX},${audioSocketHostPart}:${audioSocketPort})
 same => n,Hangup()
 same => n(error),Playback(an-error-has-occurred)
 same => n,Hangup()`}
              </pre>
              <div className="absolute top-1 right-1">
                <CopyButton className="h-6 w-6 bg-gray-200 dark:bg-gray-800 rounded">
                  {`[rapida-inbound]
exten => _X.,1,Answer()
 same => n,Set(CHANNEL(audioreadformat)=slin)
 same => n,Set(CHANNEL(audiowriteformat)=slin)
 same => n,Set(CTX=\${CURL(${webhookUrl})})
 same => n,GotoIf($["\${CTX}" = ""]?error)
 same => n,AudioSocket(\${CTX},${audioSocketHostPart}:${audioSocketPort})
 same => n,Hangup()
 same => n(error),Playback(an-error-has-occurred)
 same => n,Hangup()`}
                </CopyButton>
              </div>
            </div>
          </div>
          <InputHelper>
            Requires <code>res_audiosocket.so</code> (Asterisk 16+). Raw TCP
            port {audioSocketPort} — SLIN 16-bit 8 kHz. Best for LAN / private
            network.
          </InputHelper>
        </div>
      </FieldSet>
    </>
  );
};

/* -------------------------------------------------------------------------- */
/*  Voice Input / Output helpers                                              */
/* -------------------------------------------------------------------------- */

const VoiceInput: FC<{ deployment?: DeploymentAudioProvider }> = ({
  deployment,
}) => (
  <div className="">
    <div className="flex items-center space-x-2 border-b py-1 px-4 h-10">
      <h4 className="font-medium">Speech to text</h4>
    </div>
    {deployment?.getAudiooptionsList() ? (
      deployment?.getAudiooptionsList().length > 0 && (
        <div className="text-xs text-gray-500 dark:text-gray-400 py-3 px-3 space-y-6">
          <FieldSet>
            <FormLabel>Provider</FormLabel>
            <ProviderPill provider={deployment?.getAudioprovider()} />
          </FieldSet>
          <div className="grid grid-cols-1 gap-4">
            {deployment
              ?.getAudiooptionsList()
              .filter(d => d.getValue())
              .filter(d => d.getKey().startsWith('listen.'))
              .map((detail, index) => (
                <FieldSet key={index}>
                  <FormLabel>{detail.getKey()}</FormLabel>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
                      {detail.getValue()}
                    </code>
                    <div className="flex shrink-0 border divide-x">
                      <CopyButton className="h-7 w-7">
                        {detail.getValue()}
                      </CopyButton>
                    </div>
                  </div>
                </FieldSet>
              ))}
          </div>
        </div>
      )
    ) : (
      <YellowNoticeBlock>Voice input is not enabled</YellowNoticeBlock>
    )}
  </div>
);

const VoiceOutput: FC<{ deployment?: DeploymentAudioProvider }> = ({
  deployment,
}) => (
  <div>
    <div className="flex items-center space-x-2 border-b py-2 px-4  h-10">
      <h4 className="font-medium">Text to speech</h4>
    </div>
    {deployment?.getAudiooptionsList() ? (
      deployment?.getAudiooptionsList().length > 0 && (
        <div className="text-xs text-gray-500 dark:text-gray-400 py-3 px-3 space-y-6">
          <FieldSet>
            <FormLabel>Provider</FormLabel>
            <ProviderPill provider={deployment?.getAudioprovider()} />
          </FieldSet>
          <div className="grid grid-cols-1 gap-4">
            {deployment
              ?.getAudiooptionsList()
              .filter(d => d.getValue())
              .filter(d => d.getKey().startsWith('speak.'))
              .map((detail, index) => (
                <FieldSet key={index}>
                  <FormLabel>{detail.getKey()}</FormLabel>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 dark:bg-gray-950 bg-gray-100 px-3 py-2 font-mono text-xs min-w-0 overflow-hidden">
                      {detail.getValue()}
                    </code>

                    <div className="flex shrink-0 border divide-x">
                      <CopyButton className="h-7 w-7">
                        {detail.getValue()}
                      </CopyButton>
                    </div>
                  </div>
                </FieldSet>
              ))}
          </div>
        </div>
      )
    ) : (
      <YellowNoticeBlock>Voice output is not enabled</YellowNoticeBlock>
    )}
  </div>
);
