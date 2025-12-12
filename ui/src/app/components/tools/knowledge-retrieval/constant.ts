import { Metadata } from '@rapidaai/react';
import { SetMetadata } from '@/utils/metadata';

export const GetKnowledgeRetrievalDefaultOptions = (
  current: Metadata[],
): Metadata[] => {
  const mtds: Metadata[] = [];

  const keysToKeep = [
    'tool.search_type',
    'tool.knowledge_id',
    'tool.top_k',
    'tool.score_threshold',
  ];

  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };

  addMetadata('tool.search_type', 'hybrid');
  addMetadata('tool.top_k', '5');
  addMetadata('tool.score_threshold', '0.5');
  addMetadata('tool.knowledge_id');
  return mtds.filter(m => keysToKeep.includes(m.getKey()));
};

/**
 *
 * @param options
 * @returns
 */
export const ValidateKnowledgeRetrievalDefaultOptions = (
  options: Metadata[],
): string | undefined => {
  const requiredKeys = [
    'tool.search_type',
    'tool.knowledge_id',
    'tool.top_k',
    'tool.score_threshold',
  ];
  const allowedSearchTypes = ['semantic', 'fullText', 'hybrid'];

  // Check if all required keys are present
  for (const key of requiredKeys) {
    if (!options.some(option => option.getKey() === key)) {
      return `Please provide the required metadata key: ${key}.`;
    }
  }

  for (const option of options) {
    if (
      option.getKey() === 'search_type' &&
      !allowedSearchTypes.includes(option.getValue())
    ) {
      return `Please provide a valid search type value. Accepted values are ${allowedSearchTypes.join(', ')}.`;
    }

    if (option.getKey() === 'top_k') {
      const topK = Number(option.getValue());
      if (isNaN(topK) || topK < 1 || topK > 10) {
        return 'Please provide a valid top_k value. It must be a number between 1 and 10.';
      }
    }

    if (option.getKey() === 'score_threshold') {
      const scoreThreshold = Number(option.getValue());
      if (
        isNaN(scoreThreshold) ||
        scoreThreshold < 0.1 ||
        scoreThreshold > 0.9
      ) {
        return 'Please provide a valid score_threshold value. It must be a number between 0.1 and 0.9.';
      }
    }
  }

  return undefined; // No errors
};
