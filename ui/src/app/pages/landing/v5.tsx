import { RapidaIcon } from '@/app/components/Icon/Rapida';
import { RapidaTextIcon } from '@/app/components/Icon/RapidaText';
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/react';
import {
  BarChart,
  ChevronRight,
  Clock,
  Puzzle,
  ShieldCheck,
} from 'lucide-react';
import { VoiceAgent, ConnectionConfig, AgentConfig } from '@rapidaai/react';
import { AuthContext } from '@/context/auth-context';
import { useContext } from 'react';
import { DemoVoiceAgent } from '@/app/pages/landing/agent';
import { Usecases } from '@/app/pages/landing/components/usecase';
import { Pricing } from '@/app/pages/landing/components/pricing';
import { Question } from '@/app/pages/landing/components/faq';
import { Features } from '@/app/pages/landing/components/feature';
import { Footer } from '@/app/pages/landing/components/footer';
import { AgentStudio } from '@/app/pages/landing/components/agent-studio';
import { Testimonial } from '@/app/pages/landing/components/testimonial';
import { Orchestrator } from '@/app/pages/landing/components/orchestrator';
import { TraceIcon } from '@/app/components/Icon/Trace';
import { Header } from '@/app/pages/landing/components/header';
export const V5 = () => {
  return (
    <div
      className="font-sans text-gray-950 dark:text-white antialiased [overflow-anchor:none] container mx-auto"
      scroll-region=""
    >
      <div className="isolate">
        <Header />
        <main className="flex min-h-dvh flex-col pt-14">
          <div className="grid flex-1 grid-rows-[1fr_auto] overflow-clip grid-cols-[1fr_var(--gutter-width)_minmax(0,var(--breakpoint-2xl))_var(--gutter-width)_1fr] [--gutter-width:--spacing(6)] lg:[--gutter-width:--spacing(10)]">
            <div className="col-start-2 row-span-full row-start-1 max-sm:hidden border-x  bg-size-[10px_10px] bg-fixed bg-[repeating-linear-gradient(315deg,theme('colors.gray.200')_0,theme('colors.gray.200')_1px,transparent_0,transparent_50%)] dark:bg-[repeating-linear-gradient(315deg,theme('colors.gray.800')_0,theme('colors.gray.800')_1px,transparent_0,transparent_50%)]" />
            <div className="col-start-4 row-span-full row-start-1 max-sm:hidden border-x  bg-size-[10px_10px] bg-fixed bg-[repeating-linear-gradient(315deg,theme('colors.gray.200')_0,theme('colors.gray.200')_1px,transparent_0,transparent_50%)] dark:bg-[repeating-linear-gradient(315deg,theme('colors.gray.800')_0,theme('colors.gray.800')_1px,transparent_0,transparent_50%)]" />

            {/* hero */}
            <div className="col-start-3 row-start-1 max-sm:col-span-full max-sm:col-start-1">
              <div className="border-y mt-12 grid gap-x-10 sm:mt-20 lg:mt-24 lg:grid-cols-[3fr_2fr]">
                <div className="px-4 py-2 sm:px-2">
                  <p className="font-mono text-base font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
                    Build Voice AI
                  </p>
                  <h1 className="mt-2 text-6xl sm:text-8xl text-pretty">
                    The voice orchestration platform for scale in production
                  </h1>
                </div>
              </div>

              <div className="flex gap-2 px-4 py-2 whitespace-nowrap sm:px-2 border-b">
                <a
                  className="pl-4 px-2 py-2 gap-2 inline-flex justify-center rounded-full text-base font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 bg-blue-600 text-white"
                  href="/demo"
                >
                  Consult an expert
                  <span className="bg-blue-800 rounded-full p-2 w-6 h-6 items-center flex">
                    <svg
                      fill="currentColor"
                      aria-hidden="true"
                      viewBox="0 0 10 10"
                      className="-mr-0.5 w-4"
                    >
                      <path d="M4.85355 0.146423L9.70711 4.99998L4.85355 9.85353L4.14645 9.14642L7.79289 5.49998H0V4.49998H7.79289L4.14645 0.85353L4.85355 0.146423Z" />
                    </svg>
                  </span>
                </a>
                {/* <DemoVoiceAgent
                      voiceAgent={
                        new VoiceAgent(
                          ConnectionConfig.DefaultConnectionConfig(
                            ConnectionConfig.WithSDK({
                              ApiKey:
                                'bca387c4e5cb8fd4cbcaeb194389216959c14dcaee4069d76c52093f5f571171',
                              UserId: '2011133225207857152',
                            }),
                          ),
                          new AgentConfig('2123740035458007040'), //'2181324656735158272') //
                          // .addKeywords([user?.name!])
                          // .addArgument('name', user?.name!),
                        )
                      }
                    /> */}
              </div>
              <TabGroup className="mt-20">
                <TabList className="border-y grid grid-cols-3 divide-x">
                  <Tab className="group relative flex items-center justify-center gap-4 p-4 focus:not-data-focus:outline-hidden data-selected:text-blue-600 flex-col sm:p-6 dark:data-selected:text-blue-500">
                    <div className="absolute -inset-x-px inset-y-0 bg-blue-500/5 not-group-data-selected:hidden dark:bg-blue-500/5" />
                    <svg
                      className="w-20 shrink-0 sm:w-30 data-lift:*:transition-transform data-lift:*:duration-300 data-lift:*:group-hover:-translate-y-0.5 data-lift:*:group-data-selected:translate-y-0"
                      aria-hidden="true"
                      fill="none"
                      viewBox="0 0 120 72"
                    >
                      <path
                        className="fill-white dark:fill-gray-950"
                        d="M56.095 7 8.464 34.5c-.957.553-1.435 1.276-1.435 2v3c0 .724.478 1.448 1.435 2L56.095 69c1.913 1.105 5.015 1.105 6.928 0l47.632-27.5c.956-.552 1.435-1.276 1.435-2v-3c-.001-.724-.479-1.447-1.435-2L63.023 7c-1.913-1.104-5.015-1.104-6.928 0"
                      />
                      <path
                        stroke="currentColor"
                        strokeOpacity="0.4"
                        d="M112.09 36.5c-.001-.724-.479-1.447-1.435-2L63.023 7c-1.913-1.104-5.015-1.104-6.928 0L8.464 34.5c-.957.553-1.435 1.276-1.435 2m105.061 0c0 .724-.479 1.448-1.435 2L63.023 66c-1.913 1.105-5.015 1.105-6.928 0L8.464 38.5c-.957-.552-1.435-1.276-1.435-2m105.061 0v3c0 .724-.479 1.448-1.435 2L63.023 69c-1.913 1.105-5.015 1.105-6.928 0L8.464 41.5c-.957-.552-1.435-1.276-1.435-2v-3"
                      />
                      <path
                        fill="currentColor"
                        stroke="currentColor"
                        d="M11.062 37c-.478-.276-.478-.724 0-1L58.694 8.5c.478-.276 1.253-.276 1.732 0l2.598 1.5c.478.276.478.724 0 1L15.392 38.5c-.478.276-1.253.276-1.732 0z"
                        opacity="0.1"
                      />
                      <g
                        fill="currentColor"
                        stroke="currentColor"
                        opacity="0.1"
                      >
                        <path d="M19.723 42c-.479-.276-.479-.724 0-1l47.63-27.5c.48-.276 1.255-.276 1.733 0L89.004 25c.479.276.479.724 0 1l-47.63 27.5c-.48.276-1.255.276-1.733 0z" />
                        <path d="M34.445 31.5c-.479-.276-.479-.724 0-1L49.167 22c.478-.276 1.254-.276 1.732 0l23.383 13.5c.478.276.478.724 0 1L59.559 45c-.478.276-1.253.276-1.732 0z" />
                      </g>
                      <path
                        fill="currentColor"
                        stroke="currentColor"
                        d="M45.703 57c-.478-.276-.478-.724 0-1l47.632-27.5c.478-.276 1.254-.276 1.732 0l12.99 7.5c.479.276.479.724 0 1L60.426 64.5c-.478.276-1.254.276-1.732 0z"
                        opacity="0.1"
                      />
                      <g data-lift="true">
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M11.062 32c-.478-.276-.478-.724 0-1L58.694 3.5c.478-.276 1.253-.276 1.732 0L63.024 5c.478.276.478.724 0 1L15.392 33.5c-.478.276-1.253.276-1.732 0z"
                        />
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M19.723 37c-.479-.276-.479-.724 0-1l47.63-27.5c.48-.276 1.255-.276 1.733 0L89.004 20c.479.276.479.724 0 1l-47.63 27.5c-.48.276-1.255.276-1.733 0z"
                        />
                        <path
                          stroke="currentColor"
                          strokeOpacity="0.3"
                          d="M37.909 44.5c-.478-.276-.478-.724 0-1l9.526-5.5c.479-.276 1.254-.276 1.732 0l1.732 1c.479.276.479.724 0 1l-9.526 5.5c-.478.276-1.254.276-1.732 0z"
                        />
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M34.445 26.5c-.479-.276-.479-.724 0-1L49.167 17c.478-.276 1.254-.276 1.732 0l23.383 13.5c.478.276.478.724 0 1L59.559 40c-.478.276-1.253.276-1.732 0z"
                        />
                        <path
                          stroke="currentColor"
                          strokeOpacity="0.3"
                          d="M56.096 36c-.479-.276-.479-.724 0-1l9.526-5.5c.478-.276 1.254-.276 1.732 0l1.732 1c.478.276.478.724 0 1L59.56 37c-.479.276-1.254.276-1.732 0zM70.818 25.5c-.478-.276-.478-.724 0-1l9.526-5.5c.479-.276 1.254-.276 1.733 0l1.732 1c.478.276.478.724 0 1l-9.527 5.5c-.478.276-1.254.276-1.732 0z"
                        />
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M45.703 52c-.478-.276-.478-.724 0-1l47.632-27.5c.478-.276 1.254-.276 1.732 0l12.99 7.5c.479.276.479.724 0 1L60.426 59.5c-.478.276-1.254.276-1.732 0z"
                        />
                        <path
                          stroke="currentColor"
                          strokeOpacity="0.3"
                          d="M93.335 34.5c-.478-.276-.478-.724 0-1l6.062-3.5c.478-.276 1.254-.276 1.732 0l1.732 1c.478.276.478.724 0 1l-6.062 3.5c-.478.276-1.254.276-1.732 0zM77.746 43.5c-.478-.276-.478-.724 0-1L89.004 36c.478-.276 1.254-.276 1.732 0l1.732 1c.479.276.479.724 0 1L81.21 44.5c-.478.276-1.254.276-1.732 0z"
                        />
                      </g>
                    </svg>
                    <div className="text-center sm:text-left text-sm xl:flex-1">
                      <p className="font-mono font-semibold tracking-widest uppercase flex">
                        <span className="hidden sm:block mr-1.5">The</span>{' '}
                        Orchestrator
                      </p>
                      <p className="mt-2 max-xl:hidden text-base text-gray-600 dark:text-gray-400">
                        Connect over a hundred voice and language models with
                        memory and logic into a single intelligent workflow.
                      </p>
                    </div>
                  </Tab>
                  {/* <div className="w-px bg-gray-200 dark:bg-gray-800" /> */}
                  <Tab className="group relative flex items-center justify-center gap-4 p-4 focus:not-data-focus:outline-hidden data-selected:text-sky-600 flex-col sm:p-6 dark:data-selected:text-sky-500">
                    <div className="absolute -inset-x-px inset-y-0 bg-sky-500/5 not-group-data-selected:hidden dark:bg-sky-500/5" />
                    <svg
                      className="w-20 shrink-0 sm:w-30 data-lift:*:transition-transform data-lift:*:duration-300 data-lift:*:group-hover:-translate-y-0.5 data-lift:*:group-data-selected:translate-y-0"
                      aria-hidden="true"
                      fill="none"
                      viewBox="0 0 120 72"
                    >
                      <g data-lift="true">
                        <path
                          shapeRendering="geometricPrecision"
                          className="fill-white dark:fill-gray-950"
                          d="M56.066 6 8.435 33.5C7.478 34.053 7 34.776 7 35.5v3c0 .724.478 1.448 1.435 2L56.066 68c1.913 1.105 5.015 1.105 6.929 0l47.631-27.5c.957-.552 1.435-1.276 1.435-2v-3c0-.724-.479-1.447-1.435-2L62.995 6c-1.914-1.104-5.015-1.104-6.929 0"
                        />
                        <path
                          shapeRendering="geometricPrecision"
                          stroke="currentColor"
                          d="M112.09 35.496c-.001-.723-.479-1.447-1.435-2l-47.632-27.5c-1.913-1.104-5.015-1.104-6.928 0l-47.631 27.5c-.957.553-1.435 1.277-1.435 2m105.061 0c0 .724-.479 1.448-1.435 2l-47.632 27.5c-1.913 1.105-5.015 1.105-6.928 0l-47.631-27.5c-.957-.552-1.435-1.276-1.435-2m105.061 0v3c0 .724-.479 1.448-1.435 2l-47.632 27.5c-1.913 1.105-5.015 1.105-6.928 0l-47.631-27.5c-.957-.552-1.435-1.276-1.435-2v-3"
                        />
                        <path
                          shapeRendering="geometricPrecision"
                          stroke="currentColor"
                          strokeOpacity="0.3"
                          d="M11.062 35.996c-.478-.276-.478-.724 0-1l47.632-27.5c.478-.276 1.253-.276 1.732 0l30.31 17.5c.479.277.479.724 0 1l-47.63 27.5c-.479.276-1.255.276-1.733 0zM45.703 55.996c-.478-.276-.478-.724 0-1l47.632-27.5c.478-.276 1.254-.276 1.732 0l12.99 7.5c.479.276.479.724 0 1l-47.631 27.5c-.478.276-1.254.276-1.732 0z"
                        />
                        <circle
                          shapeRendering="geometricPrecision"
                          cx="1.5"
                          cy="1.5"
                          r="1.5"
                          fill="currentColor"
                          transform="matrix(.86603 -.5 .86603 .5 16.258 35.496)"
                        />
                        <path
                          shapeRendering="geometricPrecision"
                          stroke="currentColor"
                          strokeLinecap="round"
                          d="m22.32 33.496 3.464-2M56.961 13.496l3.465-2M49.168 17.996l4.33-2.5M42.24 21.996l3.463-2"
                        />
                        <path
                          stroke="currentColor"
                          strokeLinecap="round"
                          strokeOpacity="0.3"
                          d="m41.373 38.496 23.383-13.5"
                        />
                        <path
                          shapeRendering="geometricPrecision"
                          stroke="currentColor"
                          strokeLinecap="round"
                          d="m53.498 55.496 6.928-4M69.086 46.496l6.928-4M84.674 37.496l6.929-4"
                        />
                        <path
                          shapeRendering="geometricPrecision"
                          stroke="currentColor"
                          strokeLinecap="round"
                          strokeOpacity="0.3"
                          d="m56.096 56.996 9.526-5.5M71.684 47.996l9.526-5.5M87.273 38.996l9.526-5.5M58.693 58.496l8.66-5M74.282 49.496l8.66-5M89.87 40.496l8.66-5M46.57 38.496l18.186-10.5"
                        />
                        <rect
                          shapeRendering="geometricPrecision"
                          width={28}
                          height={2}
                          fill="currentColor"
                          rx="0.5"
                          transform="matrix(.86603 -.5 .86603 .5 33.579 34.496)"
                        />
                        <rect
                          shapeRendering="geometricPrecision"
                          width={32}
                          height={2}
                          fill="currentColor"
                          rx="0.5"
                          transform="matrix(.86603 -.5 .86603 .5 35.311 37.496)"
                        />
                        <rect
                          shapeRendering="geometricPrecision"
                          width={10}
                          height={3}
                          fill="currentColor"
                          rx="1.5"
                          transform="matrix(.86603 -.5 .86603 .5 48.301 39.996)"
                        />
                        <rect
                          shapeRendering="geometricPrecision"
                          width={10}
                          height={3}
                          fill="currentColor"
                          fillOpacity="0.3"
                          rx="1.5"
                          transform="matrix(.86603 -.5 .86603 .5 58.693 33.996)"
                        />
                      </g>
                    </svg>
                    <div className="text-center sm:text-left text-sm xl:flex-1">
                      <p className="font-mono font-semibold tracking-widest uppercase flex">
                        <span className="hidden sm:block mr-1.5">Agent</span>{' '}
                        Studio
                      </p>
                      <p className="mt-2 max-xl:hidden text-base text-gray-600 dark:text-gray-400">
                        Design, test, and launch fully customizable voice
                        experiences that are ready for production.
                      </p>
                    </div>
                  </Tab>
                  {/* <div className="w-px bg-gray-200 dark:bg-gray-800" /> */}
                  <Tab className="group relative flex items-center justify-center gap-4 p-4 focus:not-data-focus:outline-hidden data-selected:text-emerald-600 flex-col sm:p-6 dark:data-selected:text-emerald-500">
                    <div className="absolute -inset-x-px inset-y-0 bg-emerald-500/5 not-group-data-selected:hidden dark:bg-emerald-500/5" />
                    <svg
                      className="w-20 shrink-0 sm:w-30 data-lift:*:transition-transform data-lift:*:duration-300 data-lift:*:group-hover:-translate-y-0.5 data-lift:*:group-data-selected:translate-y-0"
                      aria-hidden="true"
                      fill="none"
                      viewBox="0 0 120 72"
                    >
                      <path
                        className="fill-white dark:fill-gray-950"
                        d="M56.095 6 8.464 33.5c-.957.553-1.435 1.276-1.435 2v3c0 .724.478 1.448 1.435 2L56.095 68c1.913 1.105 5.015 1.105 6.928 0l47.632-27.5c.956-.552 1.435-1.276 1.435-2v-3c-.001-.724-.479-1.447-1.435-2L63.023 6c-1.913-1.104-5.015-1.104-6.928 0"
                      />
                      <g stroke="currentColor" opacity="0.1">
                        <path
                          fill="currentColor"
                          d="M60.425 52.5c-.478-.276-.478-.724 0-1L87.272 36c.478-.276 1.254-.276 1.732 0l1.732 1c.478.276.478.724 0 1L63.89 53.5c-.478.276-1.253.276-1.732 0zM54.363 49c-.956-.552-.956-1.448 0-2l3.464-2c.957-.552 2.508-.552 3.464 0 .957.552.957 1.448 0 2l-3.464 2c-.956.552-2.507.552-3.464 0Z"
                        />
                        <path strokeLinecap="round" d="m63.89 43.5 12.124-7" />
                        <path
                          fill="currentColor"
                          d="M46.57 44.5c-.48-.276-.48-.724 0-1L73.415 28c.478-.276 1.254-.276 1.732 0l1.732 1c.478.276.478.724 0 1L50.033 45.5c-.478.276-1.254.276-1.732 0z"
                        />
                        <path strokeLinecap="round" d="m43.105 42.5 10.392-6" />
                        <path
                          fill="currentColor"
                          d="M37.043 39c-.478-.276-.478-.724 0-1L63.89 22.5c.478-.276 1.253-.276 1.732 0l1.732 1c.478.276.478.724 0 1L40.507 40c-.478.276-1.254.276-1.732 0z"
                        />
                        <path strokeLinecap="round" d="m33.579 37 10.392-6" />
                      </g>
                      <path
                        stroke="currentColor"
                        strokeOpacity="0.4"
                        d="M112.09 35.5c-.001-.724-.479-1.447-1.435-2L63.023 6c-1.913-1.104-5.015-1.104-6.928 0L8.464 33.5c-.957.553-1.435 1.276-1.435 2m105.061 0c0 .724-.479 1.448-1.435 2L63.023 65c-1.913 1.105-5.015 1.105-6.928 0L8.464 37.5c-.957-.552-1.435-1.276-1.435-2m105.061 0v3c0 .724-.479 1.448-1.435 2L63.023 68c-1.913 1.105-5.015 1.105-6.928 0L8.464 40.5c-.957-.552-1.435-1.276-1.435-2v-3"
                      />
                      <path
                        stroke="currentColor"
                        strokeOpacity="0.4"
                        d="M17.99 40c-.478-.276-.478-.724 0-1l47.632-27.5c.478-.276 1.254-.276 1.732 0L108.057 35c.478.276.478.724 0 1L60.426 63.5c-.479.276-1.254.276-1.732 0z"
                      />
                      <path
                        fill="currentColor"
                        stroke="currentColor"
                        d="M11.062 36c-.478-.276-.478-.724 0-1L58.694 7.5c.478-.276 1.253-.276 1.732 0L63.024 9c.478.276.478.724 0 1L15.392 37.5c-.478.276-1.253.276-1.732 0z"
                        opacity="0.1"
                      />
                      <g data-lift="true">
                        <path
                          className="fill-current"
                          fillOpacity="0.3"
                          stroke="currentColor"
                          d="M60.425 47.5c-.478-.276-.478-.724 0-1L87.272 31c.478-.276 1.254-.276 1.732 0l1.732 1c.478.276.478.724 0 1L63.89 48.5c-.478.276-1.253.276-1.732 0zM54.363 44c-.956-.552-.956-1.448 0-2l3.464-2c.957-.552 2.508-.552 3.464 0 .957.552.957 1.448 0 2l-3.464 2c-.956.552-2.507.552-3.464 0Z"
                        />
                        <circle
                          cx={2}
                          cy={2}
                          r={2}
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          transform="matrix(.86603 -.5 .86603 .5 56.095 41)"
                        />
                        <path
                          stroke="currentColor"
                          strokeLinecap="round"
                          d="m63.89 38.5 12.124-7"
                        />
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M46.57 39.5c-.48-.276-.48-.724 0-1L73.415 23c.478-.276 1.254-.276 1.732 0l1.732 1c.478.276.478.724 0 1L50.033 40.5c-.478.276-1.254.276-1.732 0z"
                        />
                        <path
                          stroke="currentColor"
                          strokeLinecap="round"
                          d="m43.105 37.5 10.392-6"
                        />
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M37.043 34c-.478-.276-.478-.724 0-1L63.89 17.5c.478-.276 1.253-.276 1.732 0l1.732 1c.478.276.478.724 0 1L40.507 35c-.478.276-1.254.276-1.732 0z"
                        />
                        <path
                          stroke="currentColor"
                          strokeLinecap="round"
                          d="m33.579 32 10.392-6"
                        />
                        <path
                          className="fill-white dark:fill-gray-950"
                          stroke="currentColor"
                          d="M11.062 31c-.478-.276-.478-.724 0-1L58.694 2.5c.478-.276 1.253-.276 1.732 0L63.024 4c.478.276.478.724 0 1L15.392 32.5c-.478.276-1.253.276-1.732 0z"
                        />
                      </g>
                    </svg>
                    <div className="text-center sm:text-left text-sm xl:flex-1">
                      <p className="font-mono font-semibold tracking-widest uppercase">
                        Observability
                      </p>
                      <p className="mt-2 max-xl:hidden text-base text-gray-600 dark:text-gray-400">
                        Monitor, measure, and improve accuracy and performance
                        in real time.
                      </p>
                    </div>
                  </Tab>
                </TabList>

                <div className="border-y mt-4">
                  <TabPanels className="bg-gray-950/5 dark:bg-white/5">
                    <TabPanel className={'sm:px-3 py-3'}>
                      <Orchestrator />
                    </TabPanel>
                    <TabPanel className={'sm:px-3 py-3'}>
                      <AgentStudio />
                    </TabPanel>
                    <TabPanel className={'sm:px-3 py-3'}>
                      <Observability />
                    </TabPanel>
                  </TabPanels>
                </div>
              </TabGroup>
              <Testimonial />
              <Usecases />
              <Features />
              <Pricing />
              <Question />
            </div>

            <Footer />
          </div>
        </main>
      </div>
    </div>
  );
};

