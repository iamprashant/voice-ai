import { AZURE_LANGUAGE } from '@/providers';
import { SetMetadata } from '@/utils/metadata';
import { Metadata } from '@rapidaai/react';

export const GetAzureDefaultOptions = (current: Metadata[]): Metadata[] => {
  const mtds: Metadata[] = [];
  const keysToKeep = [
    'rapida.credential_id',
    'speak.language',
    'speak.voice.id',
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
  addMetadata('speak.language', 'en-US', value =>
    AZURE_LANGUAGE().some(l => l.code === value),
  );
  addMetadata('speak.voice.id');
  return [
    ...mtds.filter(m => keysToKeep.includes(m.getKey())),
    ...current.filter(m => m.getKey().startsWith('speaker.')),
  ];
};

export const ValidateAzureOptions = (options: Metadata[]): boolean => {
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
  const languageOption = options.find(opt => opt.getKey() === 'speak.language');
  if (
    !languageOption ||
    !AZURE_LANGUAGE().some(lang => lang.code === languageOption.getValue())
  ) {
    return false;
  }

  // Validate voice
  const voiceOption = options.find(opt => opt.getKey() === 'speak.voice.id');
  if (!voiceOption) {
    return false;
  }

  return true;
};
