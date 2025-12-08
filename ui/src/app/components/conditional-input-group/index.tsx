import { cn } from '@/utils';
import { Switch } from '@headlessui/react';
import { AnimatePresence, motion } from 'framer-motion';
import { ChevronDown, Lock, Unlock } from 'lucide-react';
import { FC, HTMLAttributes, useState } from 'react';

interface ConditionalInputGroupProps extends HTMLAttributes<HTMLDivElement> {
  title?: string;
  enable: boolean;
  initiallyExpanded?: boolean;
  onChangeEnable: (b: boolean) => void;
}

export const ConditionalInputGroup: FC<ConditionalInputGroupProps> = ({
  onChangeEnable,
  enable,
  initiallyExpanded = true,
  ...props
}) => {
  const [isExpanded, setIsExpanded] = useState(initiallyExpanded);

  const handleToggleEnable = (enabled: boolean) => {
    onChangeEnable(enabled);
    if (enabled) {
      setIsExpanded(true); // Expand when enable is turned on
    } else {
      setIsExpanded(false); // Force collapse when enable is turned off
    }
  };

  return (
    <section
      {...props}
      className={cn('border m-4 rounded-[2px]', props.className)}
    >
      <div
        onClick={() => {
          if (enable) {
            setIsExpanded(!isExpanded); // Toggle only if enabled
          }
        }}
        className={cn(
          'outline-solid outline-[1.5px] outline-transparent',
          'focus-within:outline-blue-600 focus:outline-blue-600 outline-offset-[-1.5px]',
          !enable && 'rounded-b-[2px] !border-b-0',
          'px-4 group flex justify-between w-full items-center py-3 text-left rounded-t-[2px] border-b hover:bg-white dark:hover:bg-gray-950',
        )}
      >
        <div className="mr-3.5 flex items-center">
          <div className="flex-none font-semibold text-sm/6">{props.title}</div>
        </div>
        <div className="flex space-x-2 items-center justify-center">
          <Switch
            checked={enable}
            onChange={handleToggleEnable} // Use the updated handler
            className={cn(
              enable ? 'bg-blue-600 justify-end' : 'bg-gray-500 justify-start',
              'relative inline-flex shrink-0 cursor-pointer rounded-full items-center border-2 border-transparent transition-all duration-200 ease-in-out focus:outline-hidden focus-visible:ring-2 focus-visible:ring-white/75 delay-200',
              'w-9 h-6',
            )}
          >
            <span className="sr-only">Switch</span>
            <span
              className={cn(
                'pointer-events-none inline-flex items-center justify-center h-4 w-4 transform rounded-full bg-white shadow-lg ring-0 transition-transform',
              )}
            >
              {enable ? (
                <Unlock className="h-3 w-3 text-blue-600" strokeWidth={1.5} />
              ) : (
                <Lock className="h-3 w-3 text-gray-500" strokeWidth={1.5} />
              )}
            </span>
          </Switch>
          <span className="h-7 w-7 flex items-center justify-center rounded-full p-1 bg-gray-200 dark:bg-gray-800 hover:bg-gray-300 dark:hover:bg-gray-800">
            <ChevronDown
              strokeWidth={1.5}
              className={cn(
                'h-full w-full transition-all',
                isExpanded && 'rotate-180',
              )}
            />
          </span>
        </div>
      </div>
      <AnimatePresence>
        <motion.div
          className="p-6"
          initial={{ opacity: 0, height: 0 }}
          animate={{ opacity: 1, height: 'auto' }}
          exit={{ opacity: 0, height: 0 }}
          transition={{ duration: 0.3, ease: 'easeInOut' }}
          style={{ display: enable && isExpanded ? 'block' : 'none' }} // Preserve conditional rendering
        >
          {props.children}
        </motion.div>
      </AnimatePresence>
    </section>
  );
};