export const Observability = () => {
  return (
    <div className="border sm:rounded-2xl dark:bg-gray-900 bg-gray-100 pt-10 flex flex-col lg:flex-row h-fit sm:h-[500px]  w-full mx-auto overflow-hidden text-gray-600 dark:text-gray-400">
      <div className="grid grid-cols-1 gap-8 md:grid-cols-2 md:gap-24 items-center ">
        <div className="px-6">
          <span className="uppercase text-xs font-medium font-mono">
            Voice Telemetry — Every Minute Counts
          </span>
          <p className="text-4xl mt-8 font-semibold text-balance">
            When you’re running thousands of voice agents across regions and
            carriers, visibility isn’t optional — it’s survival.
          </p>

          <div className="mt-6 text-xs font-medium grid grid-cols-1 sm:grid-cols-2 md:grid-cols-1 lg:grid-cols-2 gap-2">
            <div className="inline-flex items-center gap-2 text-xs">
              <TraceIcon className="size-4" strokeWidth={1.5} />
              <span className="font-medium text-sm">
                End-to-End Traceability
              </span>
            </div>
            <div className="inline-flex items-center gap-2 text-xs">
              <Clock className="size-4" strokeWidth={1.5} />
              <span className="font-medium text-sm">
                Sub-Second Performance Insights
              </span>
            </div>
            <div className="inline-flex items-center gap-2 text-xs">
              <Puzzle className="size-4" strokeWidth={1.5} />
              <span className="font-medium text-sm">
                Operational Intelligence for Scale
              </span>
            </div>
            <div className="inline-flex items-center gap-2 text-xs">
              <BarChart className="size-4" strokeWidth={1.5} />
              <span className="font-medium text-sm">
                Real-Time Observability Dashboards
              </span>
            </div>
            <div className="inline-flex items-center gap-2 text-xs">
              <ShieldCheck className="size-4" strokeWidth={1.5} />
              <span className="font-medium text-sm">
                Built for Production Reliability
              </span>
            </div>
          </div>
        </div>
        <div className="h-full md:order-first">
          <img
            src="/images/screenshots/telemetry-dark.png"
            alt="#_"
            className="hidden dark:block bg-gray-200 shadow-box shadow-gray-500/30 overflow-hidden w-full h-full object-cover object-left rounded-tr-2xl border border-l-0 border-b-0 dark:border-gray-800"
          />
          <img
            src="/images/screenshots/telemetry-light.png"
            alt="#_"
            className="dark:hidden block bg-gray-200 shadow-box shadow-gray-500/30 overflow-hidden w-full h-full object-cover object-center  rounded-tr-2xl border border-l-0 border-b-0"
          />
        </div>
      </div>
    </div>
  );
};
