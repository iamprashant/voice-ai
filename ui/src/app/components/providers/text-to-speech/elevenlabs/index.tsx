import { Metadata } from '@rapidaai/react';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { useState } from 'react';
import {
  ELEVENLABS_LANGUAGE,
  ELEVENLABS_MODEL,
  ELEVENLABS_VOICE,
} from '@/providers';
import { CustomValueDropdown } from '@/app/components/dropdown/custom-value-dropdown';
import { Dropdown } from '@/app/components/dropdown';
import { ILinkBorderButton } from '@/app/components/form/button';
import { ExternalLink } from 'lucide-react';
export {
  GetElevanLabDefaultOptions,
  ValidateElevanLabOptions,
} from './constant';

const renderOption = (c: { name: string }) => {
  return (
    <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
      <span className="truncate capitalize">{c.name}</span>
    </span>
  );
};

export const ConfigureElevanLabTextToSpeech: React.FC<{
  onParameterChange: (parameters: Metadata[]) => void;
  parameters: Metadata[] | null;
}> = ({ onParameterChange, parameters }) => {
  /**
   *
   */
  const [filteredVoices, setFilteredVoices] = useState(ELEVENLABS_VOICE());

  /**
   *
   * @param key
   * @returns
   */
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
      <FieldSet className="col-span-1" key="speak.model">
        <FormLabel>Models</FormLabel>
        <Dropdown
          className="bg-light-background max-w-full dark:bg-gray-950"
          currentValue={ELEVENLABS_MODEL().find(
            x => x.model_id === getParamValue('speak.model'),
          )}
          setValue={(v: { model_id: string }) =>
            updateParameter('speak.model', v.model_id)
          }
          allValue={ELEVENLABS_MODEL()}
          placeholder="Select model"
          option={renderOption}
          label={renderOption}
        />
      </FieldSet>
      <FieldSet className="col-span-1" key="speak.voice.id">
        <FormLabel>Voice</FormLabel>
        <div className="flex">
          <CustomValueDropdown
            searchable
            className="bg-light-background max-w-full dark:bg-gray-950"
            currentValue={filteredVoices.find(
              x => x.voice_id === getParamValue('speak.voice.id'),
            )}
            setValue={(v: { voice_id: string }) =>
              updateParameter('speak.voice.id', v.voice_id)
            }
            allValue={filteredVoices}
            customValue
            onSearching={t => {
              const voices = ELEVENLABS_VOICE();
              const v = t.target.value;
              if (v.length > 0) {
                setFilteredVoices(
                  voices.filter(
                    voice =>
                      voice.name.toLowerCase().includes(v.toLowerCase()) ||
                      voice.voice_id.toLowerCase().includes(v.toLowerCase()) ||
                      voice.category.toLowerCase().includes(v.toLowerCase()) ||
                      voice.labels.language
                        ?.toLowerCase()
                        .includes(v.toLowerCase()),
                  ),
                );
                return;
              }
              setFilteredVoices(voices);
            }}
            onAddCustomValue={vl => {
              updateParameter('speak.voice.id', vl);
            }}
            placeholder="Select voice"
            option={renderOption}
            label={renderOption}
          />
          <ILinkBorderButton
            target="_blank"
            href={`/integration/models/elevenlabs?query=${getParamValue('speak.voice.id')}`}
            className="h-10 text-sm p-2 px-3 bg-light-background max-w-full dark:bg-gray-950 border-b"
          >
            <ExternalLink className="w-4 h-4" strokeWidth={1.5} />
          </ILinkBorderButton>
        </div>
      </FieldSet>

      <FieldSet className="col-span-1" key="speak.language">
        <FormLabel>Language</FormLabel>
        <Dropdown
          className="bg-light-background max-w-full dark:bg-gray-950"
          currentValue={ELEVENLABS_LANGUAGE().find(
            x => x.language_id === getParamValue('speak.language'),
          )}
          setValue={(v: { language_id: string }) =>
            updateParameter('speak.language', v.language_id)
          }
          allValue={ELEVENLABS_LANGUAGE()}
          placeholder="Select language"
          option={renderOption}
          label={renderOption}
        />
      </FieldSet>
    </>
  );
};
