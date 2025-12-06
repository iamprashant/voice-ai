import { SARVAM_LANGUAGE, SARVAM_SPEECH_TO_TEXT_MODEL } from '@/providers';
import { SetMetadata } from '@/utils/metadata';
import { Metadata } from '@rapidaai/react';

/**
 *
 * @param current
 * @returns
 */
export const GetSarvamDefaultOptions = (current: Metadata[]): Metadata[] => {
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
  addMetadata('listen.language', 'en-IN', value =>
    SARVAM_LANGUAGE().some(l => l.language_id === value),
  );

  // Set model
  addMetadata('listen.model', 'saarika:v2.5', value =>
    SARVAM_SPEECH_TO_TEXT_MODEL().some(m => m.model_id === value),
  );

  // Only return metadata for the keys we want to keep
  return [
    ...mtds.filter(m => keysToKeep.includes(m.getKey())),
    ...current.filter(m => m.getKey().startsWith('microphone.')),
  ];
};

/**
 *
 * @param options
 * @returns
 */
export const ValidateSarvamOptions = (options: Metadata[]): boolean => {
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
  // Validate language
  const languageOption = options.find(
    opt => opt.getKey() === 'listen.language',
  );
  if (
    !languageOption ||
    !SARVAM_LANGUAGE().some(
      lang => lang.language_id === languageOption.getValue(),
    )
  ) {
    return false;
  }

  // Validate model
  const modelOption = options.find(opt => opt.getKey() === 'listen.model');
  if (
    !modelOption ||
    !SARVAM_SPEECH_TO_TEXT_MODEL().some(
      m => m.model_id === modelOption.getValue(),
    )
  ) {
    return false;
  }

  return true;
};
