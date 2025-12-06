import { Metadata } from '@rapidaai/react';
import { ProviderComponentProps } from '@/app/components/providers';
import {
  ConfigureAssemblyAISpeechToText,
  GetAssemblyAIDefaultOptions,
  ValidateAssemblyAIOptions,
} from '@/app/components/providers/speech-to-text/assemblyai';
import {
  ConfigureAzureSpeechToText,
  GetAzureDefaultOptions,
  ValidateAzureOptions,
} from '@/app/components/providers/speech-to-text/azure';
import {
  ConfigureCartesiaSpeechToText,
  GetCartesiaDefaultOptions,
  ValidateCartesiaOptions,
} from '@/app/components/providers/speech-to-text/cartesia';
import {
  ConfigureDeepgramSpeechToText,
  GetDeepgramDefaultOptions,
  ValidateDeepgramOptions,
} from '@/app/components/providers/speech-to-text/deepgram';
import {
  ConfigureGoogleSpeechToText,
  GetGoogleDefaultOptions,
  ValidateGoogleOptions,
} from '@/app/components/providers/speech-to-text/google';
import {
  ConfigureOpenAISpeechToText,
  GetOpenAIDefaultOptions,
  ValidateOpenAIOptions,
} from '@/app/components/providers/speech-to-text/openai';
import { FC } from 'react';
import {
  ConfigureSarvamSpeechToText,
  GetSarvamDefaultOptions,
  ValidateSarvamOptions,
} from '@/app/components/providers/speech-to-text/sarvam';

export const GetDefaultSpeechToTextIfInvalid = (
  provider: string,
  parameters: Metadata[],
) => {
  switch (provider) {
    case 'google':
    case 'google-cloud':
      return GetGoogleDefaultOptions(parameters);
    case 'deepgram':
      return GetDeepgramDefaultOptions(parameters);
    case 'openai':
      return GetOpenAIDefaultOptions(parameters);
    case 'azure':
    case 'azure-cloud':
      return GetAzureDefaultOptions(parameters);
    case 'assemblyai':
      return GetAssemblyAIDefaultOptions(parameters);
    case 'cartesia':
      return GetCartesiaDefaultOptions(parameters);
    case 'sarvamai':
      return GetSarvamDefaultOptions(parameters);
    default:
      return parameters;
  }
};

export const ValidateSpeechToTextIfInvalid = (
  provider: string,
  parameters: Metadata[],
): boolean => {
  switch (provider) {
    case 'google-cloud':
    case 'google':
      return ValidateGoogleOptions(parameters);
    case 'deepgram':
      return ValidateDeepgramOptions(parameters);
    case 'openai':
      return ValidateOpenAIOptions(parameters);
    case 'azure':
    case 'azure-cloud':
      return ValidateAzureOptions(parameters);
    case 'assemblyai':
      return ValidateAssemblyAIOptions(parameters);
    case 'cartesia':
      return ValidateCartesiaOptions(parameters);
    case 'sarvamai':
      return ValidateSarvamOptions(parameters);
    default:
      return false;
  }
};

/**
 *
 * @returns
 */
export const GetDefaultMicrophoneConfig = (
  existing: Metadata[] = [],
): Metadata[] => {
  const defaultConfig = [
    // end of speech
    { key: 'microphone.eos.timeout', value: '1000' },
    { key: 'microphone.eos.provider', value: 'silence_based_eos' },
    { key: 'microphone.denoising.provider', value: 'rn_noise' },
    { key: 'microphone.vad.provider', value: 'silero_vad' },
    { key: 'microphone.vad.threshold', value: '0.5' },
  ];

  const existingKeys = new Set(existing.map(m => m.getKey()));

  const newConfigs = defaultConfig
    .filter(({ key }) => !existingKeys.has(key))
    .map(({ key, value }) => {
      const metadata = new Metadata();
      metadata.setKey(key);
      metadata.setValue(value);
      return metadata;
    });

  return [...existing, ...newConfigs];
};

export const SpeechToTextConfigComponent: FC<ProviderComponentProps> = ({
  provider,
  parameters,
  onChangeParameter,
}) => {
  switch (provider) {
    case 'google':
    case 'google-cloud':
      return (
        <ConfigureGoogleSpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );

    case 'deepgram':
      return (
        <ConfigureDeepgramSpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );
    case 'openai':
      return (
        <ConfigureOpenAISpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );
    case 'azure':
    case 'azure-cloud':
      return (
        <ConfigureAzureSpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );
    case 'assemblyai':
      return (
        <ConfigureAssemblyAISpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );
    case 'cartesia':
      return (
        <ConfigureCartesiaSpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );
    case 'sarvamai':
      return (
        <ConfigureSarvamSpeechToText
          parameters={parameters}
          onParameterChange={onChangeParameter}
        />
      );
    default:
      return null;
  }
};
