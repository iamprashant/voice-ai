import { Metadata } from '@rapidaai/react';
import { SetMetadata } from '@/utils/metadata';

export const GetEndpointDefaultOptions = (current: Metadata[]): Metadata[] => {
  const mtds: Metadata[] = [];
  const keysToKeep = ['tool.endpoint_id', 'tool.parameters'];

  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };

  addMetadata('tool.endpoint_id');
  addMetadata('tool.parameters');
  return mtds.filter(m => keysToKeep.includes(m.getKey()));
};

export const ValidateEndpointDefaultOptions = (
  options: Metadata[],
): string | undefined => {
  const requiredKeys = ['tool.endpoint_id', 'tool.parameters'];
  const foundKeys = new Set<string>();

  for (const option of options) {
    const key = option.getKey();
    foundKeys.add(key);

    if (key === 'tool.endpoint_id') {
      const value = option.getValue();
      if (typeof value !== 'string' || value === '') {
        return 'Please provide a valid value for tool.endpoint_id. It must be a non-empty string.';
      }
    }

    if (key === 'tool.parameters') {
      const value = option.getValue();
      if (typeof value !== 'string' || value === '') {
        return 'Please provide a valid value for tool.parameters. It must be a non-empty JSON string.';
      }

      try {
        const parameters = JSON.parse(value);
        if (
          typeof parameters !== 'object' ||
          parameters === null ||
          Array.isArray(parameters)
        ) {
          return 'Please ensure tool.parameters is a valid JSON object.';
        }

        const entries = Object.entries(parameters);
        if (entries.length === 0) {
          return 'Please provide parameter values within tool.parameters. It cannot be an empty object.';
        }

        for (const [paramKey, paramValue] of entries) {
          const [type, key] = paramKey.split('.');
          if (
            !type ||
            !key ||
            typeof paramValue !== 'string' ||
            paramValue === ''
          ) {
            return 'Please ensure each parameter key follows the format "type.key" and the values are non-empty strings.';
          }
        }

        const values = entries.map(([, value]) => value);
        const uniqueValues = new Set(values);
        if (values.length !== uniqueValues.size) {
          return 'Please ensure parameter values within tool.parameters are unique.';
        }
      } catch (e) {
        return 'Please provide a valid JSON string for tool.parameters.';
      }
    }
  }

  if (!requiredKeys.every(key => foundKeys.has(key))) {
    return `Please ensure all required metadata keys are present: ${requiredKeys.join(', ')}.`;
  }

  return undefined; // No errors
};
