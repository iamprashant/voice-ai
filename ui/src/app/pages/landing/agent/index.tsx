import { TooltipCursor } from '@/app/components/base/tooltip-cursor';
import { ChevronDown, Mic, X } from 'lucide-react';
import { Tooltip } from '@/app/components/base/tooltip';
import {
  useConnectAgent,
  MultibandAudioVisualizerComponent,
  useMultibandMicrophoneTrackVolume,
  useSelectInputDeviceAgent,
  VoiceAgent,
} from '@rapidaai/react';
import { FC, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { cn } from '@/styles/media';
export const DemoVoiceAgent: FC<{
  className?: string;
  placeholder?: string;
  voiceAgent: VoiceAgent;
}> = ({ className, placeholder, voiceAgent }) => {
  const localMultibandVolume = useMultibandMicrophoneTrackVolume(voiceAgent, 5);
  const { handleConnectAgent, handleDisconnectAgent, isConnected } =
    useConnectAgent(voiceAgent);
  const { devices, activeDeviceId, setActiveMediaDevice } =
    useSelectInputDeviceAgent({
      voiceAgent: voiceAgent,
      requestPermissions: true,
    });
  return (
    <div className="flex justify-center items-center flex-col gap-4">
      <TooltipCursor
        content={
          isConnected ? (
            <span>Just talk</span>
          ) : (
            <span className="flex items-center gap-1">
              <Mic className="w-4 h-4" />
              <span>Click to talk.</span>
            </span>
          )
        }
        delay={200}
      >
        <motion.div
          whileHover={{
            scale: 1.05,
          }}
          whileTap={{
            scale: 0.985,
          }}
          className={cn(
            'w-fit cursor-pointer z-10 animate-border-rotate [background:linear-gradient(45deg,#ffffff,--theme(--color-white)_50%,#ffffff)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-slate-600/.48)_80%,--theme(--color-blue-500)_86%,--theme(--color-blue-300)_90%,--theme(--color-blue-500)_94%,--theme(--color-slate-600/.48))_border-box] border border-transparent',
            'dark:[background:linear-gradient(45deg,#000000,--theme(--color-black)_50%,#000000)_padding-box,conic-gradient(from_var(--border-angle),--theme(--color-slate-600/.48)_80%,--theme(--color-blue-500)_86%,--theme(--color-blue-300)_90%,--theme(--color-blue-500)_94%,--theme(--color-slate-600/.48))_border-box]',
            'font-medium rounded-[2px] px-4 h-10 space-x-3 dark:text-white text-gray-800 flex items-center',
            isConnected && 'bg-transparent!',
          )}
          onMouseEnter={e => {
            e.preventDefault();
          }}
          onMouseLeave={e => {
            e.preventDefault();
          }}
          onClick={async () => {
            await handleConnectAgent();
          }}
        >
          {isConnected ? (
            <div className="flex items-center gap-1.5 ">
              <div className="rounded-[2px] px-1 py-1 flex items-center gap-2">
                <Mic className="w-5 h-5 opacity-50 hover:opacity-100 dark:text-white text-gray-800" />
                <MultibandAudioVisualizerComponent
                  classNames="gap-1"
                  state="connecting"
                  barWidth={3}
                  minBarHeight={2}
                  maxBarHeight={14}
                  barColor="bg-gray-700 dark:bg-white opacity-50 "
                  frequencies={
                    localMultibandVolume.length > 0
                      ? localMultibandVolume
                      : Array.from({ length: 3 }, () => [0.01])
                  }
                />
                <FlyoutLink
                  className={'flex items-center'}
                  FlyoutContent={
                    <div className="flex flex-col rounded-[2px] w-[265px] overflow-hidden bg-white dark:bg-gray-950 border-subtle border shadow-lg dark:border-gray-700">
                      <div className="p-1 space-y-1">
                        {devices.map((x, idx) => {
                          return (
                            <div
                              key={idx}
                              className={cn(
                                activeDeviceId === x.deviceId &&
                                  'bg-primary! text-white!',
                                'rounded-[2px] px-2.5 py-2.5 hover:bg-primary/10 hover:text-primary hover:rounded-[2px] text-palette-blue-700 text-sm',
                              )}
                            >
                              <div
                                className="flex gap-2 text-xs font-medium items-center"
                                onClick={() => {
                                  setActiveMediaDevice(x.deviceId);
                                }}
                              >
                                {x.label}
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  }
                >
                  {isOpen => (
                    <ChevronDown
                      className={cn(
                        'w-5 h-5 transition-transform duration-300 opacity-50 hover:opacity-100 dark:text-white',
                        isOpen ? 'rotate-180' : '',
                      )}
                    />
                  )}
                </FlyoutLink>
              </div>

              <Tooltip content="Click to end session">
                <button
                  onClick={async () => {
                    await handleDisconnectAgent();
                  }}
                  className="bg-red-500/20 backdrop-blur-xl rounded-[2px] shadow-lg p-1 border-[0.2px] border-red-500/20 text-red-500 hover:bg-red-500 hover:text-white cursor-pointer"
                >
                  <X className="w-4 h-4" strokeWidth={1.5} />
                </button>
              </Tooltip>
            </div>
          ) : (
            <>
              <span>Talk to a live agent</span>
              <svg
                className="h-4 w-4"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                role="presentation"
              >
                <path
                  d="M7.00171 7C7.00171 4.23858 9.24029 2 12.0017 2C14.7631 2 17.0017 4.23858 17.0017 7V12C17.0017 14.7614 14.7631 17 12.0017 17C9.24029 17 7.00171 14.7614 7.00171 12V7Z"
                  fill="currentcolor"
                ></path>
                <path
                  d="M5.25323 13.867L5.12021 13.385L4.15625 13.6511L4.28928 14.133C5.18021 17.3611 8.04316 19.7714 11.5016 19.9846V22H12.5016V19.9846C15.9601 19.7714 18.8231 17.3611 19.714 14.133L19.847 13.6511L18.8831 13.385L18.75 13.867C17.9331 16.827 15.2205 19 12.0016 19C8.78279 19 6.07019 16.827 5.25323 13.867Z"
                  fill="currentcolor"
                ></path>
              </svg>
            </>
          )}
        </motion.div>
      </TooltipCursor>
    </div>
  );
};

const FlyoutLink = ({ children, FlyoutContent, className }) => {
  const [open, setOpen] = useState(false);
  return (
    <div
      onMouseEnter={() => setOpen(true)}
      onMouseLeave={() => setOpen(false)}
      className="relative w-fit h-fit"
    >
      <div className={cn(className)}>
        {typeof children === 'function' ? children(open) : children}
      </div>
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: 15 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 15 }}
            transition={{ duration: 0.3, ease: 'easeOut' }}
            className="absolute left-0 bottom-12 z-50"
          >
            <div className="absolute -top-6 left-0 right-0 h-6 bg-transparent" />
            <div className="absolute left-4 -bottom-2 h-4 w-4 rotate-45 bg-white dark:bg-gray-950 border-b border-r dark:border-gray-700" />
            {FlyoutContent}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
