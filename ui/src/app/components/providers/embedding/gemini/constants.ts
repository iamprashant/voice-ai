import { SetMetadata } from '@/utils/metadata';
import { Metadata } from '@rapidaai/react';
import { GEMINI_EMBEDDING_MODEL } from '@/providers/';

export const GetGeminiEmbeddingDefaultOptions = (
  current: Metadata[],
): Metadata[] => {
  const mtds: Metadata[] = [];
  const keysToKeep = ['rapida.credential_id', 'model.id', 'model.name'];

  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };
  addMetadata('rapida.credential_id');

  addMetadata('model.id', GEMINI_EMBEDDING_MODEL()[0].id, value =>
    GEMINI_EMBEDDING_MODEL().some(model => model.id === value),
  );

  // Add validation for model.name
  addMetadata('model.name', GEMINI_EMBEDDING_MODEL()[0].name, value =>
    GEMINI_EMBEDDING_MODEL().some(model => model.name === value),
  );

  return mtds.filter(m => keysToKeep.includes(m.getKey()));
};

export const ValidateGeminiEmbeddingDefaultOptions = (
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
    return 'Please provide valid credential for gemini.';
  }
  const modelIdOption = options.find(opt => opt.getKey() === 'model.id');
  if (
    !modelIdOption ||
    !GEMINI_EMBEDDING_MODEL().some(model => model.id === modelIdOption.getValue())
  ) {
    return 'Please select a valid embedding model.';
  }

  return undefined;
};
