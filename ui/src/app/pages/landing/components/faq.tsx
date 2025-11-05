export const Question = () => (
  <section
    className="border-y mt-30 grid grid-cols-1 gap-10 lg:grid-cols-2"
    id="faqs"
  >
    <div className="lg:border-r ">
      <div className="grid grid-cols-1 gap-y-2 px-4 py-2 border-b sm:px-2 lg:border-b/half">
        <p className="text-[2.5rem]/none font-medium tracking-tight text-pretty">
          Frequently asked questions
        </p>
      </div>
    </div>
    <div className="lg:border-l ">
      <div className="grid grid-cols-1 gap-10">
        {/*  */}

        <div className="group">
          <h3 className="px-4 py-2 sm:px-2 font-mono text-[0.8125rem]/6 font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
            General
          </h3>
          <dl>
            {[
              {
                question: 'What is Rapida?',
                answer:
                  'Rapida is a voice orchestration platform that enables communications platforms to deliver intelligent voice experiences as a native product without rebuilds or brand fragmentation.',
              },
              {
                question: "What does 'voice orchestration' mean?",
                answer:
                  'Voice orchestration is a layer that enables seamless integration of voice capabilities across multiple channels and systems, allowing you to build sophisticated voice experiences with memory, context, and multi-turn conversations.',
              },
              {
                question: 'Who should use Rapida?',
                answer:
                  'Rapida is designed for startups, product teams, enterprises, and communications platforms looking to add voice AI capabilities to their products without significant engineering overhead.',
              },
              {
                question: 'What can I build with Rapida?',
                answer:
                  'You can build a wide variety of voice experiences including customer service agents, voice reminders, transaction systems, appointment bookings, surveys, collections agents, and more across multiple industries like retail, e-commerce, lending, and insurance.',
              },
              {
                question: 'What are the main components of Rapida?',
                answer:
                  'Rapida consists of three core layers: The Orchestrator (for handling conversations with memory), Prebuilt Assets (ready-made voice experiences), and Custom Assets (for building specialized voice solutions).',
              },
              {
                question:
                  'Does Rapida support multi-channel voice experiences?',
                answer:
                  'Yes, Rapida supports seamless voice experiences across voice, chat, email, and WhatsApp channels.',
              },
              {
                question: 'What integrations does Rapida offer?',
                answer:
                  'Rapida integrates with third-party systems through APIs, supports LLM integration (including OpenAI and other providers), SIP integration, CRM systems, Knowledge Management tools, and webhooks & event triggers.',
              },
              {
                question: 'How do you ensure voice quality?',
                answer:
                  'Rapida provides observability and telemetry across all interactions, monitoring metrics, measuring, and improving accuracy and performance to ensure high-quality voice experiences.',
              },
            ].map((faq, index) => (
              <details
                className="group border-t  px-4 py-3 sm:px-2"
                key={index}
              >
                <summary
                  id="faq-what-does-lifetime-access-mean-exactly"
                  className="flex w-full cursor-pointer justify-between gap-4 select-none group-open:text-blue-500 dark:group-open:text-blue-500 [&::-webkit-details-marker]:hidden"
                >
                  <div className="text-left text-sm/7 font-semibold text-pretty">
                    {faq.question}
                  </div>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 group-open:hidden"
                  >
                    <path d="M8.75 3.75a.75.75 0 0 0-1.5 0v3.5h-3.5a.75.75 0 0 0 0 1.5h3.5v3.5a.75.75 0 0 0 1.5 0v-3.5h3.5a.75.75 0 0 0 0-1.5h-3.5v-3.5Z" />
                  </svg>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 not-group-open:hidden"
                  >
                    <path d="M3.75 7.25a.75.75 0 0 0 0 1.5h8.5a.75.75 0 0 0 0-1.5h-8.5Z" />
                  </svg>
                </summary>
                <div className='mt-4 grid grid-cols-1 gap-6 text-sm/7 text-gray-600 dark:text-gray-400 [&_strong]:font-semibold [&_strong]:text-gray-950 dark:[&_strong]:text-white [&_h2]:not-first:mt-15 [&_h3]:not-first:mt-6 [&_h2]:text-lg/8 [&_h2]:font-semibold [&_h2]:text-gray-950 dark:[&_h2]:text-white [&_h3]:text-base/7 [&_h3]:font-semibold [&_h3]:text-gray-950 dark:[&_h3]:text-white [&_img]:outline [&_img]:-outline-offset-1 [&_img]:outline-black/5 [&_img]:dark:outline-white/10 [&_a]:font-semibold [&_a]:text-gray-950 [&_a]:underline [&_a]:decoration-sky-400 [&_a]:underline-offset-4 [&_a]:hover:text-sky-500 dark:[&_a]:text-white dark:[&_a]:hover:text-sky-500 [&_li]:relative [&_li]:before:absolute [&_li]:before:-top-0.5 [&_li]:before:-left-6 [&_li]:before:text-gray-300 [&_li]:before:content-["▪"] [&_ul]:pl-9 [&_pre]:overflow-x-auto [&_pre]:rounded-xl [&_pre]:border-4 [&_pre]:border-gray-950 [&_pre]:bg-gray-900 [&_pre]:p-4 [&_pre]:text-white [&_pre]:outline-1 [&_pre]:-outline-offset-5 [&_pre]:outline-white/10 dark:[&_pre]:border-[color-mix(in_oklab,var(--color-gray-950),white_10%)] [&_pre_code]:bg-gray-900 [&_code]:not-in-[pre]:font-medium [&_code]:not-in-[pre]:whitespace-nowrap [&_code]:not-in-[pre]:text-gray-950 [&_code]:not-in-[pre]:before:content-["\`"] [&_code]:not-in-[pre]:after:content-["\`"] dark:[&_code]:not-in-[pre]:text-white'>
                  <p>{faq.answer}</p>
                </div>
              </details>
            ))}
          </dl>
        </div>
        {/*  */}
        <div className="group">
          <h3 className="px-4 py-2 sm:px-2 font-mono text-[0.8125rem]/6 font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
            Compatibility
          </h3>
          <dl>
            {[
              {
                question: 'How does Rapida handle latency?',
                answer:
                  'Rapida maintains predictable, low latency performance even at scale, allowing you to focus on delivering new customer experiences instead of managing infrastructure.',
              },
              {
                question: 'Can Rapida scale to handle high volumes?',
                answer:
                  'Yes, Rapida is built for scale with support for unlimited voice minutes, flexible scaling, multi-region deployment options, and on-premises deployment capabilities.',
              },
              {
                question: 'Is there an edge network option?',
                answer:
                  'Yes, Rapida offers a Distributed Edge Network for processing voice experiences closer to your users.',
              },
              {
                question: 'What security and compliance features are included?',
                answer:
                  'Rapida includes data encryption, full conversation analysis, crystal-clear call quality with AHI metrics, BYOB (Bring Your Own LLM/STT/TTS) options, and knowledge integration for accurate contextual information.',
              },
            ].map((faq, index) => (
              <details
                key={index}
                className="group border-t  px-4 py-3 sm:px-2"
              >
                <summary
                  id="faq-what-does-lifetime-access-mean-exactly"
                  className="flex w-full cursor-pointer justify-between gap-4 select-none group-open:text-blue-500 dark:group-open:text-blue-500 [&::-webkit-details-marker]:hidden"
                >
                  <div className="text-left text-sm/7 font-semibold text-pretty">
                    {faq.question}
                  </div>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 group-open:hidden"
                  >
                    <path d="M8.75 3.75a.75.75 0 0 0-1.5 0v3.5h-3.5a.75.75 0 0 0 0 1.5h3.5v3.5a.75.75 0 0 0 1.5 0v-3.5h3.5a.75.75 0 0 0 0-1.5h-3.5v-3.5Z" />
                  </svg>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 not-group-open:hidden"
                  >
                    <path d="M3.75 7.25a.75.75 0 0 0 0 1.5h8.5a.75.75 0 0 0 0-1.5h-8.5Z" />
                  </svg>
                </summary>
                <div className='mt-4 grid grid-cols-1 gap-6 text-sm/7 text-gray-600 dark:text-gray-400 [&_strong]:font-semibold [&_strong]:text-gray-950 dark:[&_strong]:text-white [&_h2]:not-first:mt-15 [&_h3]:not-first:mt-6 [&_h2]:text-lg/8 [&_h2]:font-semibold [&_h2]:text-gray-950 dark:[&_h2]:text-white [&_h3]:text-base/7 [&_h3]:font-semibold [&_h3]:text-gray-950 dark:[&_h3]:text-white [&_img]:outline [&_img]:-outline-offset-1 [&_img]:outline-black/5 [&_img]:dark:outline-white/10 [&_a]:font-semibold [&_a]:text-gray-950 [&_a]:underline [&_a]:decoration-sky-400 [&_a]:underline-offset-4 [&_a]:hover:text-sky-500 dark:[&_a]:text-white dark:[&_a]:hover:text-sky-500 [&_li]:relative [&_li]:before:absolute [&_li]:before:-top-0.5 [&_li]:before:-left-6 [&_li]:before:text-gray-300 [&_li]:before:content-["▪"] [&_ul]:pl-9 [&_pre]:overflow-x-auto [&_pre]:rounded-xl [&_pre]:border-4 [&_pre]:border-gray-950 [&_pre]:bg-gray-900 [&_pre]:p-4 [&_pre]:text-white [&_pre]:outline-1 [&_pre]:-outline-offset-5 [&_pre]:outline-white/10 dark:[&_pre]:border-[color-mix(in_oklab,var(--color-gray-950),white_10%)] [&_pre_code]:bg-gray-900 [&_code]:not-in-[pre]:font-medium [&_code]:not-in-[pre]:whitespace-nowrap [&_code]:not-in-[pre]:text-gray-950 [&_code]:not-in-[pre]:before:content-["\`"] [&_code]:not-in-[pre]:after:content-["\`"] dark:[&_code]:not-in-[pre]:text-white'>
                  <p>{faq.answer}</p>
                </div>
              </details>
            ))}
          </dl>
        </div>
        <div className="group">
          <h3 className="px-4 py-2 sm:px-2 font-mono text-[0.8125rem]/6 font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
            Pricing
          </h3>
          <dl>
            {[
              {
                question: 'What is the pricing model?',
                answer:
                  'Rapida offers usage-based scaling with transparent pricing. The startup plan starts at US$599 per month for one-time payment or more flexible plans for large use cases.',
              },
              {
                question: "What's included in the startup plan?",
                answer:
                  'The startup plan includes 50,000 voice minutes, full SDK & orchestration suite with complete access to all development tools and APIs, and continuous updates with automatic new features and improvements.',
              },
              {
                question: 'Is there a free trial?',
                answer:
                  "Yes, you can get a demo or start building with Rapida's platform.",
              },
              {
                question:
                  'What does "unlimited projects and production environments" mean?',
                answer:
                  'You can scale across as many use cases, regions, and deployments as needed without additional licensing restrictions.',
              },
              {
                question: 'How quickly can I get started?',
                answer:
                  'You can get everything set up within days or months and jump-start your deployment with pre-built templates for common use cases and ready-made integrations.',
              },
            ].map((faq, index) => (
              <details
                key={index}
                className="group border-t  px-4 py-3 sm:px-2"
              >
                <summary
                  id="faq-what-does-lifetime-access-mean-exactly"
                  className="flex w-full cursor-pointer justify-between gap-4 select-none group-open:text-blue-500 dark:group-open:text-blue-500 [&::-webkit-details-marker]:hidden"
                >
                  <div className="text-left text-sm/7 font-semibold text-pretty">
                    {faq.question}
                  </div>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 group-open:hidden"
                  >
                    <path d="M8.75 3.75a.75.75 0 0 0-1.5 0v3.5h-3.5a.75.75 0 0 0 0 1.5h3.5v3.5a.75.75 0 0 0 1.5 0v-3.5h3.5a.75.75 0 0 0 0-1.5h-3.5v-3.5Z" />
                  </svg>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 not-group-open:hidden"
                  >
                    <path d="M3.75 7.25a.75.75 0 0 0 0 1.5h8.5a.75.75 0 0 0 0-1.5h-8.5Z" />
                  </svg>
                </summary>
                <div className='mt-4 grid grid-cols-1 gap-6 text-sm/7 text-gray-600 dark:text-gray-400 [&_strong]:font-semibold [&_strong]:text-gray-950 dark:[&_strong]:text-white [&_h2]:not-first:mt-15 [&_h3]:not-first:mt-6 [&_h2]:text-lg/8 [&_h2]:font-semibold [&_h2]:text-gray-950 dark:[&_h2]:text-white [&_h3]:text-base/7 [&_h3]:font-semibold [&_h3]:text-gray-950 dark:[&_h3]:text-white [&_img]:outline [&_img]:-outline-offset-1 [&_img]:outline-black/5 [&_img]:dark:outline-white/10 [&_a]:font-semibold [&_a]:text-gray-950 [&_a]:underline [&_a]:decoration-sky-400 [&_a]:underline-offset-4 [&_a]:hover:text-sky-500 dark:[&_a]:text-white dark:[&_a]:hover:text-sky-500 [&_li]:relative [&_li]:before:absolute [&_li]:before:-top-0.5 [&_li]:before:-left-6 [&_li]:before:text-gray-300 [&_li]:before:content-["▪"] [&_ul]:pl-9 [&_pre]:overflow-x-auto [&_pre]:rounded-xl [&_pre]:border-4 [&_pre]:border-gray-950 [&_pre]:bg-gray-900 [&_pre]:p-4 [&_pre]:text-white [&_pre]:outline-1 [&_pre]:-outline-offset-5 [&_pre]:outline-white/10 dark:[&_pre]:border-[color-mix(in_oklab,var(--color-gray-950),white_10%)] [&_pre_code]:bg-gray-900 [&_code]:not-in-[pre]:font-medium [&_code]:not-in-[pre]:whitespace-nowrap [&_code]:not-in-[pre]:text-gray-950 [&_code]:not-in-[pre]:before:content-["\`"] [&_code]:not-in-[pre]:after:content-["\`"] dark:[&_code]:not-in-[pre]:text-white'>
                  <p>{faq.answer}</p>
                </div>
              </details>
            ))}
          </dl>
        </div>
        <div className="group">
          <h3 className="px-4 py-2 sm:px-2 font-mono text-[0.8125rem]/6 font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
            Support
          </h3>
          <dl>
            {[
              {
                question: 'How can I get technical support?',
                answer:
                  'Rapida offers 24/7 support for all customers, with dedicated deployment engineers for implementation assistance.',
              },
              {
                question: 'Who should I contact for enterprise solutions?',
                answer:
                  'Visit the "Talk to us" section or contact Rapida\'s sales team for custom enterprise pricing and dedicated support.',
              },
              {
                question: 'Where can I find documentation and resources?',
                answer:
                  'Rapida provides comprehensive documentation, API references, and integration guides through the developer platform.',
              },
            ].map((faq, index) => (
              <details
                key={index}
                className="group border-t  px-4 py-3 sm:px-2"
              >
                <summary
                  id="faq-what-does-lifetime-access-mean-exactly"
                  className="flex w-full cursor-pointer justify-between gap-4 select-none group-open:text-blue-500 dark:group-open:text-blue-500 [&::-webkit-details-marker]:hidden"
                >
                  <div className="text-left text-sm/7 font-semibold text-pretty">
                    {faq.question}
                  </div>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 group-open:hidden"
                  >
                    <path d="M8.75 3.75a.75.75 0 0 0-1.5 0v3.5h-3.5a.75.75 0 0 0 0 1.5h3.5v3.5a.75.75 0 0 0 1.5 0v-3.5h3.5a.75.75 0 0 0 0-1.5h-3.5v-3.5Z" />
                  </svg>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    fill="currentColor"
                    aria-hidden="true"
                    data-slot="icon"
                    className="h-7 w-4 not-group-open:hidden"
                  >
                    <path d="M3.75 7.25a.75.75 0 0 0 0 1.5h8.5a.75.75 0 0 0 0-1.5h-8.5Z" />
                  </svg>
                </summary>
                <div className='mt-4 grid grid-cols-1 gap-6 text-sm/7 text-gray-600 dark:text-gray-400 [&_strong]:font-semibold [&_strong]:text-gray-950 dark:[&_strong]:text-white [&_h2]:not-first:mt-15 [&_h3]:not-first:mt-6 [&_h2]:text-lg/8 [&_h2]:font-semibold [&_h2]:text-gray-950 dark:[&_h2]:text-white [&_h3]:text-base/7 [&_h3]:font-semibold [&_h3]:text-gray-950 dark:[&_h3]:text-white [&_img]:outline [&_img]:-outline-offset-1 [&_img]:outline-black/5 [&_img]:dark:outline-white/10 [&_a]:font-semibold [&_a]:text-gray-950 [&_a]:underline [&_a]:decoration-sky-400 [&_a]:underline-offset-4 [&_a]:hover:text-sky-500 dark:[&_a]:text-white dark:[&_a]:hover:text-sky-500 [&_li]:relative [&_li]:before:absolute [&_li]:before:-top-0.5 [&_li]:before:-left-6 [&_li]:before:text-gray-300 [&_li]:before:content-["▪"] [&_ul]:pl-9 [&_pre]:overflow-x-auto [&_pre]:rounded-xl [&_pre]:border-4 [&_pre]:border-gray-950 [&_pre]:bg-gray-900 [&_pre]:p-4 [&_pre]:text-white [&_pre]:outline-1 [&_pre]:-outline-offset-5 [&_pre]:outline-white/10 dark:[&_pre]:border-[color-mix(in_oklab,var(--color-gray-950),white_10%)] [&_pre_code]:bg-gray-900 [&_code]:not-in-[pre]:font-medium [&_code]:not-in-[pre]:whitespace-nowrap [&_code]:not-in-[pre]:text-gray-950 [&_code]:not-in-[pre]:before:content-["\`"] [&_code]:not-in-[pre]:after:content-["\`"] dark:[&_code]:not-in-[pre]:text-white'>
                  <p>{faq.answer}</p>
                </div>
              </details>
            ))}
          </dl>
        </div>
      </div>
    </div>
  </section>
);
