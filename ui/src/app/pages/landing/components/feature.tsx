import {
  Activity,
  BarChart3,
  BookOpen,
  Cpu,
  Globe2,
  PhoneCall,
  Waves,
  Webhook,
} from 'lucide-react';

export const Features = () => {
  return (
    <section className="mt-20 sm:mt-40">
      <div className="border-y px-4 py-2 sm:px-2">
        <h2 className="max-w-3xl text-3xl font-medium tracking-tight text-pretty md:text-[2.5rem]/14">
          A Feature Rich Platform to Build and Deploy Trusted Voice-First
          Experiences
        </h2>
        <p className="mt-4 max-w-2xl text-base text-gray-600 dark:text-gray-400">
          Rapida empowers teams to create intelligent, multi-channel voice
          agents with ease. From real-time speech streaming to custom business
          logic workflows and context-aware responses
        </p>
      </div>

      <section className="border-y relative grid grid-cols-1 gap-y-6 gap-x-10 xl:gap-x-6 2xl:gap-x-10 md:grid-cols-[minmax(300px,1fr)_minmax(0,340px)] lg:grid-cols-[minmax(300px,1fr)_repeat(2,minmax(0,340px))] xl:grid-cols-[minmax(300px,1fr)_repeat(3,minmax(0,340px))] mt-16 lg:mt-20">
        <div className="grid grid-rows-[1fr_auto] md:border-r ">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <Globe2 className="w-6 h-6 opacity-70 mt-4" strokeWidth={1.5} />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Distributed Edge Network
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Deploy across a global low-latency edge network to ensure
                sub-100 ms round-trip voice and inference performance anywhere.
              </p>
            </div>
          </div>
          <div className="border-t px-4 py-2 max-md:border-y sm:px-2">
            <p className="text-sm/6 text-gray-600 dark:text-gray-400">
              Now live in Singapore and India.
            </p>
          </div>
        </div>
        <div className="border-x grid grid-rows-[1fr_auto] max-md:border-t max-xl:last:hidden max-lg:nth-[3]:hidden last:border-r-0 max-xl:nth-[3]:border-r-0 max-lg:nth-[2]:border-r-0">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <PhoneCall
                className="w-6 h-6 opacity-70 mt-4"
                strokeWidth={1.5}
              />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Telephony & SIP Integration
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Integrate seamlessly with SIP trunks and existing telephony
                systems — no number changes or infrastructure rewiring required.
              </p>
            </div>
          </div>
        </div>
        <div className="border-x grid grid-rows-[1fr_auto] max-md:border-t max-xl:last:hidden max-lg:nth-[3]:hidden last:border-r-0 max-xl:nth-[3]:border-r-0 max-lg:nth-[2]:border-r-0">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <Waves className="w-6 h-6 opacity-70 mt-4" strokeWidth={1.5} />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Noise Cancellation & Audio Enhancement
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Deliver crystal-clear call quality with real-time background
                noise suppression, gain control, and adaptive echo cancellation.
              </p>
            </div>
          </div>
        </div>
        <div className="border-x grid grid-rows-[1fr_auto] max-md:border-t max-xl:last:hidden max-lg:nth-[3]:hidden last:border-r-0 max-xl:nth-[3]:border-r-0 max-lg:nth-[2]:border-r-0">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <Cpu className="w-6 h-6 opacity-70 mt-4" strokeWidth={1.5} />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Bring Your Own LLM / STT / TTS
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Orchestrate across any AI provider — OpenAI, Anthropic,
                Deepgram, ElevenLabs, or on-prem — without vendor lock-in.
              </p>
            </div>
          </div>
        </div>
      </section>
      <section className="border-y relative grid grid-cols-1 gap-y-6 gap-x-10 xl:gap-x-6 2xl:gap-x-10 md:grid-cols-[minmax(300px,1fr)_minmax(0,340px)] lg:grid-cols-[minmax(300px,1fr)_repeat(2,minmax(0,340px))] xl:grid-cols-[minmax(300px,1fr)_repeat(3,minmax(0,340px))] mt-16 lg:mt-20">
        <div className="grid grid-rows-[1fr_auto] md:border-r ">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <Activity className="w-6 h-6 opacity-70 mt-4" strokeWidth={1.5} />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Observability & Telemetry
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Trace every millisecond across audio, network, and inference
                layers with full visibility for performance and reliability
                metrics.
              </p>
            </div>
          </div>
        </div>
        <div className="border-x grid grid-rows-[1fr_auto] max-md:border-t max-xl:last:hidden max-lg:nth-[3]:hidden last:border-r-0 max-xl:nth-[3]:border-r-0 max-lg:nth-[2]:border-r-0">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <BarChart3
                className="w-6 h-6 opacity-70 mt-4"
                strokeWidth={1.5}
              />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Post-Conversation Analysis
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Automatically summarize and score every interaction — from
                sentiment to SOP adherence and quality outcomes.
              </p>
            </div>
          </div>
        </div>
        <div className="border-x grid grid-rows-[1fr_auto] max-md:border-t max-xl:last:hidden max-lg:nth-[3]:hidden last:border-r-0 max-xl:nth-[3]:border-r-0 max-lg:nth-[2]:border-r-0">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <BookOpen className="w-6 h-6 opacity-70 mt-4" strokeWidth={1.5} />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Knowledge Integration
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Connect enterprise data sources or knowledge bases to ground
                conversations in accurate, contextual information.
              </p>
            </div>
          </div>
        </div>
        <div className="border-x grid grid-rows-[1fr_auto] max-md:border-t max-xl:last:hidden max-lg:nth-[3]:hidden last:border-r-0 max-xl:nth-[3]:border-r-0 max-lg:nth-[2]:border-r-0">
          <div className="grid grid-cols-1 items-center">
            <div className="px-4 py-2 sm:px-2">
              <Webhook className="w-6 h-6 opacity-70 mt-4" strokeWidth={1.5} />
              <div className="flex items-center gap-2 mt-4">
                <h3 className="text-base/7 font-semibold">
                  <div className="absolute inset-0" />
                  Webhooks & Event Triggers
                </h3>
              </div>

              <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                Extend workflows with programmable hooks and real-time events to
                integrate with CRMs, analytics, or automation systems.
              </p>
            </div>
          </div>
        </div>
      </section>
    </section>
  );
};
