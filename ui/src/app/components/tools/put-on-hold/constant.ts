import { Metadata } from '@rapidaai/react';
import { SetMetadata } from '@/utils/metadata';

// ============================================================================
// Constants
// ============================================================================

const REQUIRED_KEYS = ['tool.max_hold_time'];
const DEFAULT_MAX_HOLD_TIME = '5';
const MIN_HOLD_TIME = 1;
const MAX_HOLD_TIME = 10;

// ============================================================================
// Default Options
// ============================================================================

export const GetPutOnHoldDefaultOptions = (current: Metadata[]): Metadata[] => {
  const metadata: Metadata[] = [];

  const meta = SetMetadata(current, 'tool.max_hold_time', DEFAULT_MAX_HOLD_TIME);
  if (meta) metadata.push(meta);

  return metadata.filter(m => REQUIRED_KEYS.includes(m.getKey()));
};

// ============================================================================
// Validation
// ============================================================================

export const ValidatePutOnHoldDefaultOptions = (
  options: Metadata[],
): string | undefined => {
  const maxHoldTime = options
    .find(m => m.getKey() === 'tool.max_hold_time')
    ?.getValue();

  if (maxHoldTime) {
    const holdTime = parseInt(maxHoldTime, 10);
    if (isNaN(holdTime) || holdTime < MIN_HOLD_TIME || holdTime > MAX_HOLD_TIME) {
      return `Please provide a valid tool.max_hold_time value. It must be a number between ${MIN_HOLD_TIME} and ${MAX_HOLD_TIME} seconds.`;
    }
  }

  return undefined;
};
