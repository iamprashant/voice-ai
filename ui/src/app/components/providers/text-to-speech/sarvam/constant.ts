import { SARVAM_LANGUAGE, SARVAM_TEXT_TO_SPEECH_MODEL } from '@/providers';
import { SetMetadata } from '@/utils/metadata';
import { Metadata } from '@rapidaai/react';

export const GetSarvamDefaultOptions = (current: Metadata[]): Metadata[] => {
  const mtds: Metadata[] = [];

  // Define the keys we want to keep
  const keysToKeep = [
    'rapida.credential_id',
    'speak.language',
    'speak.voice.id',
    'speak.model',
  ];
  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };
  addMetadata('rapida.credential_id');
  addMetadata('speak.language');
  addMetadata('speak.voice.id');
  addMetadata('speak.model');
  return [
    ...mtds.filter(m => keysToKeep.includes(m.getKey())),
    ...current.filter(m => m.getKey().startsWith('speaker.')),
  ];
};

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

  const voiceID = options.find(opt => opt.getKey() === 'speak.voice.id');
  if (!voiceID || !voiceID.getValue() || voiceID.getValue().length === 0) {
    return false;
  }

  const validations = [
    {
      key: 'speak.language',
      validator: SARVAM_LANGUAGE(),
      field: 'language_id',
    },
    {
      key: 'speak.model',
      validator: SARVAM_TEXT_TO_SPEECH_MODEL(),
      field: 'model_id',
    },
  ];

  return validations.every(({ key, validator, field }) => {
    const option = options.find(opt => opt.getKey() === key);
    return option && validator.some(item => item[field] === option.getValue());
  });
};
