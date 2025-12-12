import { Metadata } from '@rapidaai/react';
import { SetMetadata } from '@/utils/metadata';

export const GetAPIRequestDefaultOptions = (
  current: Metadata[],
): Metadata[] => {
  const mtds: Metadata[] = [];

  const keysToKeep = [
    'tool.method',
    'tool.endpoint',
    'tool.headers',
    'tool.parameters',
  ];

  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };

  addMetadata('tool.method', 'POST');
  addMetadata('tool.endpoint');
  addMetadata('tool.headers');
  addMetadata('tool.parameters');
  return mtds.filter(m => keysToKeep.includes(m.getKey()));
};

export const ValidateAPIRequestDefaultOptions = (
  options: Metadata[],
): string | undefined => {
  const requiredKeys = ['tool.method', 'tool.endpoint', 'tool.parameters'];
  const validMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'];

  // Check if all required keys are present
  const hasAllRequiredKeys = requiredKeys.every(key =>
    options.some(option => option.getKey() === key),
  );

  if (!hasAllRequiredKeys) {
    return 'Please provide all required metadata keys: tool.method, tool.endpoint, and tool.parameters.';
  }

  // Validate method
  const methodOption = options.find(
    option => option.getKey() === 'tool.method',
  );
  if (
    methodOption &&
    !validMethods.includes(methodOption.getValue().toUpperCase())
  ) {
    return `Please provide HTTP method provided. Supported methods are: ${validMethods.join(', ')}.`;
  }

  // Validate endpoint (check for valid URL)
  const endpointOption = options.find(
    option => option.getKey() === 'tool.endpoint',
  );
  if (endpointOption) {
    try {
      new URL(endpointOption.getValue());
    } catch (error) {
      return 'Please provide a valid URL for the endpoint.';
    }
  }

  // Validate headers
  const headersOption = options.find(
    option => option.getKey() === 'tool.headers',
  );
  if (headersOption) {
    if (
      !headersOption.getValue() ||
      typeof headersOption.getValue() !== 'string'
    ) {
      return 'Please provide valid headers as a string for creating the API request tool.';
    }
    try {
      const headers = JSON.parse(headersOption.getValue());
      for (const [key, value] of Object.entries(headers)) {
        if (
          typeof key !== 'string' ||
          typeof value !== 'string' ||
          key.trim() === '' ||
          value.trim() === ''
        ) {
          return `Please provide a valid header entry detected. Header key and value must be non-empty strings. Key: ${key}, Value: ${value}.`;
        }
      }
    } catch (error) {
      return 'Please provide valid headers.';
    }
  }

  // Validate parameters
  const parametersOption = options.find(
    option => option.getKey() === 'tool.parameters',
  );
  const value = parametersOption?.getValue();
  if (typeof value !== 'string' || value === '') {
    return 'Please provide valid parameters as a non-empty string.';
  }

  try {
    const parameters = JSON.parse(value);
    if (
      typeof parameters !== 'object' ||
      parameters === null ||
      Array.isArray(parameters)
    ) {
      return 'Parameters must be a valid JSON object.';
    }

    const entries = Object.entries(parameters);

    if (entries.length === 0) {
      return 'Parameters object must contain at least one key-value pair.';
    }

    for (const [paramKey, paramValue] of entries) {
      const [type, key] = paramKey.split('.');
      if (
        !type ||
        !key ||
        typeof paramValue !== 'string' ||
        paramValue === ''
      ) {
        return `Please provide a valid parameter format. Key: ${paramKey}, Value: ${paramValue}. Ensure key is in "type.key" format and value is a non-empty string.`;
      }
    }

    const values = entries.map(([, value]) => value);
    const uniqueValues = new Set(values);
    if (values.length !== uniqueValues.size) {
      return 'Please provide a valid parameter, values must be unique.';
    }
  } catch (e) {
    return 'Please provide valid parameters, must be a valid JSON object.';
  }

  // Return undefined if all validations pass successfully
  return undefined;
};
