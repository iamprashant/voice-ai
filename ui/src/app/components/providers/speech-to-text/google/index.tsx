import { Metadata } from '@rapidaai/react';
import { Dropdown } from '@/app/components/dropdown';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { Input } from '@/app/components/form/input';
import { Slider } from '@/app/components/form/slider';
import { InputHelper } from '@/app/components/input-helper';
import {
  GOOGLE_SPEECH_TO_TEXT_LANGUGAE,
  GOOGLE_SPEECH_TO_TEXT_MODEL,
} from '@/providers/index';
export {
  GetGoogleDefaultOptions,
  ValidateGoogleOptions,
} from '@/app/components/providers/speech-to-text/google/constant';

const renderOption = (c: { name: string }) => (
  <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
    <span className="truncate capitalize">{c.name}</span>
  </span>
);

export const ConfigureGoogleSpeechToText: React.FC<{
  onParameterChange: (parameters: Metadata[]) => void;
  parameters: Metadata[] | null;
}> = ({ onParameterChange, parameters }) => {
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
      <FieldSet className="col-span-1">
        <FormLabel>Model</FormLabel>
        <Dropdown
          className="bg-light-background max-w-full dark:bg-gray-950"
          currentValue={GOOGLE_SPEECH_TO_TEXT_MODEL().find(
            x => x.id === getParamValue('listen.model'),
          )}
          setValue={(v: { id: string }) => {
            updateParameter('listen.model', v.id);
          }}
          allValue={GOOGLE_SPEECH_TO_TEXT_MODEL()}
          placeholder={`Select model`}
          option={renderOption}
          label={renderOption}
        />
      </FieldSet>
      <FieldSet className="col-span-1">
        <FormLabel>Language</FormLabel>
        <Dropdown
          className="bg-light-background max-w-full dark:bg-gray-950"
          currentValue={GOOGLE_SPEECH_TO_TEXT_LANGUGAE().find(
            x => x.code === getParamValue('listen.language'),
          )}
          setValue={(v: { code: string }) => {
            updateParameter('listen.language', v.code);
          }}
          allValue={GOOGLE_SPEECH_TO_TEXT_LANGUGAE()}
          placeholder={`Select language`}
          option={renderOption}
          label={renderOption}
        />
      </FieldSet>
      <FieldSet className="col-span-1">
        <FormLabel>Region</FormLabel>
        <Dropdown
          className="bg-light-background max-w-full dark:bg-gray-950"
          currentValue={['global'].find(
            x => x === getParamValue('listen.region'),
          )}
          setValue={(v: string) => {
            updateParameter('listen.region', v);
          }}
          allValue={['global']}
          placeholder={`Select region`}
          option={(c: string) => (
            <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
              <span className="truncate capitalize">{c}</span>
            </span>
          )}
          label={(c: string) => (
            <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
              <span className="truncate capitalize">{c}</span>
            </span>
          )}
        />
      </FieldSet>
      <FieldSet className="col-span-1">
        <FormLabel>Transcript Confidence Threshold</FormLabel>
        <div className="flex space-x-2 justify-center items-center">
          <Slider
            min={0.1}
            max={0.9}
            step={0.1}
            value={parseFloat(getParamValue('listen.threshold')) || 0.1}
            onSlide={c => {
              updateParameter('listen.threshold', c.toString());
            }}
          />
          <Input
            type="number"
            min={0.1}
            max={0.9}
            step={0.1}
            value={getParamValue('listen.threshold')}
            onChange={v => {
              updateParameter('listen.threshold', v.target.value);
            }}
            className="bg-light-background w-16"
          />
        </div>
        <InputHelper>
          Transcripts with a confidence score below this threshold will be
          filtered out.
        </InputHelper>
      </FieldSet>
    </>
  );
};
