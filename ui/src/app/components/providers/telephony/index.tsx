import { Metadata, VaultCredential } from '@rapidaai/react';
import {
  ConfigureExotelTelephony,
  ValidateExotelTelephonyOptions,
} from '@/app/components/providers/telephony/exotel';
import {
  ConfigureTwilioTelephony,
  ValidateTwilioTelephonyOptions,
} from '@/app/components/providers/telephony/twilio';
import {
  ConfigureVonageTelephony,
  ValidateVonageTelephonyOptions,
} from '@/app/components/providers/telephony/vonage';
import { Dropdown } from '@/app/components/dropdown';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { InputGroup } from '@/app/components/input-group';
import { InputHelper } from '@/app/components/input-helper';
import { CredentialDropdown } from '@/app/components/dropdown/credential-dropdown';
import { useCallback } from 'react';
import { ProviderComponentProps } from '@/app/components/providers';
import { TELEPHONY_PROVIDER } from '@/providers';

export const ValidateTelephonyOptions = (
  provider: string,
  parameters: Metadata[],
): boolean => {
  switch (provider) {
    case 'vonage':
      return ValidateVonageTelephonyOptions(parameters);
    case 'twilio':
      return ValidateTwilioTelephonyOptions(parameters);
    case 'exotel':
      return ValidateExotelTelephonyOptions(parameters);
    default:
      return false;
  }
};

/**
 *
 * @param param0
 * @returns
 */

export const ConfigureTelephonyComponent: React.FC<ProviderComponentProps> = ({
  provider,
  parameters,
  onChangeParameter,
}) => {
  switch (provider) {
    case 'exotel':
      return (
        <ConfigureExotelTelephony
          parameters={parameters || []}
          onParameterChange={onChangeParameter}
        />
      );
    case 'vonage':
      return (
        <ConfigureVonageTelephony
          parameters={parameters || []}
          onParameterChange={onChangeParameter}
        />
      );

    case 'twilio':
      return (
        <ConfigureTwilioTelephony
          parameters={parameters || []}
          onParameterChange={onChangeParameter}
        />
      );

    default:
      return null;
  }
};

/**
 *
 * @param param0
 * @returns
 */

export const TelephonyProvider: React.FC<ProviderComponentProps> = props => {
  const { provider, onChangeParameter, onChangeProvider, parameters } = props;
  const getParamValue = useCallback(
    (key: string) => {
      return parameters?.find(p => p.getKey() === key)?.getValue() ?? '';
    },
    [JSON.stringify(parameters)],
  );

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
    onChangeParameter(updatedParams);
  };

  return (
    <InputGroup title="Telephony" className="bg-white dark:bg-gray-900 ">
      <div className="flex flex-col space-y-6">
        <FieldSet>
          <FormLabel>Telephony provider</FormLabel>
          <Dropdown
            className="bg-light-background max-w-full dark:bg-gray-950"
            currentValue={TELEPHONY_PROVIDER.find(x => x.code === provider)}
            setValue={v => {
              onChangeProvider(v.code);
            }}
            allValue={TELEPHONY_PROVIDER}
            placeholder="Select telephony provider"
            option={c => {
              return (
                <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
                  <img
                    alt=""
                    loading="lazy"
                    width={16}
                    height={16}
                    className="sm:h-4 sm:w-4 w-4 h-4 align-middle block shrink-0"
                    src={c.image}
                  />
                  <span className="truncate capitalize">{c.name}</span>
                </span>
              );
            }}
            label={c => {
              return (
                <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
                  <img
                    alt=""
                    loading="lazy"
                    width={16}
                    height={16}
                    className="sm:h-4 sm:w-4 w-4 h-4 align-middle block shrink-0"
                    src={c.image}
                  />
                  <span className="truncate capitalize">{c.name}</span>
                </span>
              );
            }}
          />
          <InputHelper>
            Choose a telephony provider to handle voice communication for your
            applications. Each provider offers different capabilities, pricing
            structures, and global coverage.
          </InputHelper>
        </FieldSet>
        {provider && (
          <CredentialDropdown
            className="bg-light-background max-w-full dark:bg-gray-950"
            onChangeCredential={(c: VaultCredential) => {
              updateParameter('rapida.credential_id', c.getId());
            }}
            currentCredential={getParamValue('rapida.credential_id')}
            provider={provider}
          />
        )}
        <div className="grid grid-cols-3 gap-x-6 gap-y-3">
          <ConfigureTelephonyComponent {...props} />
        </div>
      </div>
    </InputGroup>
  );
};
