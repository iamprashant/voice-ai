import { Metadata } from '@rapidaai/react';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { Input } from '@/app/components/form/input';
import { InputHelper } from '@/app/components/input-helper';

export const ValidateSIPTelephonyOptions = (options: Metadata[]): boolean => {
  const credentialID = options.find(
    opt => opt.getKey() === 'rapida.credential_id',
  );
  if (
    !credentialID ||
    !credentialID.getValue() ||
    credentialID.getValue().length === 0
  ) {
    return false;
  }

  // Validate caller ID
  const callerId = options.find(opt => opt.getKey() === 'phone');
  if (!callerId || !callerId.getValue() || callerId.getValue().length === 0) {
    return false;
  }

  return true;
};

export const ConfigureSIPTelephony: React.FC<{
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
      <FieldSet className="col-span-2">
        <FormLabel>Caller ID</FormLabel>
        <Input
          className="bg-light-background"
          value={getParamValue('phone')}
          onChange={v => {
            updateParameter('phone', v.target.value);
          }}
          placeholder="e.g., +15551234567"
        />
        <InputHelper>
          The phone number to display as caller ID for outbound calls.
        </InputHelper>
      </FieldSet>
    </>
  );
};
