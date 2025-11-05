import { useState } from 'react';

import {
  AdjustmentsHorizontalIcon,
  ArrowsRightLeftIcon,
  CloudArrowUpIcon,
  WrenchScrewdriverIcon,
} from '@heroicons/react/20/solid';

export const AgentStudio = () => {
  const [open, setOpen] = useState(1); // Initial panel open by default

  return (
    <div className="border sm:rounded-2xl bg-white dark:bg-gray-950 flex flex-col lg:flex-row h-fit sm:h-[500px]  w-full mx-auto overflow-hidden">
      <BuildinAgentPanel open={open} setOpen={setOpen} />
      <WebsocketPanel open={open} setOpen={setOpen} />
      <AgentKitPanel open={open} setOpen={setOpen} />
    </div>
  );
};

const BuildinAgentPanel = ({ open, setOpen }) => (
  <Panel
    open={open}
    setOpen={setOpen}
    id={1}
    content={<AgentBuilder />}
    title="Agent Builder"
  />
);

const WebsocketPanel = ({ open, setOpen }) => (
  <Panel
    open={open}
    setOpen={setOpen}
    id={2}
    content={<Websocket />}
    title="Connect with Websocket"
  />
);

const AgentKitPanel = ({ open, setOpen }) => (
  <Panel
    open={open}
    setOpen={setOpen}
    id={3}
    content={<AgentKit />}
    title="Connect with AgentKit"
  />
);

const Panel = ({ open, setOpen, id, title, content }) => {
  const isOpen = open === id;
  return (
    <>
      <button
        className="bg-white dark:bg-gray-950 transition-colors p-3 border-r-[1px] border-b-[1px] flex flex-row-reverse lg:flex-col justify-end items-center gap-4 relative group"
        onClick={() => setOpen(id)}
      >
        <span
          style={{
            writingMode: 'vertical-lr',
          }}
          className="hidden lg:block text-lg rotate-180"
        >
          {title}
        </span>
        <span className="block lg:hidden text-lg">{title}</span>
        <span className="w-4 h-4 bg-white dark:bg-gray-950 group-hover:bg-slate-50 transition-colors border-r-[1px] border-b-[1px] lg:border-b-0 lg:border-t-[1px] border-slate-200 dark:border-gray-800 rotate-45 absolute bottom-0 lg:bottom-[50%] right-[50%] lg:right-0 translate-y-[50%] translate-x-[50%] z-20" />
      </button>

      {isOpen && content}
    </>
  );
};

export function AgentBuilder() {
  return (
    <div className="overflow-hidden dark:bg-gray-900 bg-gray-100 pt-10 px-6 lg:px-8 text-gray-600 dark:text-gray-400">
      <div className="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 sm:gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-2">
        <div className="lg:pt-4 lg:pr-8">
          <div className="lg:max-w-lg">
            <p className="mt-2 text-2xl font-semibold tracking-tight text-pretty">
              Build your way–with or without code
            </p>
            <p className="mt-2">
              Build only what's unique to your brand and let rapida handle the
              rest.
            </p>
            <dl className="mt-10 max-w-xl space-y-8 text-base/7 lg:max-w-none">
              {[
                {
                  name: 'plug & play, no extra steps',
                  description:
                    'Choose your AI model, connect your tools, and define exactly how your agent talks, behaves, and helps users.',
                  icon: AdjustmentsHorizontalIcon,
                },
                {
                  icon: WrenchScrewdriverIcon,
                  name: 'reliability and scale',
                  description:
                    'Manage hundreds of millions of interactions and effectively handle peak loads without compromising service quality.',
                },
              ].map(feature => (
                <div key={feature.name} className="relative flex flex-col">
                  <feature.icon
                    aria-hidden="true"
                    className="size-5 text-blue-600 dark:text-blue-400"
                  />
                  <p className="font-semibold mt-4 capitalize">
                    {feature.name}
                  </p>{' '}
                  <p className="mt-1">{feature.description}</p>
                </div>
              ))}
            </dl>
          </div>
        </div>
        <img
          alt="Product screenshot light"
          src="/images/screenshots/agent-builder-light.png"
          width={2432}
          height={1442}
          className="block w-3xl max-w-none rounded-xl border shadow-xl ring-1 ring-white/10 sm:w-228 md:-ml-4 lg:-ml-0 dark:hidden"
        />
        <img
          alt="Product screenshot dark"
          src="/images/screenshots/agent-builder-dark.png"
          width={2432}
          height={1442}
          className="hidden w-3xl max-w-none rounded-xl border shadow-xl ring-1 ring-white/10 sm:w-228 md:-ml-4 lg:-ml-0 dark:block"
        />
      </div>
    </div>
  );
}

