import { Metadata } from '@rapidaai/react';
import { SetMetadata } from '@/utils/metadata';

// ============================================================================
// Constants
// ============================================================================

const REQUIRED_KEYS = [
  'tool.search_type',
  'tool.knowledge_id',
  'tool.top_k',
  'tool.score_threshold',
];

const ALLOWED_SEARCH_TYPES = ['semantic', 'fullText', 'hybrid'];

const DEFAULTS = {
  search_type: 'hybrid',
  top_k: '5',
  score_threshold: '0.5',
} as const;

// ============================================================================
// Default Options
// ============================================================================

export const GetKnowledgeRetrievalDefaultOptions = (
  current: Metadata[],
): Metadata[] => {
  const metadata: Metadata[] = [];

  const addMetadata = (key: string, defaultValue?: string) => {
    const meta = SetMetadata(current, key, defaultValue);
    if (meta) metadata.push(meta);
  };

  addMetadata('tool.search_type', DEFAULTS.search_type);
  addMetadata('tool.top_k', DEFAULTS.top_k);
  addMetadata('tool.score_threshold', DEFAULTS.score_threshold);
  addMetadata('tool.knowledge_id');

  return metadata.filter(m => REQUIRED_KEYS.includes(m.getKey()));
};

// ============================================================================
// Validation
// ============================================================================

const getOptionValue = (
  options: Metadata[],
  key: string,
): string | undefined => {
  return options.find(opt => opt.getKey() === key)?.getValue();
};

const validateRequiredKeys = (options: Metadata[]): string | undefined => {
  for (const key of REQUIRED_KEYS) {
    if (!options.some(opt => opt.getKey() === key)) {
      return `Please provide the required metadata key: ${key}.`;
    }
  }
  return undefined;
};

const validateSearchType = (value: string | undefined): string | undefined => {
  if (value && !ALLOWED_SEARCH_TYPES.includes(value)) {
    return `Please provide a valid search type value. Accepted values are ${ALLOWED_SEARCH_TYPES.join(', ')}.`;
  }
  return undefined;
};

const validateTopK = (value: string | undefined): string | undefined => {
  if (value !== undefined) {
    const topK = Number(value);
    if (isNaN(topK) || topK < 1 || topK > 10) {
      return 'Please provide a valid top_k value. It must be a number between 1 and 10.';
    }
  }
  return undefined;
};

const validateScoreThreshold = (
  value: string | undefined,
): string | undefined => {
  if (value !== undefined) {
    const threshold = Number(value);
    if (isNaN(threshold) || threshold < 0.1 || threshold > 0.9) {
      return 'Please provide a valid score_threshold value. It must be a number between 0.1 and 0.9.';
    }
  }
  return undefined;
};

export const ValidateKnowledgeRetrievalDefaultOptions = (
  options: Metadata[],
): string | undefined => {
  // Run all validations in sequence, return first error
  return (
    validateRequiredKeys(options) ||
    validateSearchType(getOptionValue(options, 'tool.search_type')) ||
    validateTopK(getOptionValue(options, 'tool.top_k')) ||
    validateScoreThreshold(getOptionValue(options, 'tool.score_threshold'))
  );
};
