import { ChevronRight, Rocket } from 'lucide-react';

export const Pricing = () => (
  <div
    id="pricing"
    className="border-y mt-40 grid grid-cols-1 gap-x-10 lg:grid-cols-2 lg:grid-rows-[1fr_auto]"
  >
    <div className="col-start-1 row-start-1 flex flex-col gap-5 lg:border-r  lg:pb-10">
      <div className="grid grid-cols-1 gap-y-2 px-4 py-2 border-b sm:px-2 lg:border-b/half">
        <p className="text-[2.5rem]/none font-medium text-pretty">
          From pilot to production
        </p>
      </div>
      <div className="px-4 py-2 max-lg:line-y sm:px-2 lg:line-y/half">
        <p className="max-w-2xl text-sm/7 text-gray-600 dark:text-gray-400">
          Whether you’re testing your first use case or deploying across
          regions, Rapida gives you a clear path from prototype to full-scale
          rollout.
        </p>
      </div>
    </div>
    <div className="max-lg:row-start-3 lg:row-start-2 bg-gray-950/5 py-[calc(--spacing(2)+1px)] dark:bg-white/5 max-lg:-mx-px max-lg:px-[calc(--spacing(2)+1px)] max-lg:pt-0 lg:line-t/half lg:-ml-px lg:border-r lg:border-t lg:pr-2 lg:pl-[calc(--spacing(2)+1px)]">
      <div className="grid grid-cols-1 gap-y-6 rounded-2xl bg-white p-6 sm:rounded-4xl sm:p-10 dark:bg-transparent dark:outline dark:outline-white/10">
        <div className="flex flex-wrap items-center justify-between gap-6">
          <div className="flex items-center gap-x-4">
            <p className="text-5xl font-light *:font-medium">
              US$<span>599</span>
            </p>
            <div>
              <p className="text-sm/6 font-semibold">one time payment</p>
              <p className="text-sm/6 text-gray-600 dark:text-gray-400">
                plus local taxes
              </p>
            </div>
          </div>
          <a
            className="pl-4 px-2 py-2 gap-2 inline-flex justify-center rounded-full text-base font-medium focus-visible:outline-2 focus-visible:outline-offset-2 bg-blue-600 text-white"
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
        </div>
        <p className="max-w-2xl text-sm/7 text-gray-600 dark:text-gray-400">
          For startups and product teams exploring voice ai.
        </p>
        <hr className="border-gray-950/5 dark:border-white/10" />
        <div className="@container">
          <ul
            className="group grid grid-cols-1 gap-x-10 gap-y-6 text-sm/7 text-gray-600 data-dark:text-gray-300 @3xl:grid-cols-2 dark:text-gray-400"
            role="list"
          >
            {[
              {
                title: '50,000 voice minutes',
                description:
                  'perfect for testing, prototyping, and running pilot deployments without scaling costs.',
              },

              {
                title: 'Full SDK & orchestration suite',
                description:
                  'complete access to all development tools and APIs for maximum flexibility.',
              },
              {
                title: 'Continuous updates included',
                description:
                  'get new features, integrations, and improvements automatically at no extra cost.',
              },
              {
                title: 'Seamless upgrade path',
                description:
                  'grow into Enterprise without migration hassle or data loss.',
              },
            ].map(item => (
              <li
                key={item.title}
                className="grid max-w-2xl grid-cols-[auto_1fr] gap-6"
              >
                <svg
                  aria-hidden="true"
                  viewBox="0 0 22 22"
                  className="h-7 w-5.5"
                >
                  <path
                    className="fill-blue-600"
                    d="M22 11c0 6.075-4.925 11-11 11S0 17.075 0 11 4.925 0 11 0s11 4.925 11 11Z"
                  />
                  <path
                    className="fill-blue-600"
                    d="M11 21c5.523 0 10-4.477 10-10S16.523 1 11 1 1 5.477 1 11s4.477 10 10 10Zm0 1c6.075 0 11-4.925 11-11S17.075 0 11 0 0 4.925 0 11s4.925 11 11 11Z"
                    clipRule="evenodd"
                    fillRule="evenodd"
                  />
                  <path
                    className="fill-white"
                    d="m14.684 7.82-4.079 6.992L7.293 11.5 8 10.793l2.395 2.395 3.425-5.872.864.504Z"
                    clipRule="evenodd"
                    fillRule="evenodd"
                  />
                </svg>
                <p>
                  <strong className="font-semibold text-gray-950 group-data-dark:text-white dark:text-white">
                    {item.title}
                  </strong>{' '}
                  — {item.description}
                </p>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
    <div className="lg:col-start-2 lg:row-span-2 flex bg-gray-950/5 py-[calc(--spacing(2)+1px)] dark:bg-white/5 max-lg:-mx-px max-lg:px-[calc(--spacing(2)+1px)] lg:-mr-px lg:border-l  lg:pr-[calc(--spacing(2)+1px)] lg:pl-2">
      <div className="grid grid-cols-1 gap-y-6 rounded-2xl bg-gray-950 p-6 sm:rounded-4xl sm:p-10 dark:bg-gray-950 dark:inset-ring dark:inset-ring-white/15">
        <p className="text-gray-300 font-mono text-[0.8125rem]/6 font-medium tracking-widest text-pretty uppercase">
          Get everything with <span className="sr-only">Rapida</span>
        </p>

        <div>
          <div className="flex flex-wrap items-center justify-between gap-6">
            <div className="flex items-center gap-x-4">
              <p className="text-6xl font-light text-white *:font-medium">
                Custom
              </p>
              <div>
                <p className="text-sm/6 font-semibold text-white">
                  Usage-based
                </p>
                <p className="text-sm/6 text-gray-300">
                  scaling with transparent pricing
                </p>
              </div>
            </div>
            <a
              className="pl-4 px-2 py-2 gap-2 inline-flex justify-center rounded-full text-base font-medium focus-visible:outline-2 focus-visible:outline-offset-2 bg-blue-600 text-white"
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
          </div>
        </div>

        <hr className="border-t-white/10" />
        <div className="@container">
          <ul
            className="group grid grid-cols-1 gap-x-10 gap-y-6 text-sm/7 text-gray-600 data-dark:text-gray-300 @3xl:grid-cols-2 dark:text-gray-400"
            role="list"
          >
            {[
              {
                title: 'Unlimited voice minutes with flexible scaling',
                description:
                  'buy once and pay only for what you use with transparent, usage-based pricing.',
              },
              {
                title: 'Multi-region deployment with on-premises option',
                description:
                  'deploy globally, in private clouds, or on your own infrastructure with complete control.',
              },
              {
                title: 'Unlimited projects and production environments',
                description:
                  'scale across as many use cases, regions, and deployments as you need.',
              },
              {
                title: 'Dedicated Deployment engineer and 24x7 support',
                description:
                  'white-glove implementation, architecture guidance, and mission-critical SLA support.',
              },
            ].map(item => (
              <li
                key={item.title}
                className="grid max-w-2xl grid-cols-[auto_1fr] gap-6"
              >
                <svg
                  aria-hidden="true"
                  viewBox="0 0 22 22"
                  className="h-7 w-5.5"
                >
                  <path
                    className="fill-blue-600"
                    d="M22 11c0 6.075-4.925 11-11 11S0 17.075 0 11 4.925 0 11 0s11 4.925 11 11Z"
                  />
                  <path
                    className="fill-blue-600"
                    d="M11 21c5.523 0 10-4.477 10-10S16.523 1 11 1 1 5.477 1 11s4.477 10 10 10Zm0 1c6.075 0 11-4.925 11-11S17.075 0 11 0 0 4.925 0 11s4.925 11 11 11Z"
                    clipRule="evenodd"
                    fillRule="evenodd"
                  />
                  <path
                    className="fill-white"
                    d="m14.684 7.82-4.079 6.992L7.293 11.5 8 10.793l2.395 2.395 3.425-5.872.864.504Z"
                    clipRule="evenodd"
                    fillRule="evenodd"
                  />
                </svg>
                <p>
                  <strong className="font-semibold text-white">
                    {item.title}
                  </strong>{' '}
                  —{' '}
                  <span className="text-gray-300 dark:text-gray-300">
                    {item.description}
                  </span>
                </p>
              </li>
            ))}
          </ul>
        </div>
        <p className="mt-4 rounded-xl bg-white/10 p-6 text-sm/7 text-gray-300 sm:rounded-2xl dark:bg-white/5">
          <Rocket className="mr-2 inline-block h-7 w-4" />
          <strong className="font-semibold text-white">
            Get live in days, not months
          </strong>{' '}
          — Jump-start your deployment with pre-built templates for common use
          cases and ready-made integrations.
        </p>
      </div>
    </div>
  </div>
);