export function Websocket() {
  return (
    <div className="overflow-hidden dark:bg-gray-900 bg-gray-100 pt-10 px-6 lg:px-8 text-gray-600 dark:text-gray-400">
      <div className="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 sm:gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-2">
        <div className="lg:pt-4 lg:pr-8">
          <div className="lg:max-w-lg">
            <p className="mt-2 text-2xl font-semibold tracking-tight text-pretty">
              Connect with WebSocket — run any workflow, anywhere
            </p>
            <p className="mt-3 text-lg/8">
              Instantly link your workflow engine—local.
            </p>
            <dl className="mt-10 max-w-xl space-y-8 text-base/7 lg:max-w-none">
              {[
                {
                  name: 'connect anything, run anywhere',
                  description:
                    'Connect any workflow engine—local or remote—and bring real-time execution directly into your voice experiences.',
                  icon: CloudArrowUpIcon,
                },
              ].map(feature => (
                <div key={feature.name} className="relative flex flex-col">
                  <feature.icon
                    aria-hidden="true"
                    className="size-5 text-blue-600 dark:text-blue-400"
                  />
                  <p className="font-semibold mt-4 capitalize">
                    {feature.name}
                  </p>{' '}
                  <p className="mt-2">{feature.description}</p>
                </div>
              ))}
            </dl>
          </div>
        </div>
        <img
          alt="Product screenshot light"
          src="/images/screenshots/websocket-light.png"
          width={2432}
          height={1442}
          className="block w-3xl max-w-none rounded-xl border shadow-xl ring-1 ring-white/10 sm:w-228 md:-ml-4 lg:-ml-0 dark:hidden"
        />
        <img
          alt="Product screenshot dark"
          src="/images/screenshots/websocket-dark.png"
          width={2432}
          height={1442}
          className="hidden w-3xl max-w-none rounded-xl border shadow-xl ring-1 ring-white/10 sm:w-228 md:-ml-4 lg:-ml-0 dark:block"
        />
      </div>
    </div>
  );
}

export function AgentKit() {
  return (
    <div className="overflow-hidden dark:bg-gray-900 bg-gray-100 pt-10 px-6 lg:px-8 text-gray-600 dark:text-gray-400">
      <div className="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 sm:gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-2">
        <div className="lg:pt-4 lg:pr-8">
          <div className="lg:max-w-lg">
            <p className="mt-2 text-2xl font-semibold tracking-tight text-pretty">
              AgentKit — power your platform with enterprise-grade voice
            </p>
            <p className="mt-3 text-lg/8">
              A gRPC-optimized SDK that enables real-time voice capabilities to
              their existing agentic platforms.
            </p>
            <dl className="mt-10 max-w-xl space-y-8 text-base/7 lg:max-w-none">
              {[
                {
                  name: 'enterprise-scale lifecycle management',
                  description:
                    'Gain full control over your agent’s lifecycle—from event orchestration to tool calls—backed by gRPC performance and production-grade reliability.',
                  icon: ArrowsRightLeftIcon,
                },
              ].map(feature => (
                <div key={feature.name} className="relative flex flex-col">
                  <feature.icon
                    aria-hidden="true"
                    className="size-5 text-blue-600 dark:text-blue-400"
                  />
                  <p className="font-semibold mt-4 capitalize">
                    {feature.name}
                  </p>{' '}
                  <p className="mt-2">{feature.description}</p>
                </div>
              ))}
            </dl>
          </div>
        </div>
        <img
          alt="Product screenshot light"
          src="/images/screenshots/agent-builder-light.png"
          width={2432}
          height={1442}
          className="block w-3xl max-w-none rounded-xl border shadow-xl ring-1 ring-white/10 sm:w-228 md:-ml-4 lg:-ml-0 dark:hidden"
        />
        <img
          alt="Product screenshot dark"
          src="/images/screenshots/agent-builder-dark.png"
          width={2432}
          height={1442}
          className="hidden w-3xl max-w-none rounded-xl border shadow-xl ring-1 ring-white/10 sm:w-228 md:-ml-4 lg:-ml-0 dark:block"
        />
      </div>
    </div>
  );
}
