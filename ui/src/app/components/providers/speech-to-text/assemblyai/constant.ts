import {
  ASSEMBLYAI_LANGUAGE,
  ASSEMBLYAI_SPEECH_TO_TEXT_MODEL,
} from '@/providers';
import { SetMetadata } from '@/utils/metadata';
import { Metadata } from '@rapidaai/react';

export const GetAssemblyAIDefaultOptions = (
  current: Metadata[],
): Metadata[] => {
  const mtds: Metadata[] = [];

  // Define the keys we want to keep
  const keysToKeep = [
    'rapida.credential_id',
    'listen.language',
    'listen.model',
    'listen.threshold',
  ];

  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };

  // Set language
  addMetadata('listen.language', 'en', value =>
    ASSEMBLYAI_LANGUAGE().some(l => l.code === value),
  );

  // Set model
  addMetadata('listen.model', 'slam-1', value =>
    ASSEMBLYAI_SPEECH_TO_TEXT_MODEL().some(m => m.id === value),
  );

  // Set threshold
  addMetadata('listen.threshold', '0.5');
  addMetadata('rapida.credential_id');

  // Only return metadata for the keys we want to keep
  return [
    ...mtds.filter(m => keysToKeep.includes(m.getKey())),
    ...current.filter(m => m.getKey().startsWith('microphone.')),
  ];
};

export const ValidateAssemblyAIOptions = (
  options: Metadata[],
): string | undefined => {
  const credentialID = options.find(
    opt => opt.getKey() === 'rapida.credential_id',
  );
  if (
    !credentialID ||
    !credentialID.getValue() ||
    credentialID.getValue().length === 0
  ) {
    return 'Please provide valid assembly ai credential for speech to text.';
  }
  // Validate language
  const languageOption = options.find(
    opt => opt.getKey() === 'listen.language',
  );
  if (
    !languageOption ||
    !ASSEMBLYAI_LANGUAGE().some(lang => lang.code === languageOption.getValue())
  ) {
    return 'Please provide valid language for speech to text.';
  }

  // Validate model
  const modelOption = options.find(opt => opt.getKey() === 'listen.model');
  if (
    !modelOption ||
    !ASSEMBLYAI_SPEECH_TO_TEXT_MODEL().some(
      m => m.id === modelOption.getValue(),
    )
  ) {
    return 'Please provide valid model for speech to text.';
  }

  return undefined;
};
