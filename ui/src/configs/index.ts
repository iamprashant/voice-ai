/* eslint-disable import/no-mutable-exports */
import { InputVarType } from '@/models/common';
import { ConnectionConfig } from '@rapidaai/react';

export interface SentryAnalyticsConfig {
  dsn: string;
  tracePropagationTargets: RegExp[];
}

export interface RapidaConfig {
  connection: {
    assistant: string;
    web: string;
    endpoint: string;
  };
  analytics?: SentryAnalyticsConfig;
}

export const getConfig = (): RapidaConfig => {
  if (process.env.NODE_ENV === 'production') {
    return {
      connection: {
        assistant: 'https://assistant-01.rapida.ai',
        web: 'https://api.rapida.ai',
        endpoint: 'https://api.rapida.ai',
      },
      analytics: {
        dsn: 'https://15153cb4befe6a0ae4249f10ff87c0b6@o4506771747831808.ingest.sentry.io/4506771748945920',
        tracePropagationTargets: [/^https:\/\/rapida\.ai\/api/],
      },
    };
  }
  return {
    connection: {
      assistant: 'http://assistant.rapida.local',
      web: 'http://dev.rapida.local',
      endpoint: 'http://dev.rapida.local',
    },
  };
};

export const CONFIG = getConfig();
export const connectionConfig = new ConnectionConfig(CONFIG.connection);

//
export const MAX_PROMPT_MESSAGE_LENGTH = 10;
export const MAX_VAR_KEY_LENGHT = 100;
export const VAR_ITEM_TEMPLATE = {
  name: '',
  type: 'string',
  defaultvalue: '',
};

export const zhRegex = /^[\u4E00-\u9FA5]$/m;
const MAX_ZN_VAR_NAME_LENGHT = 8;
const MAX_EN_VAR_VALUE_LENGHT = 30;
export const getMaxVarNameLength = (value: string) => {
  if (zhRegex.test(value)) return MAX_ZN_VAR_NAME_LENGHT;
  return MAX_EN_VAR_VALUE_LENGHT;
};

export const CONTEXT_PLACEHOLDER_TEXT = '{{#context#}}';
export const HISTORY_PLACEHOLDER_TEXT = '{{#histories#}}';
export const QUERY_PLACEHOLDER_TEXT = '{{#query#}}';
export const PRE_PROMPT_PLACEHOLDER_TEXT = '{{#pre_prompt#}}';

export const SUPPORTED_PROMPT_VARIABLE_TYPE = () => {
  return [
    InputVarType.textInput,
    InputVarType.paragraph,
    InputVarType.number,
    InputVarType.url,
    InputVarType.json,
  ];
};
