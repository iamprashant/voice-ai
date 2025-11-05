import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/react';
import { PhoneIncoming, PhoneOutgoing, Play } from 'lucide-react';
export const Usecases = () => {
  return (
    <section className="mt-20 sm:mt-40">
      <div className="border-y px-4 py-2 sm:px-2">
        <h2 className="max-w-3xl text-3xl font-medium tracking-tight text-pretty md:text-[2.5rem]/14">
          One Platform. Many Voice Experiences
        </h2>
        <p className="mt-4 max-w-2xl text-base text-gray-600 dark:text-gray-400">
          Build, deploy, and scale production-ready Voice AI. Connect your
          systems, define memory and logic, and deliver brand-native
          conversations.
        </p>
      </div>
      <TabGroup className="mt-16 lg:mt-20">
        <div className="border-y">
          <TabList className="flex whitespace-nowrap">
            <Tab className="border-r p-3 hover:bg-gray-950/2 focus:not-data-focus:outline-hidden data-selected:bg-gray-950/2.5 sm:px-6 dark:hover:bg-white/4 dark:data-selected:bg-white/5 cursor-pointer">
              <p className="font-mono text-sm font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
                BFSI
              </p>
            </Tab>
            <Tab className="border-r p-3 hover:bg-gray-950/2 focus:not-data-focus:outline-hidden data-selected:bg-gray-950/2.5 sm:px-6 dark:hover:bg-white/4 dark:data-selected:bg-white/5 cursor-pointer">
              <p className="font-mono text-sm font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
                Retail
              </p>
            </Tab>
            <Tab className="border-r p-3 hover:bg-gray-950/2 focus:not-data-focus:outline-hidden data-selected:bg-gray-950/2.5 sm:px-6 dark:hover:bg-white/4 dark:data-selected:bg-white/5 cursor-pointer">
              <p className="font-mono text-sm font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
                Ecommerce
              </p>
            </Tab>
            <Tab className="border-r p-3 hover:bg-gray-950/2 focus:not-data-focus:outline-hidden data-selected:bg-gray-950/2.5 sm:px-6 dark:hover:bg-white/4 dark:data-selected:bg-white/5 cursor-pointer">
              <p className="font-mono text-sm font-medium tracking-widest text-pretty uppercase text-gray-600 dark:text-gray-400">
                Coaching & Education
              </p>
            </Tab>
          </TabList>
        </div>
        <TabPanels>
          <TabPanel className="relative mt-4">
            <div className="pointer-events-none absolute inset-0 z-10 grid grid-cols-1 gap-2 max-sm:hidden sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4">
              <div className="border-r" />
              <div className="border-x" />
              <div className="border-x" />
              <div className="border-l" />
            </div>
            <ul className="gap-2 sm:grid sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4 divide-y border-t">
              {[
                {
                  name: 'Collections Reminder',
                  subtitle:
                    'Automated outbound calls to remind customers of due payments and capture promises-to-pay (PTP).',
                  primary_metric: 'Right-Party Contact & Payment Conversion',
                  persona_name: 'Lending / Credit Card Customer',
                  category: 'outbound',
                  url: 'https://example.com/call-recordings/collections-reminder.mp3',
                  minute: 5,
                },
                {
                  name: 'Document Chase',
                  subtitle:
                    'Outbound calls guiding customers through required documentation for applications or KYB.',
                  primary_metric: 'First-Time-Right Docs & Time to Complete',
                  persona_name: 'Loan Applicant / Merchant',
                  category: 'outbound',
                  url: 'https://example.com/call-recordings/document-chase.mp3',
                  minute: 5,
                },
                {
                  name: 'Appointment Booking',
                  subtitle:
                    'Outbound voice calls to schedule medical, insurance, video-KYC, or wealth advisory appointments.',
                  primary_metric: 'Booked Rate & No-Show Rate',
                  persona_name: 'Insurance / Lending Customer',
                  category: 'outbound',
                  url: 'https://example.com/call-recordings/appointment-booking.mp3',
                  minute: 5,
                },
                {
                  name: 'Status Check',
                  subtitle:
                    'Inbound calls to check the status of applications, claims, or payouts with AI capturing intent and providing next steps.',
                  primary_metric: 'Containment & First-Contact Resolution',
                  persona_name: 'Loan / Insurance / Payment Customer',
                  category: 'inbound',
                  url: 'https://example.com/call-recordings/status-check.mp3',
                  minute: 5,
                },
                {
                  name: 'Dispute Intake',
                  subtitle:
                    'Inbound calls capturing dispute or chargeback requests with AI guidance to ensure first-time-right submissions.',
                  primary_metric:
                    'First-Time-Right Intake & Invalid Dispute Rate',
                  persona_name: 'Cardholder / Payment Customer',
                  category: 'inbound',
                  url: 'https://example.com/call-recordings/dispute-intake.mp3',
                  minute: 5,
                },
                {
                  name: 'IVR Deflection',
                  subtitle:
                    'Inbound AI-powered IVR handling top repetitive intents, FAQs, and routing to reduce live agent load.',
                  primary_metric: 'Containment & Average Wait Time',
                  persona_name: 'All BFSI Customers',
                  category: 'inbound',
                  url: 'https://example.com/call-recordings/ivr-deflection.mp3',
                  minute: 5,
                },
                {
                  name: 'Early-Stage Collections',
                  subtitle:
                    'Outbound calls for high-volume early-stage collections, capturing payment plans to improve promise-to-pay rates.',
                  primary_metric: '% Promises Kept & Roll Rate',
                  persona_name: 'Lending / Cards / BNPL Customer',
                  category: 'outbound',
                  url: 'https://example.com/call-recordings/early-stage-collections.mp3',
                  minute: 5,
                },
                {
                  name: 'Hot Lead Recovery',
                  subtitle:
                    'Outbound calls to re-engage abandoned or minutes-fresh leads, increasing completion of applications.',
                  primary_metric: 'Resume Rate & Time to Resume',
                  persona_name: 'Lending / Insurance Prospect',
                  category: 'outbound',
                  url: 'https://example.com/call-recordings/hot-lead-recovery.mp3',
                  minute: 5,
                },
              ].map((us, index) => (
                <li
                  key={index}
                  className="block transition hover:bg-gray-950/2.5 focus-visible:relative focus-visible:z-20 focus-visible:bg-white/75 focus-visible:backdrop-blur-xs dark:hover:bg-white/5 border-l border-b border-r"
                >
                  <div className="px-4 py-2 sm:px-2">
                    {us.category === 'outbound' ? (
                      <PhoneOutgoing
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    ) : (
                      <PhoneIncoming
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    )}
                    <div className="flex items-center gap-2 mt-4">
                      <h3 className="text-base/7 font-semibold">
                        <div className="absolute inset-0" />
                        {us.name}
                      </h3>
                    </div>

                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                      {us.subtitle}
                    </p>
                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                      {us.primary_metric}
                    </p>
                  </div>
                  <div className="px-4 py-2 sm:px-2 border-t">
                    <figcaption className="flex space-x-2">
                      <Play
                        className="aspect-square size-9 rounded-full outline -outline-offset-1 outline-gray-950/5 dark:outline-white/10 p-2.5 bg-blue-600/20 fill-blue-600 text-blue-600"
                        strokeWidth={1.5}
                      />
                      <div className="text-sm line-height">
                        <p className="font-medium">{us.persona_name}</p>
                        <p className="text-gray-600 dark:text-gray-400">
                          {us.minute} minutes
                        </p>
                      </div>
                    </figcaption>
                  </div>
                </li>
              ))}
            </ul>
            <div className="pointer-events-none absolute inset-x-0 -bottom-0.5 z-10 flex h-100 items-end justify-center bg-linear-to-b to-white pb-8 max-sm:hidden dark:to-gray-950">
              <a
                className="pointer-events-auto gap-2 inline-flex justify-center rounded-full text-base  font-medium  focus-visible:outline-2 focus-visible:outline-offset-2 bg-gray-950 text-white hover:bg-gray-800 focus-visible:outline-gray-950 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white dark:hover:bg-gray-200 dark:focus-visible:outline-white dark:focus-visible:outline-white px-4 py-2"
                href="/demo"
              >
                Consult an expert
                <svg
                  fill="currentColor"
                  aria-hidden="true"
                  viewBox="0 0 10 10"
                  className="-mr-0.5 w-2.5"
                >
                  <path d="M4.85355 0.146423L9.70711 4.99998L4.85355 9.85353L4.14645 9.14642L7.79289 5.49998H0V4.49998H7.79289L4.14645 0.85353L4.85355 0.146423Z" />
                </svg>
              </a>
            </div>
          </TabPanel>

          {/* retail */}
          <TabPanel className="relative mt-4">
            <div className="pointer-events-none absolute inset-0 z-10 grid grid-cols-1 gap-2 max-sm:hidden sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4">
              <div className="border-r" />
              <div className="border-x" />
              <div className="border-x" />
              <div className="border-l" />
            </div>
            <ul className="gap-2 sm:grid sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4 divide-y border-t">
              {[
                {
                  name: 'Order Status & Delivery ETA',
                  subtitle:
                    "AI answers 'Where is my order?' calls instantly with real-time tracking and updates.",
                  category: 'Inbound',
                  primary_metric:
                    'Containment Rate, AHT, First-Contact Resolution',
                  persona_name: 'CX Head / Operations Manager',
                  url: 'https://example.com/demo/order-status-call',
                  minute: 4,
                },
                {
                  name: 'COD Payment Reminder',
                  subtitle:
                    'Outbound AI calls to confirm delivery and payment readiness, reducing failed COD deliveries.',
                  category: 'Outbound',
                  primary_metric:
                    'Right-Party Contact, Confirmation Rate, Failed Delivery Rate',
                  persona_name: 'Logistics Manager / E-Commerce Lead',
                  url: 'https://example.com/demo/cod-reminder-call',
                  minute: 4,
                },
                {
                  name: 'Abandoned Cart Recovery',
                  subtitle:
                    'Voice AI follows up on abandoned carts within minutes to recover lost sales.',
                  category: 'Outbound',
                  primary_metric:
                    'Call-to-Conversion Rate, Cart Recovery %, Containment',
                  persona_name: 'Growth / Retention Manager',
                  url: 'https://example.com/demo/cart-recovery-call',
                  minute: 5,
                },
                {
                  name: 'Feedback & CSAT Capture',
                  subtitle:
                    'Post-purchase feedback calls with natural voice interactions for better response rates.',
                  category: 'Outbound',
                  primary_metric:
                    'Response Rate, CSAT %, Feedback Completion Rate',
                  persona_name: 'Customer Experience Manager',
                  url: 'https://example.com/demo/feedback-call',
                  minute: 4,
                },
                {
                  name: 'Return / Exchange Assistance',
                  subtitle:
                    'Inbound call automation to initiate, track, or schedule returns and exchanges.',
                  category: 'Inbound',
                  primary_metric:
                    'Containment Rate, AHT, First-Time-Right Resolution',
                  persona_name: 'Customer Support Head',
                  url: 'https://example.com/demo/return-assistance-call',
                  minute: 5,
                },
                {
                  name: 'Store Visit Reminders',
                  subtitle:
                    'Outbound reminder calls for in-store pickup, loyalty events, or scheduled appointments.',
                  category: 'Outbound',
                  primary_metric:
                    'Attendance Rate, Confirmation %, No-Show Rate',
                  persona_name: 'Retail Operations Manager',
                  url: 'https://example.com/demo/store-reminder-call',
                  minute: 3,
                },
                {
                  name: 'Loyalty Upsell Calls',
                  subtitle:
                    'AI reaches existing customers with personalized offers or loyalty plan renewals.',
                  category: 'Outbound',
                  primary_metric:
                    'Offer Acceptance Rate, Conversion %, Containment',
                  persona_name: 'CRM / Loyalty Lead',
                  url: 'https://example.com/demo/loyalty-upsell-call',
                  minute: 5,
                },
                {
                  name: 'Fraud Verification (Order/Payment)',
                  subtitle:
                    'AI calls to verify suspicious or high-value transactions before shipping.',
                  category: 'Outbound',
                  primary_metric:
                    'Verification Completion Rate, False Positive %, Containment',
                  persona_name: 'Fraud Operations Manager',
                  url: 'https://example.com/demo/fraud-verification-call',
                  minute: 4,
                },
              ].map((us, index) => (
                <li
                  key={index}
                  className="block transition hover:bg-gray-950/2.5 focus-visible:relative focus-visible:z-20 focus-visible:bg-white/75 focus-visible:backdrop-blur-xs dark:hover:bg-white/5 border-l border-b border-r"
                >
                  <div className="px-4 py-2 sm:px-2">
                    {us.category === 'outbound' ? (
                      <PhoneOutgoing
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    ) : (
                      <PhoneIncoming
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    )}
                    <div className="flex items-center gap-2 mt-4">
                      <h3 className="text-base/7 font-semibold">
                        <div className="absolute inset-0" />
                        {us.name}
                      </h3>
                    </div>

                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400 line-clamp-2">
                      {us.subtitle}
                    </p>
                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                      {us.primary_metric}
                    </p>
                  </div>
                  <div className="px-4 py-2 sm:px-2 border-t">
                    <figcaption className="flex space-x-2">
                      <Play
                        className="aspect-square size-9 rounded-full outline -outline-offset-1 outline-gray-950/5 dark:outline-white/10 p-2.5 bg-blue-600/20 fill-blue-600 text-blue-600"
                        strokeWidth={1.5}
                      />
                      <div className="text-sm line-height">
                        <p className="font-medium">{us.persona_name}</p>
                        <p className="text-gray-600 dark:text-gray-400">
                          {us.minute} minutes
                        </p>
                      </div>
                    </figcaption>
                  </div>
                </li>
              ))}
            </ul>
            <div className="pointer-events-none absolute inset-x-0 -bottom-0.5 z-10 flex h-100 items-end justify-center bg-linear-to-b to-white pb-8 max-sm:hidden dark:to-gray-950">
              <a
                className="pointer-events-auto gap-2 inline-flex justify-center rounded-full text-base  font-medium  focus-visible:outline-2 focus-visible:outline-offset-2 bg-gray-950 text-white hover:bg-gray-800 focus-visible:outline-gray-950 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white dark:hover:bg-gray-200 dark:focus-visible:outline-white dark:focus-visible:outline-white px-4 py-2"
                href="/demo"
              >
                See it in action
                <svg
                  fill="currentColor"
                  aria-hidden="true"
                  viewBox="0 0 10 10"
                  className="-mr-0.5 w-2.5"
                >
                  <path d="M4.85355 0.146423L9.70711 4.99998L4.85355 9.85353L4.14645 9.14642L7.79289 5.49998H0V4.49998H7.79289L4.14645 0.85353L4.85355 0.146423Z" />
                </svg>
              </a>
            </div>
          </TabPanel>

          {/* education */}
          <TabPanel className="relative mt-4">
            <div className="pointer-events-none absolute inset-0 z-10 grid grid-cols-1 gap-2 max-sm:hidden sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4">
              <div className="border-r" />
              <div className="border-x" />
              <div className="border-x" />
              <div className="border-l" />
            </div>
            <ul className="gap-2 sm:grid sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4 divide-y border-t">
              {[
                {
                  name: 'Enrollment Inquiry',
                  subtitle:
                    'AI handles inbound admission and course queries, captures lead intent, and schedules counselor calls.',
                  category: 'Inbound',
                  primary_metric:
                    'Containment Rate, Lead Capture %, Conversion-to-Callback',
                  persona_name: 'Admissions Head / Marketing Lead',
                  url: 'https://example.com/demo/enrollment-inquiry-call',
                  minute: 5,
                },
                {
                  name: 'Fee Payment Reminder',
                  subtitle:
                    'Automated outbound calls reminding students or parents about upcoming fee deadlines.',
                  category: 'Outbound',
                  primary_metric:
                    'Right-Party Contact, Payment Rate, Promise-to-Pay % (PTP)',
                  persona_name: 'Finance / Collections Manager',
                  url: 'https://example.com/demo/fee-reminder-call',
                  minute: 4,
                },
                {
                  name: 'Class Schedule Update',
                  subtitle:
                    'Outbound AI calls to notify students of rescheduled or canceled classes with confirmation capture.',
                  category: 'Outbound',
                  primary_metric:
                    'Reach Rate, Acknowledgment %, No-Show Reduction',
                  persona_name: 'Operations / Student Experience Manager',
                  url: 'https://example.com/demo/class-schedule-call',
                  minute: 3,
                },
                {
                  name: 'Demo Class Booking',
                  subtitle:
                    'AI calls prospective students to confirm demo class slots and send reminders.',
                  category: 'Outbound',
                  primary_metric: 'Booking Rate, Confirmation %, Show Rate',
                  persona_name: 'Admissions / Marketing Manager',
                  url: 'https://example.com/demo/demo-class-booking-call',
                  minute: 4,
                },
                {
                  name: 'Attendance Follow-Up',
                  subtitle:
                    'Outbound voice reminders for absent students to understand reasons and improve attendance.',
                  category: 'Outbound',
                  primary_metric:
                    'Response Rate, Attendance Improvement %, Containment',
                  persona_name: 'Academic Coordinator / Program Lead',
                  url: 'https://example.com/demo/attendance-followup-call',
                  minute: 4,
                },
                {
                  name: 'Exam / Assignment Reminders',
                  subtitle:
                    'Outbound AI calls to remind students of upcoming exams or assignment submissions.',
                  category: 'Outbound',
                  primary_metric:
                    'Reach Rate, Submission Rate, No-Show / Late Submission %',
                  persona_name: 'Academic Operations Manager',
                  url: 'https://example.com/demo/exam-reminder-call',
                  minute: 3,
                },
                {
                  name: 'Placement Support',
                  subtitle:
                    'AI assists with interview scheduling and verification calls during placement season.',
                  category: 'Inbound',
                  primary_metric: 'Booking Rate, First-Time-Right %, AHT',
                  persona_name: 'Placement Officer / Career Services Head',
                  url: 'https://example.com/demo/placement-support-call',
                  minute: 5,
                },
                {
                  name: 'Feedback & NPS Capture',
                  subtitle:
                    'Post-course feedback voice calls to gather student satisfaction and actionable insights.',
                  category: 'Outbound',
                  primary_metric: 'Response Rate, NPS %, Completion Rate',
                  persona_name: 'Student Experience / QA Head',
                  url: 'https://example.com/demo/feedback-nps-call',
                  minute: 4,
                },
              ].map((us, index) => (
                <li
                  key={index}
                  className="block transition hover:bg-gray-950/2.5 focus-visible:relative focus-visible:z-20 focus-visible:bg-white/75 focus-visible:backdrop-blur-xs dark:hover:bg-white/5 border-l border-b border-r"
                >
                  <div className="px-4 py-2 sm:px-2">
                    {us.category === 'outbound' ? (
                      <PhoneOutgoing
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    ) : (
                      <PhoneIncoming
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    )}
                    <div className="flex items-center gap-2 mt-4">
                      <h3 className="text-base/7 font-semibold">
                        <div className="absolute inset-0" />
                        {us.name}
                      </h3>
                    </div>

                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400 line-clamp-2">
                      {us.subtitle}
                    </p>
                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                      {us.primary_metric}
                    </p>
                  </div>
                  <div className="px-4 py-2 sm:px-2 border-t">
                    <figcaption className="flex space-x-2">
                      <Play
                        className="aspect-square size-9 rounded-full outline -outline-offset-1 outline-gray-950/5 dark:outline-white/10 p-2.5 bg-blue-600/20 fill-blue-600 text-blue-600"
                        strokeWidth={1.5}
                      />
                      <div className="text-sm line-height">
                        <p className="font-medium">{us.persona_name}</p>
                        <p className="text-gray-600 dark:text-gray-400">
                          {us.minute} minutes
                        </p>
                      </div>
                    </figcaption>
                  </div>
                </li>
              ))}
            </ul>
            <div className="pointer-events-none absolute inset-x-0 -bottom-0.5 z-10 flex h-100 items-end justify-center bg-linear-to-b to-white pb-8 max-sm:hidden dark:to-gray-950">
              <a
                className="pointer-events-auto gap-2 inline-flex justify-center rounded-full text-base  font-medium  focus-visible:outline-2 focus-visible:outline-offset-2 bg-gray-950 text-white hover:bg-gray-800 focus-visible:outline-gray-950 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white dark:hover:bg-gray-200 dark:focus-visible:outline-white dark:focus-visible:outline-white px-4 py-2"
                href="/demo"
              >
                Consult an expert
                <svg
                  fill="currentColor"
                  aria-hidden="true"
                  viewBox="0 0 10 10"
                  className="-mr-0.5 w-2.5"
                >
                  <path d="M4.85355 0.146423L9.70711 4.99998L4.85355 9.85353L4.14645 9.14642L7.79289 5.49998H0V4.49998H7.79289L4.14645 0.85353L4.85355 0.146423Z" />
                </svg>
              </a>
            </div>
          </TabPanel>

          {/* coaching */}

          <TabPanel className="relative mt-4">
            <div className="pointer-events-none absolute inset-0 z-10 grid grid-cols-1 gap-2 max-sm:hidden sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4">
              <div className="border-r" />
              <div className="border-x" />
              <div className="border-x" />
              <div className="border-l" />
            </div>
            <ul className="gap-2 sm:grid sm:grid-cols-2 sm:gap-x-5 sm:gap-y-10 md:gap-10 lg:grid-cols-3 xl:grid-cols-4 divide-y border-t">
              {[
                {
                  name: 'GROW Session Simulation',
                  subtitle:
                    'AI conducts a structured coaching call using the GROW model — guiding through Goal, Reality, Options, and Will with reflection capture.',
                  category: 'Inbound',
                  primary_metric:
                    'Completion Rate, Reflection Depth %, Goal Clarity Score',
                  persona_name: 'Coach / L&D Head',
                  url: 'https://example.com/demo/grow-session-call',
                  minute: 6,
                },
                {
                  name: 'Goal Progress Check-In',
                  subtitle:
                    'Outbound AI follow-up to track progress toward goals and capture blockers since the last session.',
                  category: 'Outbound',
                  primary_metric:
                    'Response Rate, Progress Update %, Follow-Up Scheduled %',
                  persona_name: 'Program Manager / Coach',
                  url: 'https://example.com/demo/goal-checkin-call',
                  minute: 4,
                },
                {
                  name: 'Session Reminder & Preparation',
                  subtitle:
                    'AI calls coachees to confirm upcoming session attendance and gather quick prep reflections.',
                  category: 'Outbound',
                  primary_metric:
                    'Confirmation Rate, No-Show %, Preparation Completion %',
                  persona_name: 'Operations / Coaching Coordinator',
                  url: 'https://example.com/demo/session-reminder-call',
                  minute: 3,
                },
                {
                  name: 'Feedback & Reflection Capture',
                  subtitle:
                    'Post-session AI calls capture voice reflections and feedback for program improvement.',
                  category: 'Outbound',
                  primary_metric:
                    'Feedback Completion %, Reflection Quality %, NPS %',
                  persona_name: 'Learning Experience Manager / QA Lead',
                  url: 'https://example.com/demo/feedback-reflection-call',
                  minute: 4,
                },
                {
                  name: 'Coach Availability Booking',
                  subtitle:
                    'Inbound AI helps coachees schedule or reschedule sessions with preferred coaches via voice.',
                  category: 'Inbound',
                  primary_metric: 'Booking Success %, Containment Rate, AHT',
                  persona_name: 'Scheduling Coordinator / Coach Ops',
                  url: 'https://example.com/demo/coach-booking-call',
                  minute: 4,
                },
                {
                  name: 'Post-Program Follow-Up',
                  subtitle:
                    'Outbound AI conducts follow-up conversations after program completion to track sustained growth.',
                  category: 'Outbound',
                  primary_metric:
                    'Response Rate, Goal Retention %, Re-Enrollment %',
                  persona_name: 'Program Lead / L&D Director',
                  url: 'https://example.com/demo/post-program-followup-call',
                  minute: 4,
                },
                {
                  name: 'New Cohort Orientation',
                  subtitle:
                    'Outbound AI introduces new learners to the program, shares resources, and confirms kickoff attendance.',
                  category: 'Outbound',
                  primary_metric:
                    'Engagement Rate, Orientation Attendance %, Resource Access %',
                  persona_name: 'Program Coordinator / L&D Manager',
                  url: 'https://example.com/demo/orientation-call',
                  minute: 4,
                },
                {
                  name: 'Accountability Partner Role-Play',
                  subtitle:
                    'Inbound role-play where AI simulates a supportive accountability partner to reinforce habit change.',
                  category: 'Inbound',
                  primary_metric:
                    'Engagement Duration, Reflection Quality %, Self-Reported Motivation %',
                  persona_name: 'Coach / Behavioral Program Lead',
                  url: 'https://example.com/demo/accountability-roleplay-call',
                  minute: 5,
                },
              ].map((us, index) => (
                <li
                  key={index}
                  className="block transition hover:bg-gray-950/2.5 focus-visible:relative focus-visible:z-20 focus-visible:bg-white/75 focus-visible:backdrop-blur-xs dark:hover:bg-white/5 border-l border-b border-r"
                >
                  <div className="px-4 py-2 sm:px-2">
                    {us.category === 'outbound' ? (
                      <PhoneOutgoing
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    ) : (
                      <PhoneIncoming
                        className=" w-6 h-6 opacity-70 mt-4"
                        strokeWidth={1.5}
                      />
                    )}
                    <div className="flex items-center gap-2 mt-4">
                      <h3 className="text-base/7 font-semibold">
                        <div className="absolute inset-0" />
                        {us.name}
                      </h3>
                    </div>

                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400 line-clamp-2">
                      {us.subtitle}
                    </p>
                    <p className="mt-4 text-sm/6 text-gray-600 md:max-w-2xs dark:text-gray-400">
                      {us.primary_metric}
                    </p>
                  </div>
                  <div className="px-4 py-2 sm:px-2 border-t">
                    <figcaption className="flex space-x-2">
                      <Play
                        className="aspect-square size-9 rounded-full outline -outline-offset-1 outline-gray-950/5 dark:outline-white/10 p-2.5 bg-blue-600/20 fill-blue-600 text-blue-600"
                        strokeWidth={1.5}
                      />
                      <div className="text-sm line-height">
                        <p className="font-medium">{us.persona_name}</p>
                        <p className="text-gray-600 dark:text-gray-400">
                          {us.minute} minutes
                        </p>
                      </div>
                    </figcaption>
                  </div>
                </li>
              ))}
            </ul>
            <div className="pointer-events-none absolute inset-x-0 -bottom-0.5 z-10 flex h-100 items-end justify-center bg-linear-to-b to-white pb-8 max-sm:hidden dark:to-gray-950">
              <a
                className="pointer-events-auto gap-2 inline-flex justify-center rounded-full text-base  font-medium  focus-visible:outline-2 focus-visible:outline-offset-2 bg-gray-950 text-white hover:bg-gray-800 focus-visible:outline-gray-950 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white dark:hover:bg-gray-200 dark:focus-visible:outline-white dark:focus-visible:outline-white px-4 py-2"
                href="/demo"
              >
                Consult an expert
                <svg
                  fill="currentColor"
                  aria-hidden="true"
                  viewBox="0 0 10 10"
                  className="-mr-0.5 w-2.5"
                >
                  <path d="M4.85355 0.146423L9.70711 4.99998L4.85355 9.85353L4.14645 9.14642L7.79289 5.49998H0V4.49998H7.79289L4.14645 0.85353L4.85355 0.146423Z" />
                </svg>
              </a>
            </div>
          </TabPanel>
        </TabPanels>
      </TabGroup>
      <div className="border-y mt-4 px-4 py-2 sm:hidden">
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
      </div>
    </section>
  );
};
