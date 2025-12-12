import { Metadata } from '@rapidaai/react';
import { SetMetadata } from '@/utils/metadata';

export const GetPutOnHoldDefaultOptions = (current: Metadata[]): Metadata[] => {
  const mtds: Metadata[] = [];

  const keysToKeep = ['tool.max_hold_time'];

  const addMetadata = (
    key: string,
    defaultValue?: string,
    validationFn?: (value: string) => boolean,
  ) => {
    const metadata = SetMetadata(current, key, defaultValue, validationFn);
    if (metadata) mtds.push(metadata);
  };

  addMetadata('tool.max_hold_time', '5');
  return mtds.filter(m => keysToKeep.includes(m.getKey()));
};

/**
 *
 * @param options
 * @returns
 */
export const ValidatePutOnHoldDefaultOptions = (
  options: Metadata[],
): string | undefined => {
  const maxHoldTimeSec = options
    .find(m => m.getKey() === 'tool.max_hold_time')
    ?.getValue();

  if (maxHoldTimeSec) {
    const holdTime = parseInt(maxHoldTimeSec, 10);
    if (isNaN(holdTime) || holdTime < 1 || holdTime > 10) {
      return 'Please provide a valid tool.max_hold_time value. It must be a number between 1 and 10 seconds.';
    }
  }

  return undefined; // No errors
};
