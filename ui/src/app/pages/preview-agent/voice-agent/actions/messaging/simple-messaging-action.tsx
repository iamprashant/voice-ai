import {
  Assistant,
  useConnectAgent,
  useInputModeToggleAgent,
  VoiceAgent,
} from '@rapidaai/react';
import { AudioLines, Loader2, Send, StopCircleIcon } from 'lucide-react';
import { FC, HTMLAttributes } from 'react';
import { useForm } from 'react-hook-form';
import { cn } from '@/utils';
import { ScalableTextarea } from '@/app/components/form/textarea';

interface SimpleMessagingAcitonProps extends HTMLAttributes<HTMLDivElement> {
  placeholder?: string;
  voiceAgent: VoiceAgent;
  assistant: Assistant | null;
}
export const SimpleMessagingAction: FC<SimpleMessagingAcitonProps> = ({
  className,
  voiceAgent,
  assistant,
  placeholder,
}) => {
  const { handleVoiceToggle } = useInputModeToggleAgent(voiceAgent);
  const {
    handleConnectAgent,
    handleDisconnectAgent,
    isConnected,
    isConnecting,
  } = useConnectAgent(voiceAgent);

  const {
    register,
    handleSubmit,
    reset,
    formState: { isValid },
  } = useForm({
    mode: 'onChange',
  });

  const onSubmitForm = data => {
    voiceAgent?.onSendText(data.message);
    reset();
  };

  return (
    <div>
      <form
        className={cn(
          'relative flex items-center gap-4 focus-within:border-primary  dark:border-gray-700 bg-light-background focus-within:bg-white',
        )}
        onSubmit={handleSubmit(onSubmitForm)}
      >
        <ScalableTextarea
          placeholder={placeholder}
          wrapperClassName="bg-light-background p-0"
          className="bg-light-background focus-within:bg-white"
          {...register('message', {
            required: 'Please write your message.',
          })}
          required
          onKeyDown={(e: React.KeyboardEvent<HTMLTextAreaElement>) => {
            if (e.key === 'Enter' && !e.shiftKey) {
              handleSubmit(onSubmitForm)(e);
            }
          }}
        />

        <div className="absolute rounded-b-lg right-2 bottom-2 my-auto w-fit">
          <div className="flex flex-row border divide-x">
            {isValid ? (
              <button
                aria-label="Starting Voice"
                type="submit"
                className="group h-9 px-3 flex flex-row items-center justify-center transition-all duration-300 hover:opacity-80 overflow-hidden w-fit bg-blue-600 dark:bg-blue-500 text-white"
              >
                <Send className="w-4.5 h-4.5 flex-shrink-0" strokeWidth={1.5} />
                <span className="text-sm overflow-hidden ml-2 font-medium">
                  Send
                </span>
              </button>
            ) : (
              <button
                aria-label="Starting Voice"
                type="button"
                disabled={isConnecting}
                onClick={async () => {
                  await handleVoiceToggle();
                  !isConnected && (await handleConnectAgent());
                }}
                className="group h-9 px-3 flex flex-row items-center justify-center transition-all duration-300 hover:opacity-80 overflow-hidden w-fit bg-blue-600 dark:bg-blue-500 text-white"
              >
                {isConnecting ? (
                  <Loader2
                    className="w-4.5 h-4.5 flex-shrink-0 animate-spin"
                    strokeWidth={1.5}
                  />
                ) : (
                  <AudioLines
                    className="w-4.5 h-4.5 flex-shrink-0"
                    strokeWidth={1.5}
                  />
                )}
                <span className="text-sm overflow-hidden ml-2 font-medium">
                  {isConnecting ? 'Connecting...' : 'Voice'}
                </span>
              </button>
            )}
            {(isConnected || isConnecting) && (
              <button
                aria-label="Stoping Voice"
                type="button"
                disabled={!isConnected && !isConnecting}
                onClick={async () => {
                  await handleDisconnectAgent();
                }}
                className="group h-9 px-3 flex flex-row items-center justify-center transition-all duration-300 hover:opacity-80 overflow-hidden w-fit bg-red-500 text-white"
              >
                <StopCircleIcon className="w-4 h-4 !border-white" />
                <span className="max-w-0 group-hover:max-w-xs transition-all duration-200 origin-left scale-x-0 group-hover:scale-x-100 group-hover:opacity-100 opacity-0 whitespace-nowrap text-sm overflow-hidden group-hover:ml-2 font-medium">
                  Stop
                </span>
              </button>
            )}
          </div>
        </div>
      </form>
    </div>
  );
};
