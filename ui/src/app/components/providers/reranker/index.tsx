import { Metadata } from '@rapidaai/react';
import { Dropdown } from '@/app/components/dropdown';
import { ProviderConfig } from '@/app/components/providers';
import { ConfigureCohereRerankerModel } from '@/app/components/providers/reranker/cohere';
import { GetCohereRerankerDefaultOptions } from '@/app/components/providers/reranker/cohere/constants';
import { cn } from '@/utils';
import { FC } from 'react';
import { RERANKER_PROVIDER } from '@/providers';

export const GetDefaultRerankerConfigIfInvalid = (
  provider: string,
  parameters: Metadata[],
): Metadata[] => {
  switch (provider) {
    case 'cohere':
      return GetCohereRerankerDefaultOptions(parameters);
    default:
      return parameters;
  }
};

export const RerankerConfigComponent: FC<{
  inputClass?: string;
  config: ProviderConfig;
  updateConfig: (config: Partial<ProviderConfig>) => void;
  disabled?: boolean;
}> = ({ config, updateConfig, disabled, inputClass }) => {
  switch (config.provider) {
    case 'cohere':
      return (
        <ConfigureCohereRerankerModel
          inputClass={inputClass}
          parameters={config.parameters}
          onParameterChange={(params: Metadata[]) =>
            updateConfig({ parameters: params })
          }
          disabled={disabled}
        />
      );
    default:
      return null;
  }
};

export const RerankerProvider: React.FC<{
  inputClass?: string;
  onChangeProvider: (i: string, v: string) => void;
  onChangeConfig: (config: ProviderConfig) => void;
  config: ProviderConfig;
  disabled?: boolean;
}> = ({ onChangeProvider, onChangeConfig, config, disabled, inputClass }) => {
  const updateConfig = (newConfig: Partial<ProviderConfig>) => {
    onChangeConfig({ ...config, ...newConfig } as ProviderConfig);
  };
  return (
    <div
      className={cn(
        'p-px',
        'outline-solid outline-transparent',
        'focus-within:outline-blue-600 focus:outline-blue-600 -outline-offset-1',
        'border-b border-gray-300 dark:border-gray-700',
        'dark:focus-within:border-blue-600 focus-within:border-blue-600',
        'transition-all duration-200 ease-in-out',
        'flex relative',
      )}
    >
      <div className="w-44 relative">
        <Dropdown
          disable={disabled}
          className={cn(
            'bg-light-background max-w-full dark:bg-gray-950 focus-within:border-none! focus-within:outline-hidden! border-none! outline-hidden',
            inputClass,
          )}
          currentValue={RERANKER_PROVIDER.find(x => x.code === config.provider)}
          setValue={v => {
            onChangeProvider(v.id, v.name);
          }}
          allValue={RERANKER_PROVIDER}
          placeholder="Select provider"
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
      </div>
      {/*  */}
      <RerankerConfigComponent
        inputClass={inputClass}
        config={config}
        updateConfig={updateConfig}
        disabled={disabled}
      />
    </div>
  );
};
