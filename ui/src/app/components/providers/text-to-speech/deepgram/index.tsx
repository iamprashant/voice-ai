import { Metadata } from '@rapidaai/react';
import { Dropdown } from '@/app/components/providers/text-to-speech/deepgram/drop';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { DEEPGRAM_VOICE } from '@/providers';
import { IBorderButton } from '@/app/components/form/button';
import { ExternalLink } from 'lucide-react';
import { useGlobalNavigation } from '@/hooks/use-global-navigator';
export { GetDeepgramDefaultOptions } from './constant';

const renderOption = (c: { icon: React.ReactNode; name: string }) => (
  <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
    {c.icon}
    <span className="truncate capitalize">{c.name}</span>
  </span>
);

export const ConfigureDeepgramTextToSpeech: React.FC<{
  onParameterChange: (parameters: Metadata[]) => void;
  parameters: Metadata[] | null;
}> = ({ onParameterChange, parameters }) => {
  const { goTo } = useGlobalNavigation();
  const getParamValue = (key: string) =>
    parameters?.find(p => p.getKey() === key)?.getValue() ?? '';
  //
  const updateParameter = (key: string, value: string) => {
    const updatedParams = [...(parameters || [])];
    const existingIndex = updatedParams.findIndex(p => p.getKey() === key);
    const newParam = new Metadata();
    newParam.setKey(key);
    newParam.setValue(value);
    if (existingIndex >= 0) {
      updatedParams[existingIndex] = newParam;
    } else {
      updatedParams.push(newParam);
    }
    onParameterChange(updatedParams);
  };

  return (
    <>
      <FieldSet className="col-span-2">
        <FormLabel>Voice</FormLabel>
        <div className="flex">
          <Dropdown
            searchable
            className="bg-light-background max-w-full dark:bg-gray-950"
            currentValue={DEEPGRAM_VOICE().find(
              x => x.code === getParamValue('speak.voice.id'),
            )}
            setValue={(v: { code: string }) => {
              updateParameter('speak.voice.id', v.code);
            }}
            allValue={DEEPGRAM_VOICE()}
            placeholder={`Select voice`}
            option={renderOption}
            label={renderOption}
            customValue
            onAddCustomValue={() => {}}
          />
          <IBorderButton
            onClick={() => {
              goTo(
                `/integration/models/deepgram?params=${getParamValue('speak.voice.id')}`,
              );
            }}
            className="h-10 text-sm rounded-[2px] p-2 px-3"
          >
            <ExternalLink className="w-4 h-4" strokeWidth={1.5} />
          </IBorderButton>
        </div>
      </FieldSet>
    </>
  );
};
