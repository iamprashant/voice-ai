import { Metadata } from '@rapidaai/react';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { DEEPGRAM_VOICE } from '@/providers';
import { ILinkBorderButton } from '@/app/components/form/button';
import { ExternalLink } from 'lucide-react';
import { useState } from 'react';
import { CustomValueDropdown } from '@/app/components/dropdown/custom-value-dropdown';
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
  /**
   *
   */
  const [filteredVoices, setFilteredVoices] = useState(DEEPGRAM_VOICE());

  /**
   *
   * @param key
   * @returns
   */
  const getParamValue = (key: string) =>
    parameters?.find(p => p.getKey() === key)?.getValue() ?? '';

  /**
   *
   * @param key
   * @param value
   */
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

  /**
   *
   */
  return (
    <>
      <FieldSet className="col-span-2">
        <FormLabel>Voice</FormLabel>
        <div className="flex">
          <CustomValueDropdown
            searchable
            className="bg-light-background max-w-full dark:bg-gray-950"
            currentValue={DEEPGRAM_VOICE().find(
              x => x.code === getParamValue('speak.voice.id'),
            )}
            setValue={(v: { code: string }) => {
              updateParameter('speak.voice.id', v.code);
            }}
            allValue={filteredVoices}
            placeholder={`Select voice`}
            option={renderOption}
            label={renderOption}
            customValue
            onSearching={t => {
              const voices = DEEPGRAM_VOICE();
              const v = t.target.value;
              if (v.length > 0) {
                setFilteredVoices(
                  voices.filter(
                    voice =>
                      voice.name.toLowerCase().includes(v.toLowerCase()) ||
                      voice.code?.toLowerCase().includes(v.toLowerCase()) ||
                      voice.age?.toLowerCase().includes(v.toLowerCase()) ||
                      voice.accent?.toLowerCase().includes(v.toLowerCase()) ||
                      voice.gender?.toLowerCase().includes(v.toLowerCase()),
                  ),
                );
                return;
              }
              setFilteredVoices(voices);
            }}
            onAddCustomValue={vl => {
              updateParameter('speak.voice.id', vl);
            }}
          />
          <ILinkBorderButton
            target="_blank"
            href={`/integration/models/deepgram?query=${getParamValue('speak.voice.id')}`}
            className="h-10 text-sm p-2 px-3 bg-light-background max-w-full dark:bg-gray-950 border-b"
          >
            <ExternalLink className="w-4 h-4" strokeWidth={1.5} />
          </ILinkBorderButton>
        </div>
      </FieldSet>
    </>
  );
};
