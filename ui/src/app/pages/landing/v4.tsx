import { RapidaIcon } from '@/app/components/Icon/Rapida';
import { RapidaTextIcon } from '@/app/components/Icon/RapidaText';
import SmoothScroll from '@/app/pages/landing/smooth-scroll';
import { VoiceAgent, ConnectionConfig, AgentConfig } from '@rapidaai/react';
import { DemoVoiceAgent } from '@/app/pages/landing/agent';
import { AuthContext } from '@/context/auth-context';
import { useContext } from 'react';

import { IBlueBGButton } from '@/app/components/Form/Button';
import { ChevronRight } from 'lucide-react';
import { motion } from 'framer-motion';

export const V4 = () => {
  const { isAuthenticated } = useContext(AuthContext);
  return (
    <SmoothScroll className="w-full min-h-dvh min-w-full antialiased text-base text-gray-700 dark:text-gray-400 dark:bg-gray-950 ">
      <div className="max-w-7xl mx-auto">
        <div className="bg-background flex border-x dark:border-gray-900 px-4 lg:px-8">
          <a
            className="text-text-primary block grow py-2 text-sm/[1.2] "
            href="/auth/signup"
          >
            We're getting ready to launch our public version —
            <span className="underline font-medium">Join the waitlist</span>
          </a>
          <button className="text-text-secondary hover:text-text-primary group flex cursor-pointer items-center text-xs/[1.2] transition-colors space-x-1.5">
            <span>Dismiss</span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width={24}
              height={24}
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth={1.5}
              strokeLinecap="round"
              strokeLinejoin="round"
              className="lucide lucide-x text-text-secondary group-hover:text-text-primary w-3.5 transition-colors"
            >
              <path d="M18 6 6 18" />
              <path d="m6 6 12 12" />
            </svg>
          </button>
        </div>
      </div>
      <header className="z-50 border-y dark:border-gray-900 sticky left-0 right-0 top-0">
        <div className="max-w-7xl mx-auto">
          <div className="relative px-4 sm:px-8 bg-paper flex h-18.25 items-center justify-between w-full border-x dark:border-gray-900">
            <div className="block">
              <a
                // className="w-32 shrink-0 focus:outline-offset-4 focus:outline-alpha2 flex text-blue-600"
                href="/"
                className="flex items-center shrink-0 space-x-1.5 ml-1 text-blue-600 dark:text-blue-500"
              >
                <RapidaIcon className="h-8 w-8" />
                <RapidaTextIcon className="h-6 shrink-0 ml-1" />
              </a>
            </div>

            <nav className="flex h-full w-full select-none flex-wrap items-center leading-none">
              <div className="flex ml-auto items-center whitespace-nowrap">
                <a
                  href={
                    isAuthenticated && isAuthenticated()
                      ? '/dashboard'
                      : '/auth/signin'
                  }
                  className="group relative inline-flex items-center justify-center overflow-hidden border border-blue-600 pl-8 pr-3 font-medium  text-white transition duration-300 ease-out rounded-[2px]"
                >
                  <ChevronRight
                    className="absolute w-4 h-4 mr-2 z-10 left-2 my-auto text-white group-hover:text-blue-600 dark:group-hover:text-blue-500"
                    strokeWidth={1.5}
                  />
                  <span className="ease absolute inset-0 flex h-full w-full -translate-x-full items-center justify-center duration-500 group-hover:translate-x-0 text-white group-hover:text-blue-600 dark:group-hover:text-blue-500">
                    {isAuthenticated && isAuthenticated()
                      ? 'Dashboard'
                      : 'Sign in'}
                  </span>
                  <span className="ease absolute left-0 pl-3 right-0 flex h-full w-full transform items-center justify-center transition-all duration-500 group-hover:translate-x-full bg-blue-600">
                    Start building
                  </span>
                  <span className="invisible relative">
                    <ChevronRight className="w-4 h-4" strokeWidth={1.5} />
                    Start building
                  </span>
                </a>
              </div>
            </nav>
          </div>
        </div>
      </header>
      <main>
        <div className="max-w-7xl mx-auto ">
          <div className="py-20 border-r border-l dark:border-gray-900"></div>
        </div>
        <div className="max-w-7xl mx-auto px-4 sm:px-8">
          <div className="flex flex-col gap-6 relative col-span-2">
            <div className="space-y-10 pt-10 md:pt-0 mb-8 md:mb-0">
              <h4 className="text-blue-600 block text-sm font-medium uppercase tracking-[1.8px]">
                Build Voice AI
              </h4>
              <h1 className="font-medium text-3xl md:text-4xl lg:text-5xl tracking-tight mt-10">
                The voice orchestration platform
                <br /> for{' '}
                <span className="text-blue-600">scale in production</span>
              </h1>
            </div>

            <div className="flex gap-3">
              <a
                target="_blank"
                href="https://cal.com/prashant-srivastav-u8duzh/30min"
                className="w-fit"
                rel="noreferrer"
              >
                <IBlueBGButton className="px-4 font-medium">
                  Get a demo
                </IBlueBGButton>
              </a>

              <DemoVoiceAgent
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
              />
            </div>
          </div>
        </div>

        <div className="max-w-7xl mx-auto hidden md:block">
          <div className="py-20 border-r border-l dark:border-gray-900"></div>
        </div>
        <div className="w-full border-b dark:text-gray-500 border-t dark:border-gray-900 hidden md:block">
          <div className="align-center grid items-center overflow-hidden mx-auto max-w-7xl grid-cols-6">
            <div className="flex py-3 items-center justify-center border-x dark:border-gray-900 text-center text-base font-medium lg:justify-start px-4 h-full">
              <span className="font-normal">
                Powering apps for multiple languages
              </span>
            </div>
            <div className="relative py-3 items-center justify-center overflow-hidden lg:border-r dark:border-gray-900 lg:border-t-0">
              <div className="flex-col flex h-full w-full items-center justify-center space-y-2">
                <div className="font-semibold text-2xl">हिन्दी</div>
                <div className="latin">Hindi</div>
              </div>
            </div>
            <div className="relative py-3 items-center justify-center overflow-hidden border-r border-t lg:border-r dark:border-gray-900 lg:border-t-0">
              <div className="flex-col flex h-full w-full items-center justify-center space-y-2">
                <div className="font-semibold text-2xl">తెలుగు</div>
                <div className="latin">Telugu</div>
              </div>
            </div>
            <div className="relative py-3 items-center justify-center overflow-hidden border-r border-t lg:border-r dark:border-gray-900 lg:border-t-0">
              <div className="flex-col flex h-full w-full items-center justify-center space-y-2">
                <div className="font-semibold text-2xl">ಕನ್ನಡ</div>
                <div className="latin">Kannada</div>
              </div>
            </div>
            <div className="relative py-3 items-center justify-center overflow-hidden border-r border-t lg:border-r dark:border-gray-900 lg:border-t-0">
              <div className="flex-col flex h-full w-full items-center justify-center space-y-2">
                <div className="font-semibold text-2xl">தமிழ்</div>
                <div className="latin">Tamil</div>
              </div>
            </div>
            <div className="relative py-3 items-center justify-center overflow-hidden border-r border-t lg:border-r dark:border-gray-900 lg:border-t-0">
              <div className="flex-col flex h-full w-full items-center justify-center space-y-2">
                <div className="font-semibold text-2xl">Bahasa</div>
                <div className="latin">Indonesian</div>
              </div>
            </div>
          </div>
        </div>
        <div className="max-w-7xl mx-auto">
          <div className="py-20 border-r border-l dark:border-gray-900"></div>
        </div>

        <section className="relative mx-auto w-full max-w-7xl px-4 sm:px-8 pt-16 lg:pt-32 dark:text-gray-400">
          <div className="max-w-3xl space-y-2">
            <div className="text-base uppercase tracking-wider">
              capabilities
            </div>
            <h2 className="z-20 text-3xl leading-[1.15em] lg:text-5xl tracking-tight lg:leading-[1.15em]">
              A feature-rich platform
            </h2>
          </div>
        </section>
        <section className="pt-8">
          <div className="border-y px-4 sm:px-8 dark:border-gray-900">
            <div className="grid grid-cols-1 md:grid-cols-3 border-x max-w-7xl mx-auto dark:border-gray-900">
              <div className="flex justify-center border-r border-b dark:border-gray-900  dark:hover:bg-gray-900 hover:bg-gray-50">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px ">
                  <div className="relative h-full w-full dark:border-gray-800">
                    <div className="relative flex flex-col gap-4 space-y-3 px-6 py-8 ">
                      <div className="relative h-full py-3 w-fit">
                        <svg
                          className="relative h-6 w-6"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 40 40"
                          fill="none"
                        >
                          <path
                            d="M6.66699 34.166L6.16891 34.1223L6.1212 34.666H6.66699V34.166ZM33.3337 34.166V34.666H33.8794L33.8317 34.1223L33.3337 34.166ZM26.167 10.8327C26.167 14.2384 23.4061 16.9993 20.0003 16.9993V17.9993C23.9584 17.9993 27.167 14.7907 27.167 10.8327H26.167ZM20.0003 16.9993C16.5946 16.9993 13.8337 14.2384 13.8337 10.8327H12.8337C12.8337 14.7907 16.0423 17.9993 20.0003 17.9993V16.9993ZM13.8337 10.8327C13.8337 7.42693 16.5946 4.66602 20.0003 4.66602V3.66602C16.0423 3.66602 12.8337 6.87464 12.8337 10.8327H13.8337ZM20.0003 4.66602C23.4061 4.66602 26.167 7.42693 26.167 10.8327H27.167C27.167 6.87464 23.9584 3.66602 20.0003 3.66602V4.66602ZM7.16508 34.2097C7.80325 26.9364 12.8631 21.3327 20.0003 21.3327V20.3327C12.2391 20.3327 6.84253 26.445 6.16891 34.1223L7.16508 34.2097ZM20.0003 21.3327C27.1376 21.3327 32.1974 26.9364 32.8356 34.2097L33.8317 34.1223C33.1581 26.445 27.7616 20.3327 20.0003 20.3327V21.3327ZM6.66699 34.666H33.3337V33.666H6.66699V34.666Z"
                            fill="currentColor"
                          />
                        </svg>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg font-medium">
                          SSO and RBAC for teams
                        </div>
                        <div className="text-base dark:text-gray-500">
                          Set up roles and collaborate on your realtime project
                          with your team.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="flex justify-center border-r border-b dark:border-gray-900 hover:bg-gray-50 dark:hover:bg-gray-900">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px ">
                  <div className="relative h-full w-full">
                    <div className="relative flex flex-col gap-4 space-y-3 px-6 py-8 text-fg0 contain-content">
                      <div className="relative h-full w-fit overflow-hidden transition-all duration-300 p-px ">
                        <div className="relative h-full w-full py-3">
                          <svg
                            className="relative h-6 w-6"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 40 40"
                            fill="none"
                          >
                            <path
                              d="M10.833 30.6257H5.83301V5.83398H34.1663V30.6257H29.1663M10.833 30.6257V34.1673M10.833 30.6257H29.1663M29.1663 30.6257V34.1673M19.1663 18.334H12.9163M27.4997 18.334C27.4997 22.4761 24.1418 25.834 19.9997 25.834C15.8575 25.834 12.4997 22.4761 12.4997 18.334C12.4997 14.1918 15.8575 10.834 19.9997 10.834C24.1418 10.834 27.4997 14.1918 27.4997 18.334Z"
                              stroke="currentColor"
                              strokeLinecap="square"
                            />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg font-medium">
                          Enterprise-grade security
                        </div>
                        <div className="text-base dark:text-gray-500">
                          End-to-end encryption, SOC2 Type 2, GDPR, CCPA, and
                          HIPAA compliant.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="flex justify-center border-b dark:border-gray-900 hover:bg-gray-50 dark:hover:bg-gray-900">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px ">
                  <div className="relative h-full w-full ">
                    <div className="relative flex flex-col gap-4 space-y-3 px-6 py-8 text-fg0 contain-content">
                      <div className="relative h-full w-fit overflow-hidden  transition-all duration-300 p-px ">
                        <div className="relative h-full w-full py-3">
                          <svg
                            className="relative h-6 w-6"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 20 20"
                            fill="none"
                          >
                            <path
                              d="M6.24959 2.91699V17.0837M2.91699 7.91699V12.0837M9.99959 6.25033V13.7503M13.7503 4.58366V15.417M17.0837 7.91699V12.0837"
                              stroke="currentColor"
                              strokeLinecap="square"
                            />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg font-medium">
                          Noise and echo cancellation
                        </div>
                        <div className="text-base dark:text-gray-500">
                          Remove background artifacts from audio streams.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="flex justify-center border-r dark:border-gray-900 hover:bg-gray-50 dark:hover:bg-gray-900">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px ">
                  <div className="relative h-full w-full ">
                    <div className="relative flex flex-col gap-4 space-y-3 px-6 py-8 text-fg0 contain-content">
                      <div className="relative h-full w-fit overflow-hidden transition-all duration-300 p-px ">
                        <div className="relative h-full w-full py-3">
                          <svg
                            className="relative h-6 w-6"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 40 40"
                            fill="none"
                          >
                            <path
                              d="M4.16699 10.8333V7.5H7.50033M22.5003 7.5H27.5003M12.5003 7.5H17.5003M32.3151 7.5H35.8337V10.8333M17.5003 32.5H12.5003M7.50033 32.5H4.16699V28.9286M4.16699 22.5V17.5"
                              stroke="currentColor"
                              strokeLinecap="square"
                            />
                            <circle
                              cx="30.8337"
                              cy="25.8327"
                              r="4.16667"
                              stroke="currentColor"
                            />
                            <ellipse
                              cx="30.833"
                              cy="25.834"
                              rx="2.5"
                              ry="2.5"
                              stroke="currentColor"
                            />
                            <circle
                              cx="30.8333"
                              cy="25.8333"
                              r="0.833333"
                              stroke="currentColor"
                            />
                            <circle
                              cx="30.8333"
                              cy="25.8333"
                              r="8.33333"
                              stroke="currentColor"
                            />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg font-medium">
                          Session recording
                        </div>
                        <div className="text-base dark:text-gray-500">
                          Composite, record, and store sessions in any
                          S3-compatible cloud bucket.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="flex justify-center border-r dark:border-gray-900 hover:bg-gray-50 dark:hover:bg-gray-900">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px ">
                  <div className="relative h-full w-full ">
                    <div className="relative flex flex-col gap-4 space-y-3 px-6 py-8 text-fg0 contain-content">
                      <div className="relative h-full w-fit overflow-hidden  transition-all duration-300 p-px ">
                        <div className="relative h-full w-full py-3">
                          <svg
                            className="relative h-6 w-6 -translate-x-[2px] -translate-y-[2px]"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 40 40"
                            fill="none"
                          >
                            <path
                              d="M22.4997 10.834H35.833V35.834H10.833V22.5006M16.6663 16.6673L5.83301 5.83398M17.4997 5.83398V17.5006L5.83301 17.5006"
                              stroke="currentColor"
                              strokeLinecap="square"
                            />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg font-medium">
                          Stream ingestion
                        </div>
                        <div className="text-base dark:text-gray-500">
                          Import external video or audio streams from a myriad
                          of formats.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="flex justify-center hover:bg-gray-50 dark:hover:bg-gray-900">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px ">
                  <div className="relative h-full w-full ">
                    <div className="relative flex flex-col gap-4 space-y-3 px-6 py-8 text-fg0 contain-content">
                      <div className="relative h-full w-fit overflow-hidden  transition-all duration-300 p-px ">
                        <div className="relative h-full w-full py-3">
                          <svg
                            className="relative h-6 w-6"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 40 40"
                            fill="none"
                          >
                            <path
                              d="M5.83301 11.6673H22.083M34.1663 28.334H21.2497M5.83301 28.334H8.74967M22.4997 11.6673C22.4997 8.4444 25.1101 5.83398 28.333 5.83398C31.5559 5.83398 34.1663 8.4444 34.1663 11.6673C34.1663 14.8902 31.5559 17.5007 28.333 17.5007C25.1101 17.5007 22.4997 14.8902 22.4997 11.6673ZM20.833 28.334C20.833 31.5569 18.2226 34.1673 14.9997 34.1673C11.7768 34.1673 9.16634 31.5569 9.16634 28.334C9.16634 25.1111 11.7768 22.5007 14.9997 22.5007C18.2226 22.5007 20.833 25.1111 20.833 28.334Z"
                              stroke="currentColor"
                              strokeLinecap="square"
                            />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg font-medium">
                          Moderation tools
                        </div>
                        <div className="text-base dark:text-gray-500">
                          Restrict user access to publish or subscribe to any
                          session or stream.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="relative mx-auto w-full max-w-7xl px-4 sm:px-8 pt-16 lg:pt-32 dark:text-gray-400">
          <div className="flex w-full flex-col gap-10 text-fg0">
            <div className="max-w-3xl space-y-2 text-fg0 undefined">
              <div className=" text-base uppercase tracking-wider">
                Customizability
              </div>
              <h2 className="z-20 text-3xl leading-[1.15em] lg:text-5xl tracking-tight lg:leading-[1.15em] ">
                <span className="">
                  Plug in any model or voice, and talk to it everywhere.
                </span>
              </h2>
            </div>
          </div>
        </section>

        <section className="pt-8">
          <div className="border-y px-4 sm:px-8 dark:border-gray-900">
            <div className="max-w-7xl mx-auto flex flex-col md:grid md:h-min md:auto-rows-min md:grid-cols-9 dark:border-gray-900 border-x ">
              <div className="md:col-span-3">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px border-r dark:border-gray-900 group">
                  <div className="relative h-full group w-full p-6 lg:p-8 lg:pt-8">
                    <div className="relative flex flex-col items-start gap-6  transition-all duration-300">
                      <div className="pt-3.5 pb-1 md:py-1 justify-center gap-2.5 lg:gap-3 items-center flex flex-1 flex-wrap group-hover:grayscale-0 grayscale">
                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[62px] lg:max-w-[62px] h-auto  dark:border-gray-900 flex items-center justify-center">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="p-2 w-12 h-12"
                            viewBox="0 0 24 24"
                          >
                            <path
                              fill="currentColor"
                              d="M11.203 24H1.517a.364.364 0 0 1-.258-.62l6.239-6.275a.366.366 0 0 1 .259-.108h3.52c2.723 0 5.025-2.127 5.107-4.845a5.004 5.004 0 0 0-4.999-5.148H7.613v4.646c0 .2-.164.364-.365.364H.968a.365.365 0 0 1-.363-.364V.364C.605.164.768 0 .969 0h10.416c6.684 0 12.111 5.485 12.01 12.187C23.293 18.77 17.794 24 11.202 24z"
                            />
                          </svg>
                        </div>
                        <div className="border  bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto  dark:border-gray-900">
                          <img
                            className="w-full invert dark:invert-0"
                            src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIYAAACDCAQAAADivaupAAAAIGNIUk0AAHomAACAhAAA+gAAAIDoAAB1MAAA6mAAADqYAAAXcJy6UTwAAAACYktHRAD/h4/MvwAAAAlwSFlzAAAuIwAALiMBeKU/dgAAAAd0SU1FB+gIFhUdFfudmaEAAAO0SURBVHja7Zw9TxRRFIbfmVkUEgkLaGcMIVZGogj+AQtblcRfYKn+BaOFH7GxNrEyGisTg4n8AysSg181WCiFCiiGhd2ZY0Fz954zYXDnXpfkfba7e7J35pmvc8+9OwAhhBBCCCGEEEIIIYQQQgghhBBCCAlJEqsjqbwR1SPrphFLRum+SI+RB1TGLG4i89r+4B5WVOQVXEXR1ZJiBXexGXFrQyKQOdH8kmnRkbeNyM9yNPSpkUb0UaCj2jrmyV8YbZ3wl0lMGX0PZThQhgNlOFCGA2U4BEi6pPwL3VuGvGJqWRpZX4oeJgNNMIPjXrYgOIFX6kzcxiwm1OAkwWtvvxP8wEW0vDbBIr4G2YO6EEgqz0Wk3fUReSmD0uj6ZHJMFqXwIgu5I4kXmcpZWZWOF9mSuTozsVBjk8z47RRt5F5bBxkSYytEZas5Gmpkk9Z7z4s+anWvCNkreK/ImhN0Pk0cKMOBMhwow4EyHHp8muzzLi/GV1YCaSeVScW4/yUDQIZhtVEpdvDbyxQa2Maol2cIRrCJDS9XzbBl9JNjDannMsdOnTJ6dCvAJJ5gzNtIwWO89VKkHOdx3bssE2ziIZZV6zesqjxjECfVoROsYL2fxiZDOIMx1bqOD6ptEtNK/gaW8b5SPy18rG2vg8kQlWIDhZ9rAgIkyFV/BdIqp2ecuS4+TRwow4EyHCjDgTIcKMMhlIyqZZfg86f7oY5KV1tNChfmTgraqswbYTo5rowBI5Wyh186shFv7dDehLpMqu5iH6ngDbQLynCgDAfKcKAMB8pw6D3PSNQMKHZrlSqbEjOyqLIuNs6q4d5lbGHJqIE2MaVqoON4Z9RAJwC7BupRWgOtzUUNBWG7Ov4Il1V1/A1uqOp4E89wWlXHH+C+OjNOYR7jqjp+DfP9dGbkxrFJcAjDqvUw1lS9NMERjKjIIaOfDKNoem0FDtVmoncZJevdk/JwdSewh3So1MolCeGgDAfKcKAMB8pwiLvATfbzAKgQ2XdLEmxyQCVdBQbUBGwDuZqtzaCT/AIZOshV0mX9TeefCVB2k7IVwhO4YKwQXsB3pegcZowVwgv2CuG+qhxqGSWfS8Y/zzZkyoi8ZUR+kqb9u/UR4DKxj5TsTgz4/eXIKv6v1YisGz5NHCjDgTIcKMOBMhwowyGmjNR4kNsTz9ZWRZiijjk2+YKnxlsSfhqRS3hhvCWhhcDEfZlIpSph9ci6iT5qrT2SEEIIIYQQQgghhBBCCCGEEEIIIYQcGP4CpMaZHanWolAAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjQtMDgtMjJUMjE6Mjk6MTMrMDA6MDBvw8jzAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDI0LTA4LTIyVDIxOjI5OjEzKzAwOjAwHp5wTwAAACh0RVh0ZGF0ZTp0aW1lc3RhbXAAMjAyNC0wOC0yMlQyMToyOToyMSswMDowMFCbR1oAAAAZdEVYdFNvZnR3YXJlAHd3dy5pbmtzY2FwZS5vcmeb7jwaAAAAAElFTkSuQmCC"
                          />
                        </div>

                        <div className="border flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto  dark:border-gray-900">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            x="0px"
                            y="0px"
                            className="p-2 w-12 h-12"
                            viewBox="0 0 48 48"
                          >
                            <linearGradient
                              id="k8yl7~hDat~FaoWq8WjN6a_VLKafOkk3sBX_gr1"
                              x1="-1254.397"
                              x2="-1261.911"
                              y1="877.268"
                              y2="899.466"
                              gradientTransform="translate(1981.75 -1362.063) scale(1.5625)"
                              gradientUnits="userSpaceOnUse"
                            >
                              <stop offset="0" stopColor="#114a8b"></stop>
                              <stop offset="1" stopColor="#0669bc"></stop>
                            </linearGradient>
                            <path
                              fill="url(#k8yl7~hDat~FaoWq8WjN6a_VLKafOkk3sBX_gr1)"
                              d="M17.634,6h11.305L17.203,40.773c-0.247,0.733-0.934,1.226-1.708,1.226H6.697 c-0.994,0-1.8-0.806-1.8-1.8c0-0.196,0.032-0.39,0.094-0.576L15.926,7.227C16.173,6.494,16.86,6,17.634,6L17.634,6z"
                            ></path>
                            <path
                              fill="#0078d4"
                              d="M34.062,29.324H16.135c-0.458-0.001-0.83,0.371-0.831,0.829c0,0.231,0.095,0.451,0.264,0.608 l11.52,10.752C27.423,41.826,27.865,42,28.324,42h10.151L34.062,29.324z"
                            ></path>
                            <linearGradient
                              id="k8yl7~hDat~FaoWq8WjN6b_VLKafOkk3sBX_gr2"
                              x1="-1252.05"
                              x2="-1253.788"
                              y1="887.612"
                              y2="888.2"
                              gradientTransform="translate(1981.75 -1362.063) scale(1.5625)"
                              gradientUnits="userSpaceOnUse"
                            >
                              <stop offset="0" stopOpacity=".3"></stop>
                              <stop offset=".071" stopOpacity=".2"></stop>
                              <stop offset=".321" stopOpacity=".1"></stop>
                              <stop offset=".623" stopOpacity=".05"></stop>
                              <stop offset="1" stopOpacity="0"></stop>
                            </linearGradient>
                            <path
                              fill="url(#k8yl7~hDat~FaoWq8WjN6b_VLKafOkk3sBX_gr2)"
                              d="M17.634,6c-0.783-0.003-1.476,0.504-1.712,1.25L5.005,39.595 c-0.335,0.934,0.151,1.964,1.085,2.299C6.286,41.964,6.493,42,6.702,42h9.026c0.684-0.122,1.25-0.603,1.481-1.259l2.177-6.416 l7.776,7.253c0.326,0.27,0.735,0.419,1.158,0.422h10.114l-4.436-12.676l-12.931,0.003L28.98,6H17.634z"
                            ></path>
                            <linearGradient
                              id="k8yl7~hDat~FaoWq8WjN6c_VLKafOkk3sBX_gr3"
                              x1="-1252.952"
                              x2="-1244.704"
                              y1="876.6"
                              y2="898.575"
                              gradientTransform="translate(1981.75 -1362.063) scale(1.5625)"
                              gradientUnits="userSpaceOnUse"
                            >
                              <stop offset="0" stopColor="#3ccbf4"></stop>
                              <stop offset="1" stopColor="#2892df"></stop>
                            </linearGradient>
                            <path
                              fill="url(#k8yl7~hDat~FaoWq8WjN6c_VLKafOkk3sBX_gr3)"
                              d="M32.074,7.225C31.827,6.493,31.141,6,30.368,6h-12.6c0.772,0,1.459,0.493,1.705,1.224 l10.935,32.399c0.318,0.942-0.188,1.963-1.13,2.281C29.093,41.968,28.899,42,28.703,42h12.6c0.994,0,1.8-0.806,1.8-1.801 c0-0.196-0.032-0.39-0.095-0.575L32.074,7.225z"
                            ></path>
                          </svg>
                        </div>
                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[62px] lg:max-w-[62px] h-auto  dark:border-gray-900 flex items-center justify-center">
                          <svg
                            className="p-2 w-12 h-12 dark:fill-white fill-black"
                            viewBox="0 0 876 876"
                          >
                            <path d="M468 292H528V584H468V292Z" />
                            <path d="M348 292H408V584H348V292Z" />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-base leading-tight lg:text-xl">
                          Bring your own voices
                        </div>
                        <div className="text-base dark:text-gray-500">
                          We support built-in: ElevenLabs, PlayHT, LMNT,
                          Deepgram, Cartesia, Rime, OpenAI, Azure.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="md:col-span-3">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px border-r dark:border-gray-900 group">
                  <div className="relative h-full group w-full p-6 lg:p-8 lg:pt-8">
                    <div className="relative flex flex-col items-start gap-6  transition-all duration-300">
                      <div className="pt-3.5 pb-1 md:py-1 justify-center gap-2.5 lg:gap-3 items-center flex flex-1 flex-wrap group-hover:grayscale-0 grayscale">
                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto  dark:border-gray-900">
                          <svg
                            viewBox="-1 -.1 949.1 959.8"
                            xmlns="http://www.w3.org/2000/svg"
                            className="p-2 w-12 h-12"
                          >
                            <path d="m925.8 456.3c10.4 23.2 17 48 19.7 73.3 2.6 25.3 1.3 50.9-4.1 75.8-5.3 24.9-14.5 48.8-27.3 70.8-8.4 14.7-18.3 28.5-29.7 41.2-11.3 12.6-23.9 24-37.6 34-13.8 10-28.5 18.4-44.1 25.3-15.5 6.8-31.7 12-48.3 15.4-7.8 24.2-19.4 47.1-34.4 67.7-14.9 20.6-33 38.7-53.6 53.6-20.6 15-43.4 26.6-67.6 34.4-24.2 7.9-49.5 11.8-75 11.8-16.9.1-33.9-1.7-50.5-5.1-16.5-3.5-32.7-8.8-48.2-15.7s-30.2-15.5-43.9-25.5c-13.6-10-26.2-21.5-37.4-34.2-25 5.4-50.6 6.7-75.9 4.1-25.3-2.7-50.1-9.3-73.4-19.7-23.2-10.3-44.7-24.3-63.6-41.4s-35-37.1-47.7-59.1c-8.5-14.7-15.5-30.2-20.8-46.3s-8.8-32.7-10.6-49.6c-1.8-16.8-1.7-33.8.1-50.7 1.8-16.8 5.5-33.4 10.8-49.5-17-18.9-31-40.4-41.4-63.6-10.3-23.3-17-48-19.6-73.3-2.7-25.3-1.3-50.9 4-75.8s14.5-48.8 27.3-70.8c8.4-14.7 18.3-28.6 29.6-41.2s24-24 37.7-34 28.5-18.5 44-25.3c15.6-6.9 31.8-12 48.4-15.4 7.8-24.3 19.4-47.1 34.3-67.7 15-20.6 33.1-38.7 53.7-53.7 20.6-14.9 43.4-26.5 67.6-34.4 24.2-7.8 49.5-11.8 75-11.7 16.9-.1 33.9 1.6 50.5 5.1s32.8 8.7 48.3 15.6c15.5 7 30.2 15.5 43.9 25.5 13.7 10.1 26.3 21.5 37.5 34.2 24.9-5.3 50.5-6.6 75.8-4s50 9.3 73.3 19.6c23.2 10.4 44.7 24.3 63.6 41.4 18.9 17 35 36.9 47.7 59 8.5 14.6 15.5 30.1 20.8 46.3 5.3 16.1 8.9 32.7 10.6 49.6 1.8 16.9 1.8 33.9-.1 50.8-1.8 16.9-5.5 33.5-10.8 49.6 17.1 18.9 31 40.3 41.4 63.6zm-333.2 426.9c21.8-9 41.6-22.3 58.3-39s30-36.5 39-58.4c9-21.8 13.7-45.2 13.7-68.8v-223q-.1-.3-.2-.7-.1-.3-.3-.6-.2-.3-.5-.5-.3-.3-.6-.4l-80.7-46.6v269.4c0 2.7-.4 5.5-1.1 8.1-.7 2.7-1.7 5.2-3.1 7.6s-3 4.6-5 6.5a32.1 32.1 0 0 1 -6.5 5l-191.1 110.3c-1.6 1-4.3 2.4-5.7 3.2 7.9 6.7 16.5 12.6 25.5 17.8 9.1 5.2 18.5 9.6 28.3 13.2 9.8 3.5 19.9 6.2 30.1 8 10.3 1.8 20.7 2.7 31.1 2.7 23.6 0 47-4.7 68.8-13.8zm-455.1-151.4c11.9 20.5 27.6 38.3 46.3 52.7 18.8 14.4 40.1 24.9 62.9 31s46.6 7.7 70 4.6 45.9-10.7 66.4-22.5l193.2-111.5.5-.5q.2-.2.3-.6.2-.3.3-.6v-94l-233.2 134.9c-2.4 1.4-4.9 2.4-7.5 3.2-2.7.7-5.4 1-8.2 1-2.7 0-5.4-.3-8.1-1-2.6-.8-5.2-1.8-7.6-3.2l-191.1-110.4c-1.7-1-4.2-2.5-5.6-3.4-1.8 10.3-2.7 20.7-2.7 31.1s1 20.8 2.8 31.1c1.8 10.2 4.6 20.3 8.1 30.1 3.6 9.8 8 19.2 13.2 28.2zm-50.2-417c-11.8 20.5-19.4 43.1-22.5 66.5s-1.5 47.1 4.6 70c6.1 22.8 16.6 44.1 31 62.9 14.4 18.7 32.3 34.4 52.7 46.2l193.1 111.6q.3.1.7.2h.7q.4 0 .7-.2.3-.1.6-.3l81-46.8-233.2-134.6c-2.3-1.4-4.5-3.1-6.5-5a32.1 32.1 0 0 1 -5-6.5c-1.3-2.4-2.4-4.9-3.1-7.6-.7-2.6-1.1-5.3-1-8.1v-227.1c-9.8 3.6-19.3 8-28.3 13.2-9 5.3-17.5 11.3-25.5 18-7.9 6.7-15.3 14.1-22 22.1-6.7 7.9-12.6 16.5-17.8 25.5zm663.3 154.4c2.4 1.4 4.6 3 6.6 5 1.9 1.9 3.6 4.1 5 6.5 1.3 2.4 2.4 5 3.1 7.6.6 2.7 1 5.4.9 8.2v227.1c32.1-11.8 60.1-32.5 80.8-59.7 20.8-27.2 33.3-59.7 36.2-93.7s-3.9-68.2-19.7-98.5-39.9-55.5-69.5-72.5l-193.1-111.6q-.3-.1-.7-.2h-.7q-.3.1-.7.2-.3.1-.6.3l-80.6 46.6 233.2 134.7zm80.5-121h-.1v.1zm-.1-.1c5.8-33.6 1.9-68.2-11.3-99.7-13.1-31.5-35-58.6-63-78.2-28-19.5-61-30.7-95.1-32.2-34.2-1.4-68 6.9-97.6 23.9l-193.1 111.5q-.3.2-.5.5l-.4.6q-.1.3-.2.7-.1.3-.1.7v93.2l233.2-134.7c2.4-1.4 5-2.4 7.6-3.2 2.7-.7 5.4-1 8.1-1 2.8 0 5.5.3 8.2 1 2.6.8 5.1 1.8 7.5 3.2l191.1 110.4c1.7 1 4.2 2.4 5.6 3.3zm-505.3-103.2c0-2.7.4-5.4 1.1-8.1.7-2.6 1.7-5.2 3.1-7.6 1.4-2.3 3-4.5 5-6.5 1.9-1.9 4.1-3.6 6.5-4.9l191.1-110.3c1.8-1.1 4.3-2.5 5.7-3.2-26.2-21.9-58.2-35.9-92.1-40.2-33.9-4.4-68.3 1-99.2 15.5-31 14.5-57.2 37.6-75.5 66.4-18.3 28.9-28 62.3-28 96.5v223q.1.4.2.7.1.3.3.6.2.3.5.6.2.2.6.4l80.7 46.6zm43.8 294.7 103.9 60 103.9-60v-119.9l-103.8-60-103.9 60z" />
                          </svg>
                        </div>

                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto  dark:border-gray-900">
                          <svg
                            viewBox="0 0 256 176"
                            version="1.1"
                            className="p-2 w-12 h-12 dark:invert"
                            xmlns="http://www.w3.org/2000/svg"
                            preserveAspectRatio="xMidYMid"
                          >
                            <g fill="#181818">
                              <path d="M147.486878,0 C147.486878,0 217.568251,175.780074 217.568251,175.780074 C217.568251,175.780074 256,175.780074 256,175.780074 C256,175.780074 185.918621,0 185.918621,0 C185.918621,0 147.486878,0 147.486878,0 C147.486878,0 147.486878,0 147.486878,0 Z" />
                              <path d="M66.1828124,106.221191 C66.1828124,106.221191 90.1624677,44.4471185 90.1624677,44.4471185 C90.1624677,44.4471185 114.142128,106.221191 114.142128,106.221191 C114.142128,106.221191 66.1828124,106.221191 66.1828124,106.221191 C66.1828124,106.221191 66.1828124,106.221191 66.1828124,106.221191 Z M70.0705318,0 C70.0705318,0 0,175.780074 0,175.780074 C0,175.780074 39.179211,175.780074 39.179211,175.780074 C39.179211,175.780074 53.5097704,138.86606 53.5097704,138.86606 C53.5097704,138.86606 126.817544,138.86606 126.817544,138.86606 C126.817544,138.86606 141.145724,175.780074 141.145724,175.780074 C141.145724,175.780074 180.324935,175.780074 180.324935,175.780074 C180.324935,175.780074 110.254409,0 110.254409,0 C110.254409,0 70.0705318,0 70.0705318,0 C70.0705318,0 70.0705318,0 70.0705318,0 Z" />
                            </g>
                          </svg>
                        </div>
                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto  dark:border-gray-900">
                          <svg
                            className="p-2 w-12 h-12"
                            viewBox="0 0 256 233"
                            version="1.1"
                            xmlns="http://www.w3.org/2000/svg"
                            preserveAspectRatio="xMidYMid"
                          >
                            <g>
                              <rect
                                fill="#000000"
                                x="186.181818"
                                y={0}
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#F7D046"
                                x="209.454545"
                                y={0}
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x={0}
                                y={0}
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x={0}
                                y="46.5454545"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x={0}
                                y="93.0909091"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x={0}
                                y="139.636364"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x={0}
                                y="186.181818"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#F7D046"
                                x="23.2727273"
                                y={0}
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#F2A73B"
                                x="209.454545"
                                y="46.5454545"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#F2A73B"
                                x="23.2727273"
                                y="46.5454545"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x="139.636364"
                                y="46.5454545"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#F2A73B"
                                x="162.909091"
                                y="46.5454545"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#F2A73B"
                                x="69.8181818"
                                y="46.5454545"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EE792F"
                                x="116.363636"
                                y="93.0909091"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EE792F"
                                x="162.909091"
                                y="93.0909091"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EE792F"
                                x="69.8181818"
                                y="93.0909091"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x="93.0909091"
                                y="139.636364"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EB5829"
                                x="116.363636"
                                y="139.636364"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EE792F"
                                x="209.454545"
                                y="93.0909091"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EE792F"
                                x="23.2727273"
                                y="93.0909091"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x="186.181818"
                                y="139.636364"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EB5829"
                                x="209.454545"
                                y="139.636364"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#000000"
                                x="186.181818"
                                y="186.181818"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EB5829"
                                x="23.2727273"
                                y="139.636364"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EA3326"
                                x="209.454545"
                                y="186.181818"
                                width="46.5454545"
                                height="46.5454545"
                              />
                              <rect
                                fill="#EA3326"
                                x="23.2727273"
                                y="186.181818"
                                width="46.5454545"
                                height="46.5454545"
                              />
                            </g>
                          </svg>
                        </div>

                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto  dark:border-gray-900">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 24 24"
                            className="p-2 w-12 h-12"
                          >
                            <path
                              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                              fill="#4285F4"
                            />
                            <path
                              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                              fill="#34A853"
                            />
                            <path
                              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                              fill="#FBBC05"
                            />
                            <path
                              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                              fill="#EA4335"
                            />
                            <path d="M1 1h22v22H1z" fill="none" />
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg leading-tight lg:text-xl">
                          Bring your own models.
                        </div>
                        <div className="text-base dark:text-gray-500">
                          We support built-in: OpenAI, Groq, Mistral,
                          OpenRouter, Together, Anyscale or bring your own
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="md:col-span-3">
                <div className="relative h-full w-full overflow-hidden transition-all duration-300 p-px  dark:border-gray-900 group">
                  <div className="relative h-full group w-full p-6 lg:p-8 lg:pt-8">
                    <div className="relative flex flex-col items-start gap-6 rounded-[2px] transition-all duration-300">
                      <div className="pt-3.5 pb-1 md:py-1 justify-center gap-2.5 lg:gap-3 items-center flex flex-1 flex-wrap group-hover:grayscale-0 grayscale">
                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto rounded-[2px] dark:border-gray-900">
                          <svg
                            role="img"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                            className="p-2 w-12 h-12"
                          >
                            <path d="M9.279 11.617l-4.54-10.07H0l6.797 15.296a.084.084 0 0 0 .153 0zm9.898-10.07s-6.148 13.868-6.917 15.565c-1.838 4.056-3.2 5.07-4.588 5.289a.026.026 0 0 0 .004.052h4.34c1.911 0 3.219-1.285 5.06-5.341C17.72 15.694 24 1.547 24 1.547z" />
                          </svg>
                        </div>

                        <div className="border bg-opacity-8 flex-1 min-w-[20%] max-w-[66px] lg:max-w-[62px] h-auto rounded-[2px] dark:border-gray-900">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            xmlnsXlink="http://www.w3.org/1999/xlink"
                            aria-label="Twilio"
                            role="img"
                            className="p-2 w-12 h-12"
                            viewBox="0 0 512 512"
                          >
                            <g fill="#f22f46">
                              <circle cx={256} cy={256} r={256} />
                              <circle cx={256} cy={256} fill="#fff" r={188} />
                              <circle cx={193} cy={193} r={53} id="c" />
                              <use xlinkHref="#c" x={126} />
                              <use xlinkHref="#c" y={126} />
                              <use xlinkHref="#c" x={126} y={126} />
                            </g>
                          </svg>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="text-lg leading-tight lg:text-xl">
                          Connect to telephony network
                        </div>
                        <div className="text-base dark:text-gray-500">
                          Allow users to dial into sessions or make programmatic
                          calls out to phones.
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
        {/* section */}
        <section className="relative mx-auto w-full max-w-7xl px-4 sm:px-8 pt-16 lg:pt-32 dark:text-gray-400">
          <div className="space-between left-0 flex w-full flex-col justify-between gap-8 lg:flex-row">
            <div className="max-w-3xl space-y-2  lg:w-1/2">
              <div className=" text-base uppercase tracking-wider">
                global scale
              </div>
              <h2 className="z-20 text-3xl leading-[1.15em] lg:text-5xl tracking-tight lg:leading-[1.15em] ">
                <span className="">The{/* */} </span>
                <span className="">backbone{/* */} </span>
                <span className="">of{/* */} </span>
                <span className="">the{/* */} </span>
                <span className="">
                  realtime conversation experience{/* */}{' '}
                </span>
              </h2>
              <p className="pt-2 text-lg lg:text-xl dark:text-gray-600">
                Our agents are built to ensure a smooth and reliable customer
                experience, with every call delivered to completion without
                interruptions or errors.
              </p>
            </div>
            <div className="flex pb-8 lg:w-1/2 lg:items-end lg:justify-end">
              <div className="flex gap-12 text-fg1">
                <div className="flex flex-col gap-2">
                  <div className="text-primary dark:text-blue-500">Uptime</div>
                  <div className="text-glow-white text-2xl text-fg0 lg:text-3xl lg:">
                    99.99%
                  </div>
                </div>
                <div className="flex flex-col gap-2">
                  <div className="text-primary dark:text-blue-500">Latency</div>
                  <div className="text-glow-white text-2xl text-fg0 lg:text-3xl lg:">
                    &lt;100ms
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/*  */}
        <div className="max-w-7xl mx-auto">
          <div className="py-20 border-r border-l dark:border-gray-900"></div>
        </div>
      </main>
      <footer className="border-t dark:border-gray-900">
        <div className="max-w-7xl mx-auto">
          <div className="relative px-4 sm:px-8 border-x dark:border-gray-900 pt-16 dark:text-gray-500">
            <div className="xs:flex-row flex flex-col justify-between gap-4">
              <div className="w-fit flex shrink-0 text-blue-600 dark:text-blue-500 items-center">
                <RapidaIcon className="h-8 w-8 shrink-0" />
                <RapidaTextIcon className="h-6 shrink-0 ml-1" />
              </div>
              <p className="text-sm">
                Build voice assistants, <br />
                unleashing the true potential of businesses with the power of
                generative ai
              </p>
              <div className="flex text-black">
                <a
                  href="https://x.com/rapidaai"
                  target="_blank"
                  aria-label="Twitter"
                  className="hover:border-blue-600 last:border-r-border group flex h-9 w-9 items-center justify-center border border-r-transparent transition-colors undefined"
                  rel="noreferrer"
                >
                  <svg
                    width={18}
                    height={18}
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 512 512"
                    className="fill-black dark:fill-white"
                  >
                    <path d="M389.2 48h70.6L305.6 224.2 487 464H345L233.7 318.6 106.5 464H35.8L200.7 275.5 26.8 48H172.4L272.9 180.9 389.2 48zM364.4 421.8h39.1L151.1 88h-42L364.4 421.8z" />
                  </svg>
                </a>
                <a
                  href="https://www.linkedin.com/company/rapida-ai"
                  aria-label="LinkedIn"
                  target="_blank"
                  className="hover:border-blue-600 last:border-r-border group flex h-9 w-9 items-center justify-center border border-r-transparent transition-colors"
                  rel="noreferrer"
                >
                  <svg
                    width={18}
                    height={18}
                    viewBox="0 0 13 13"
                    fill="none"
                    className="fill-black dark:fill-white"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M11.875 0.125C12.3398 0.125 12.75 0.535156 12.75 1.02734V11.5C12.75 11.9922 12.3398 12.375 11.875 12.375H1.34766C0.882812 12.375 0.5 11.9922 0.5 11.5V1.02734C0.5 0.535156 0.882812 0.125 1.34766 0.125H11.875ZM4.19141 10.625V4.80078H2.38672V10.625H4.19141ZM3.28906 3.98047C3.86328 3.98047 4.32812 3.51562 4.32812 2.94141C4.32812 2.36719 3.86328 1.875 3.28906 1.875C2.6875 1.875 2.22266 2.36719 2.22266 2.94141C2.22266 3.51562 2.6875 3.98047 3.28906 3.98047ZM11 10.625V7.42578C11 5.86719 10.6445 4.63672 8.8125 4.63672C7.9375 4.63672 7.33594 5.12891 7.08984 5.59375H7.0625V4.80078H5.33984V10.625H7.14453V7.75391C7.14453 6.98828 7.28125 6.25 8.23828 6.25C9.16797 6.25 9.16797 7.125 9.16797 7.78125V10.625H11Z" />
                  </svg>
                </a>
                <a
                  href="https://www.youtube.com/@RapidaAI"
                  aria-label="YouTube"
                  target="_blank"
                  className="hover:border-blue-600 last:border-r-border group flex h-9 w-9 items-center justify-center border border-r-transparent transition-colors undefined"
                  rel="noreferrer"
                >
                  <svg
                    width={18}
                    height={18}
                    viewBox="0 0 16 11"
                    fill="none"
                    className="fill-black dark:fill-white"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M15.0117 1.66797C15.3398 2.81641 15.3398 5.27734 15.3398 5.27734C15.3398 5.27734 15.3398 7.71094 15.0117 8.88672C14.8477 9.54297 14.3281 10.0352 13.6992 10.1992C12.5234 10.5 7.875 10.5 7.875 10.5C7.875 10.5 3.19922 10.5 2.02344 10.1992C1.39453 10.0352 0.875 9.54297 0.710938 8.88672C0.382812 7.71094 0.382812 5.27734 0.382812 5.27734C0.382812 5.27734 0.382812 2.81641 0.710938 1.66797C0.875 1.01172 1.39453 0.492188 2.02344 0.328125C3.19922 0 7.875 0 7.875 0C7.875 0 12.5234 0 13.6992 0.328125C14.3281 0.492188 14.8477 1.01172 15.0117 1.66797ZM6.34375 7.49219L10.2266 5.27734L6.34375 3.0625V7.49219Z" />
                  </svg>
                </a>
                <a
                  href="https://github.com/rapidaai"
                  aria-label="GitHub"
                  className="hover:border-blue-600 border-r group flex h-9 w-9 items-center justify-center border transition-colors"
                >
                  <svg
                    viewBox="0 0 14 14"
                    fill="none"
                    width={18}
                    height={18}
                    className="fill-black dark:fill-white"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M4.51172 11.1328C4.51172 11.0781 4.45703 11.0234 4.375 11.0234C4.29297 11.0234 4.23828 11.0781 4.23828 11.1328C4.23828 11.1875 4.29297 11.2422 4.375 11.2148C4.45703 11.2148 4.51172 11.1875 4.51172 11.1328ZM3.66406 10.9961C3.69141 10.9414 3.77344 10.9141 3.85547 10.9414C3.9375 10.9688 3.96484 11.0234 3.96484 11.0781C3.9375 11.1328 3.85547 11.1602 3.80078 11.1328C3.71875 11.1328 3.66406 11.0508 3.66406 10.9961ZM4.89453 10.9688C4.94922 10.9414 5.03125 10.9961 5.03125 11.0508C5.05859 11.1055 5.00391 11.1328 4.92188 11.1602C4.83984 11.1875 4.75781 11.1602 4.75781 11.1055C4.75781 11.0234 4.8125 10.9688 4.89453 10.9688ZM6.67188 0.46875C10.4727 0.46875 13.5625 3.36719 13.5625 7.14062C13.5625 10.1758 11.7031 12.7734 8.96875 13.6758C8.61328 13.7578 8.47656 13.5391 8.47656 13.3477C8.47656 13.1289 8.50391 11.9805 8.50391 11.0781C8.50391 10.4219 8.28516 10.0117 8.03906 9.79297C9.57031 9.62891 11.1836 9.41016 11.1836 6.78516C11.1836 6.01953 10.9102 5.66406 10.4727 5.17188C10.5273 4.98047 10.7734 4.26953 10.3906 3.3125C9.81641 3.12109 8.50391 4.05078 8.50391 4.05078C7.95703 3.88672 7.38281 3.83203 6.78125 3.83203C6.20703 3.83203 5.63281 3.88672 5.08594 4.05078C5.08594 4.05078 3.74609 3.14844 3.19922 3.3125C2.81641 4.26953 3.03516 4.98047 3.11719 5.17188C2.67969 5.66406 2.46094 6.01953 2.46094 6.78516C2.46094 9.41016 4.01953 9.62891 5.55078 9.79297C5.33203 9.98438 5.16797 10.2852 5.11328 10.7227C4.70312 10.9141 3.71875 11.2148 3.11719 10.1484C2.73438 9.49219 2.05078 9.4375 2.05078 9.4375C1.39453 9.4375 2.02344 9.875 2.02344 9.875C2.46094 10.0664 2.76172 10.8594 2.76172 10.8594C3.17188 12.0898 5.08594 11.6797 5.08594 11.6797C5.08594 12.2539 5.08594 13.1836 5.08594 13.375C5.08594 13.5391 4.97656 13.7578 4.62109 13.7031C1.88672 12.7734 0 10.1758 0 7.14062C0 3.36719 2.89844 0.46875 6.67188 0.46875ZM2.65234 9.90234C2.67969 9.875 2.73438 9.90234 2.78906 9.92969C2.84375 9.98438 2.84375 10.0664 2.81641 10.0938C2.76172 10.1211 2.70703 10.0938 2.65234 10.0664C2.625 10.0117 2.59766 9.92969 2.65234 9.90234ZM2.35156 9.68359C2.37891 9.65625 2.40625 9.65625 2.46094 9.68359C2.51562 9.71094 2.54297 9.73828 2.54297 9.76562C2.51562 9.82031 2.46094 9.82031 2.40625 9.79297C2.35156 9.76562 2.32422 9.73828 2.35156 9.68359ZM3.22656 10.668C3.28125 10.6133 3.36328 10.6406 3.41797 10.6953C3.47266 10.75 3.47266 10.832 3.44531 10.8594C3.41797 10.9141 3.33594 10.8867 3.28125 10.832C3.19922 10.7773 3.19922 10.6953 3.22656 10.668ZM2.92578 10.2578C2.98047 10.2305 3.03516 10.2578 3.08984 10.3125C3.11719 10.3672 3.11719 10.4492 3.08984 10.4766C3.03516 10.5039 2.98047 10.4766 2.92578 10.4219C2.87109 10.3672 2.87109 10.2852 2.92578 10.2578Z" />
                  </svg>
                </a>
              </div>
            </div>
            <div className="pb-6" />
          </div>
        </div>
        <div className="border-t dark:border-gray-900 dark:text-gray-500">
          <div className="max-w-7xl mx-auto">
            <div className="relative px-4 sm:px-8 border-x dark:border-gray-900 py-4">
              <div className="text-tertiary flex flex-wrap justify-between gap-x-10 gap-y-1 text-sm">
                <p>© RapidaAI | Singapore, India</p>
                <p>Rapida is a registered trademark of Rapida Pte. Ltd.</p>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </SmoothScroll>
  );
};
