import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { Input } from '@/app/components/form/input';
import { Slider } from '@/app/components/form/slider';
import { Textarea } from '@/app/components/form/textarea';
import { InputGroup } from '@/app/components/input-group';
import { InputHelper } from '@/app/components/input-helper';
import { cn } from '@/utils';
import { FC } from 'react';

export interface ExperienceConfig {
  greeting?: string;
  messageOnError?: string;
  idealTimeout?: string;
  idealMessage?: string;
  maxCallDuration?: string;
}

export const ConfigureExperience: FC<{
  experienceConfig: ExperienceConfig;
  setExperienceConfig: (config: ExperienceConfig) => void;
}> = ({ experienceConfig, setExperienceConfig }) => {
  const {
    greeting,
    messageOnError,
    idealTimeout,
    idealMessage,
    maxCallDuration,
  } = experienceConfig;

  /**
   *
   * @param newGreeting
   */
  const onChangeGreeting = (newGreeting: string) => {
    setExperienceConfig({ ...experienceConfig, greeting: newGreeting });
  };

  const onChangeMessageOnError = (newMessageOnError: string) => {
    setExperienceConfig({
      ...experienceConfig,
      messageOnError: newMessageOnError,
    });
  };

  const onChangeIdealMessage = (idealMessage: string) => {
    setExperienceConfig({
      ...experienceConfig,
      idealMessage: idealMessage,
    });
  };

  const onChangeIdealTimeout = (idealTimeout: string) => {
    setExperienceConfig({
      ...experienceConfig,
      idealTimeout: idealTimeout,
    });
  };
  const onChangeMaxCallDuration = (duration: string) => {
    setExperienceConfig({
      ...experienceConfig,
      maxCallDuration: duration,
    });
  };

  return (
    <InputGroup
      title="General Experience"
      className="bg-white dark:bg-gray-900 "
    >
      <div className={cn('flex max-w-3xl')}>
        <div className="flex flex-col space-y-6 w-full">
          <FieldSet>
            <FormLabel>Greeting</FormLabel>
            <Textarea
              row={2}
              className="bg-light-background"
              value={greeting || ''}
              onChange={e => onChangeGreeting(e.target.value)}
              placeholder={
                'Write a custom greeting message. You can use {{variable}} to include dynamic content.'
              }
            />
          </FieldSet>
        </div>
      </div>
      <InputGroup
        title={'Advanced Experience Configuration'}
        className="mx-0 my-0 mt-6"
        initiallyExpanded={false}
      >
        <div className={cn('flex max-w-3xl')}>
          <div className="flex flex-col space-y-6 w-full">
            <FieldSet>
              <FormLabel>Error Message</FormLabel>
              <Input
                className="bg-light-background"
                placeholder="Message that will be send to the user when error occured"
                value={messageOnError || ''}
                onChange={e => onChangeMessageOnError(e.target.value)}
              />
            </FieldSet>
            <FieldSet>
              <FormLabel>Idle Silence Timeout (millisecond)</FormLabel>
              <div className="flex space-x-2 justify-center items-center">
                <Slider
                  min={3000}
                  max={10000}
                  step={500}
                  value={idealTimeout && parseInt(idealTimeout)}
                  onSlide={(v: number) => {
                    onChangeIdealTimeout(v.toString());
                  }}
                />
                <Input
                  className="bg-light-background w-16"
                  value={idealTimeout}
                  onChange={e => {
                    onChangeIdealTimeout(e.target.value);
                  }}
                />
              </div>
              <InputHelper>
                Duration of silence after which Rapida will interrupt the user
                (3000-10000ms).
              </InputHelper>
            </FieldSet>
            <FieldSet className="relative col-span-1">
              <FormLabel>Idle Message</FormLabel>
              <Input
                className="bg-light-background"
                placeholder="Message that the assistant will speak when the user hasn't responded."
                value={idealMessage}
                onChange={e => onChangeIdealMessage(e.target.value)}
              />
              <InputHelper>
                Message that the assistant will speak when the user hasn't
                responded.
              </InputHelper>
            </FieldSet>
            <FieldSet>
              <FormLabel>Maximum Session Duration (millisecond)</FormLabel>
              <div className="flex space-x-2 justify-center items-center">
                <Slider
                  min={5000}
                  max={15000}
                  step={1000}
                  value={maxCallDuration && parseInt(maxCallDuration)}
                  onSlide={(v: number) => {
                    onChangeMaxCallDuration(v.toString());
                  }}
                />
                <Input
                  className="bg-light-background w-16"
                  value={maxCallDuration}
                  onChange={e => {
                    onChangeMaxCallDuration(e.target.value);
                  }}
                />
              </div>
              <InputHelper>
                Maximum Session Duration. Set the time limit for sessions values
                it should be between 5000ms and 15000ms.
              </InputHelper>
            </FieldSet>
          </div>
        </div>
      </InputGroup>
    </InputGroup>
  );
};
