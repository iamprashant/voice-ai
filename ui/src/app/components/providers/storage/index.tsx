import { Metadata, VaultCredential } from '@rapidaai/react';
import { Dropdown } from '@/app/components/dropdown';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import {
  ConfigureAwsStorage,
  ValidateAwsStorageOptions,
} from '@/app/components/providers/storage/aws';
import {
  ConfigureAzureStorage,
  ValidateAzureStorageOptions,
} from '@/app/components/providers/storage/azure';
import {
  ConfigureGoogleStorage,
  ValidateGoogleCloudStorageOptions,
} from '@/app/components/providers/storage/google';
import { cn } from '@/utils';
import { CredentialDropdown } from '@/app/components/dropdown/credential-dropdown';
import { useCallback } from 'react';
import { ProviderComponentProps } from '@/app/components/providers';
import { STORAGE_PROVIDER } from '@/providers';

export const ValidateStorageOptions = (
  provider: string,
  parameters: Metadata[],
): boolean => {
  switch (provider) {
    case 'azure-cloud':
      return ValidateAzureStorageOptions(parameters);
    case 'google-cloud':
      return ValidateGoogleCloudStorageOptions(parameters);
    case 'aws-cloud':
      return ValidateAwsStorageOptions(parameters);
    default:
      return false;
  }
};

/**
 *
 * @param param0
 * @returns
 */
export const CloudStorageProvider: React.FC<ProviderComponentProps> = ({
  parameters,
  provider,
  onChangeParameter,
  onChangeProvider,
}) => {
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

  const renderConfigComponent = () => {
    switch (provider) {
      case 'azure-cloud':
        return (
          <ConfigureAzureStorage
            parameters={parameters || []}
            onParameterChange={(params: Metadata[]) =>
              onChangeParameter(params)
            }
          />
        );
      case 'google-cloud':
        return (
          <ConfigureGoogleStorage
            parameters={parameters || []}
            onParameterChange={(params: Metadata[]) =>
              onChangeParameter(params)
            }
          />
        );

      case 'aws-cloud':
        return (
          <ConfigureAwsStorage
            parameters={parameters || []}
            onParameterChange={(params: Metadata[]) =>
              onChangeParameter(params)
            }
          />
        );

      default:
        return null;
    }
  };
  return (
    <div className={cn('px-6 pb-6 pt-2 flex gap-8 pl-8')}>
      <div className="space-y-6 w-full max-w-6xl">
        <FieldSet className="relative col-span-1">
          <FormLabel>Provider</FormLabel>
          <Dropdown
            className="bg-light-background max-w-full dark:bg-gray-950"
            currentValue={STORAGE_PROVIDER.find(x => x.code === provider)}
            setValue={v => {
              onChangeProvider(v.code);
            }}
            allValue={STORAGE_PROVIDER}
            placeholder="Select storage provide"
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
          {renderConfigComponent()}
        </div>
      </div>
    </div>
  );
};
