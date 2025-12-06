import { Metadata } from '@rapidaai/react';
import { Dropdown } from '@/app/components/dropdown';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { SARVAM_LANGUAGE, SARVAM_SPEECH_TO_TEXT_MODEL } from '@/providers';
export {
  GetSarvamDefaultOptions,
  ValidateSarvamOptions,
} from '@/app/components/providers/speech-to-text/sarvam/constant';

const renderOption = (c: { icon: React.ReactNode; name: string }) => (
  <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
    {c.icon}
    <span className="truncate capitalize">{c.name}</span>
  </span>
);

export const ConfigureSarvamSpeechToText: React.FC<{
  onParameterChange: (parameters: Metadata[]) => void;
  parameters: Metadata[] | null;
}> = ({ onParameterChange, parameters }) => {
  const getParamValue = (key: string) =>
    parameters?.find(p => p.getKey() === key)?.getValue() ?? '';
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
          currentValue={SARVAM_SPEECH_TO_TEXT_MODEL().find(
            x => x.model_id === getParamValue('listen.model'),
          )}
          setValue={(v: { model_id: string }) => {
            updateParameter('listen.model', v.model_id);
          }}
          allValue={SARVAM_SPEECH_TO_TEXT_MODEL()}
          placeholder={`Select model`}
          option={renderOption}
          label={renderOption}
        />
      </FieldSet>
      <FieldSet className="col-span-1">
        <FormLabel>Language</FormLabel>
        <Dropdown
          className="bg-light-background max-w-full dark:bg-gray-950"
          currentValue={SARVAM_LANGUAGE().find(
            x => x.language_id === getParamValue('listen.language'),
          )}
          setValue={(v: { language_id: string }) => {
            updateParameter('listen.language', v.language_id);
          }}
          allValue={SARVAM_LANGUAGE()}
          placeholder={`Select language`}
          option={renderOption}
          label={renderOption}
        />
      </FieldSet>
    </>
  );
};
