import { FC, HTMLAttributes, useState, useMemo, useCallback } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import {
  ChevronDown,
  CircleFadingPlus,
  MessageSquareText,
  Mic,
  MicOff,
  X,
  Check,
} from 'lucide-react';

import {
  useConnectAgent,
  MultibandAudioVisualizerComponent,
  useMultibandMicrophoneTrackVolume,
  useSelectInputDeviceAgent,
  useInputModeToggleAgent,
  useMuteAgent,
  VoiceAgent,
} from '@rapidaai/react';
import { cn } from '@/utils';
import { Tooltip } from '@/app/components/base/tooltip';
import TooltipPlus from '@/app/components/base/tooltip-plus';

/**
 *
 */
interface AudioMessagingActionProps extends HTMLAttributes<HTMLDivElement> {
  voiceAgent: VoiceAgent;
  placeholder?: string;
}

/**
 *
 * @param param0
 * @returns
 */
export const AudioMessagingAction: FC<AudioMessagingActionProps> = ({
  className,
  voiceAgent,
  placeholder,
}) => {
  // Use higher band count and more responsive settings for better voice visualization
  const localMultibandVolume = useMultibandMicrophoneTrackVolume(
    voiceAgent,
    5, // 5 frequency bands
    0.05, // Lower loPass for more bass frequencies
    0.85, // Higher hiPass for more treble
  );

  const { handleDisconnectAgent, handleConnectAgent, isConnected } =
    useConnectAgent(voiceAgent);

  const { handleTextToggle } = useInputModeToggleAgent(voiceAgent);

  // Mute/unmute functionality
  const { isMuted, handleToggleMute } = useMuteAgent(voiceAgent);

  const { devices, activeDeviceId, setActiveMediaDevice } =
    useSelectInputDeviceAgent({
      voiceAgent: voiceAgent,
      requestPermissions: true,
    });

  // Get the currently selected device label
  const activeDeviceLabel = useMemo(() => {
    const device = devices.find(d => d.deviceId === activeDeviceId);
    if (device) {
      // Truncate long device names
      const label = device.label || 'Unknown Device';
      return label.length > 25 ? label.substring(0, 22) + '...' : label;
    }
    return 'Select Microphone';
  }, [devices, activeDeviceId]);

  // Handle smooth device switching
  const handleDeviceChange = useCallback(
    async (deviceId: string) => {
      if (deviceId !== activeDeviceId) {
        await setActiveMediaDevice(deviceId);
      }
    },
    [activeDeviceId, setActiveMediaDevice],
  );

  // Compute visualizer frequencies - when muted, show flat bars
  const visualizerFrequencies = useMemo(() => {
    if (isMuted) {
      // Show minimal flat bars when muted
      return Array.from({ length: 5 }, () => [0.02]);
    }
    return localMultibandVolume.length > 0
      ? localMultibandVolume
      : Array.from({ length: 5 }, () => [0.02]);
  }, [isMuted, localMultibandVolume]);

  return (
    <div className={cn('relative flex items-center p-2 py-3 gap-4', className)}>
      <div className="flex items-center justify-center w-full">
        {!isConnected ? (
          <button
            onClick={async () => {
              await handleConnectAgent();
            }}
            className={cn(
              'flex items-center gap-1.5 border-[0.5px] border-primary/10 bg-gray-100 dark:bg-gray-950 rounded-[2px] p-1 shadow-lg px-4 py-2',
            )}
          >
            <CircleFadingPlus className="w-4 h-4" strokeWidth={1.5} />
            <span className="font-medium text-sm">Click to talk</span>
          </button>
        ) : (
          <div className="flex items-center gap-1.5 dark:border-gray-900 bg-gray-100 dark:bg-gray-950 rounded-[2px] p-1 border">
            {/* Mute/Unmute Button */}
            <TooltipPlus
              className="bg-white dark:bg-gray-950 border-[0.5px] rounded-[2px] px-0 py-0"
              popupContent={
                <div className="px-3 py-2 text-sm text-gray-600 dark:text-gray-500">
                  {isMuted ? 'Click to unmute' : 'Click to mute'}
                </div>
              }
            >
              <button
                onClick={handleToggleMute}
                className={cn(
                  'rounded-[2px] p-2 transition-all duration-200',
                  isMuted
                    ? 'bg-red-500/20 text-red-500 hover:bg-red-500 hover:text-muted'
                    : 'bg-transparent hover:bg-primary/10 text-gray-500',
                )}
              >
                {isMuted ? (
                  <MicOff className="w-4.5 h-4.5" strokeWidth={1.5} />
                ) : (
                  <Mic className="w-4.5 h-4.5" strokeWidth={1.5} />
                )}
              </button>
            </TooltipPlus>

            {/* Audio Visualizer */}
            <div className="rounded-[2px] px-2 flex items-center gap-2">
              <MultibandAudioVisualizerComponent
                classNames="gap-1"
                state={isMuted ? 'disconnected' : 'listening'}
                barWidth={4}
                minBarHeight={3}
                maxBarHeight={18}
                barColor={cn(
                  'rounded-full transition-all duration-150',
                  isMuted
                    ? 'bg-red-400/50 dark:bg-red-500/50'
                    : 'bg-primary dark:bg-primary opacity-80',
                )}
                frequencies={visualizerFrequencies}
              />

              {/* Device Selector Dropdown */}
              <FlyoutLink
                FlyoutContent={
                  <div className="flex flex-col rounded-[2px] w-[300px] overflow-hidden bg-white dark:bg-gray-950 border-subtle border shadow-lg dark:border-gray-700">
                    <div className="px-3 py-2 border-b border-gray-200 dark:border-gray-700">
                      <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                        Select Microphone
                      </span>
                    </div>
                    <div className="p-1 space-y-0.5 max-h-[200px] overflow-y-auto">
                      {devices.map((device, idx) => {
                        const isActive = activeDeviceId === device.deviceId;
                        return (
                          <button
                            key={device.deviceId || idx}
                            className={cn(
                              'w-full text-left rounded-[2px] px-3 py-2.5 text-sm transition-all duration-150',
                              'flex items-center justify-between gap-2',
                              isActive
                                ? 'bg-primary text-white'
                                : 'hover:bg-primary/10 text-gray-700 dark:text-gray-200',
                            )}
                            onClick={() => handleDeviceChange(device.deviceId)}
                          >
                            <div className="flex items-center gap-2 min-w-0 flex-1">
                              <span className="truncate">
                                {device.label || `Microphone ${idx + 1}`}
                              </span>
                            </div>
                            {isActive && (
                              <Check
                                className="w-4 h-4 shrink-0 text-white"
                                strokeWidth={2}
                              />
                            )}
                          </button>
                        );
                      })}
                    </div>
                  </div>
                }
                className="flex items-center cursor-pointer"
              >
                {isOpen => (
                  <div className="flex items-center gap-1 px-2 py-2 rounded-[2px] hover:bg-primary/10 transition-colors">
                    <span className="text-xs font-medium text-gray-600 dark:text-gray-300 max-w-[120px] truncate">
                      {activeDeviceLabel}
                    </span>
                    <ChevronDown
                      strokeWidth={1.5}
                      className={cn(
                        'w-4 h-4 transition-transform duration-300 text-gray-500 dark:text-gray-400',
                        isOpen ? 'rotate-180' : '',
                      )}
                    />
                  </div>
                )}
              </FlyoutLink>
            </div>

            {/* Text Mode Button */}
            <TooltipPlus
              className="bg-white dark:bg-gray-950 border-[0.5px] rounded-[2px] px-0 py-0"
              popupContent={
                <div className="px-3 py-2 text-sm text-gray-600 dark:text-gray-500">
                  Click for text mode
                </div>
              }
            >
              <button
                onClick={async () => {
                  await handleTextToggle();
                }}
                className="bg-primary/20 backdrop-blur-xl rounded-[2px] p-2 border-[0.2px] border-primary/20 text-primary hover:bg-primary hover:text-white transition-all duration-200"
              >
                <MessageSquareText className="w-4 h-4" strokeWidth="1.5" />
              </button>
            </TooltipPlus>

            {/* End Session Button */}
            <TooltipPlus
              className="bg-white dark:bg-gray-950 border-[0.5px] rounded-[2px] px-0 py-0"
              popupContent={
                <div className="px-3 py-2 text-sm text-gray-600 dark:text-gray-500">
                  Click to end session
                </div>
              }
            >
              <button
                onClick={async () => {
                  await handleDisconnectAgent();
                  await handleTextToggle();
                }}
                className="bg-red-500/20 backdrop-blur-xl rounded-[2px]  p-2 border-[0.2px] border-red-500/20 text-red-500 hover:bg-red-500 hover:text-white cursor-pointer transition-all duration-200"
              >
                <X className="w-4 h-4" strokeWidth="1.5" />
              </button>
            </TooltipPlus>
          </div>
        )}
      </div>
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
