export const Footer = () => {
  return (
    <footer className="col-start-3 row-start-2 max-sm:col-span-full max-sm:col-start-1 @container mb-16 grid w-full">
      <div className="border-y mt-16 flex flex-col sm:items-center text-sm/7 text-nowrap text-gray-600 @lg:flex-row dark:text-gray-400 px-4 sm:px-2 py-2 sm:py-0">
        <p className="sm:px-3">
          © {/* */}2025{/* */} Rapida.ai. All rights reserved.
        </p>
        <div className="h-full hidden sm:block w-px bg-gray-200 dark:bg-gray-800" />
        <a
          className="hover:text-gray-950 dark:hover:text-white sm:px-3"
          href="/static/privacy-policy"
        >
          Privacy Policy
        </a>
        <div className="h-full hidden sm:block w-px bg-gray-200 dark:bg-gray-800" />
        <a
          className="hover:text-gray-950 dark:hover:text-white sm:px-3"
          href="/static/privacy-policy"
        >
          Terms and Conditions
        </a>
        <div className="hidden border-t sm:border-t-0 w-fit mt-2 sm:mt-0 self-center border-x  sm:flex">
          <a
            href="https://x.com/rapidaai"
            target="_blank"
            aria-label="Twitter"
            className="hover:border-blue-600 border border-transparent  group flex h-9 w-9 items-center justify-center  transition-colors undefined"
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
          <div className="h-9 w-[1px] bg-gray-200 dark:bg-gray-800 z-10" />
          <a
            href="https://www.linkedin.com/company/rapida-ai"
            aria-label="LinkedIn"
            target="_blank"
            className="hover:border-blue-600 border border-transparent  group flex h-9 w-9 items-center justify-center  transition-colors"
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
          <div className="h-9 w-[1px] bg-gray-200 dark:bg-gray-800 z-10" />
          <a
            href="https://www.youtube.com/@RapidaAI"
            aria-label="YouTube"
            target="_blank"
            className="hover:border-blue-600 border border-transparent  group flex h-9 w-9 items-center justify-center  transition-colors undefined"
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
          <div className="h-9 w-[1px] bg-gray-200 dark:bg-gray-800 z-10" />
          <a
            href="https://github.com/rapidaai"
            aria-label="GitHub"
            className="hover:border-blue-600 border border-transparent group flex h-9 w-9 items-center justify-center transition-colors"
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
    </footer>
  );
};
