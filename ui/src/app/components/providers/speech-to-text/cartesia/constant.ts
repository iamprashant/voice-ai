import { CARTESIA_LANGUAGE, CARTESIA_SPEECH_TO_TEXT_MODEL } from '@/providers';
import { SetMetadata } from '@/utils/metadata';
import { Metadata } from '@rapidaai/react';

export const GetCartesiaDefaultOptions = (current: Metadata[]): Metadata[] => {
  const mtds: Metadata[] = [];

  // Define the keys we want to keep
  const keysToKeep = [
    'rapida.credential_id',
    'listen.language',
    'listen.model',
  ];

  // Function to create or update metadata
  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };

  addMetadata('rapida.credential_id');
  // Set language
  addMetadata('listen.language', 'en', value =>
    CARTESIA_LANGUAGE().some(l => l.code === value),
  );

  // Set model
  addMetadata('listen.model', 'ink-whisper', value =>
    CARTESIA_SPEECH_TO_TEXT_MODEL().some(m => m.id === value),
  );

  // Only return metadata for the keys we want to keep
  return [
    ...mtds.filter(m => keysToKeep.includes(m.getKey())),
    ...current.filter(m => m.getKey().startsWith('microphone.')),
  ];
};

export const ValidateCartesiaOptions = (
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
    return 'Please provide valid cartesia credential for speech to text.';
  }
  // Validate language
  const languageOption = options.find(
    opt => opt.getKey() === 'listen.language',
  );
  if (
    !languageOption ||
    !CARTESIA_LANGUAGE().some(lang => lang.code === languageOption.getValue())
  ) {
    return 'Please provide valid cartesia language options for speech to text.';
  }

  // Validate model
  const modelOption = options.find(opt => opt.getKey() === 'listen.model');
  if (
    !modelOption ||
    !CARTESIA_SPEECH_TO_TEXT_MODEL().some(
      model => model.id === modelOption.getValue(),
    )
  ) {
    return 'Please provide valid cartesia model for speech to text.';
  }

  return undefined;
};
