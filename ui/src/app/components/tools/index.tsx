import { Metadata } from '@rapidaai/react';
import { FC, useCallback, useMemo } from 'react';
import { Dropdown } from '@/app/components/dropdown';
import { FieldSet } from '@/app/components/form/fieldset';
import { FormLabel } from '@/app/components/form-label';
import { cn } from '@/utils';
import { ConfigureAPIRequest } from '@/app/components/tools/api-request';
import {
  GetAPIRequestDefaultOptions,
  ValidateAPIRequestDefaultOptions,
} from '@/app/components/tools/api-request/constant';
import { ConfigureEndOfConversation } from '@/app/components/tools/end-of-conversation';
import {
  GetEndOfConversationDefaultOptions,
  ValidateEndOfConversationDefaultOptions,
} from '@/app/components/tools/end-of-conversation/constant';
import { ConfigureEndpoint } from '@/app/components/tools/endpoint';
import {
  GetEndpointDefaultOptions,
  ValidateEndpointDefaultOptions,
} from '@/app/components/tools/endpoint/constant';
import { ConfigureKnowledgeRetrieval } from '@/app/components/tools/knowledge-retrieval';
import {
  GetKnowledgeRetrievalDefaultOptions,
  ValidateKnowledgeRetrievalDefaultOptions,
} from '@/app/components/tools/knowledge-retrieval/constant';
import { ConfigurePutOnHold } from '@/app/components/tools/put-on-hold';
import {
  GetPutOnHoldDefaultOptions,
  ValidatePutOnHoldDefaultOptions,
} from '@/app/components/tools/put-on-hold/constant';
import {
  APIRequestToolDefintion,
  BUILDIN_TOOLS,
  EndOfConverstaionToolDefintion,
  EndpointToolDefintion,
  PutOnHoldToolDefintion,
  KnowledgeRetrievalToolDefintion,
} from '@/llm-tools';

// ============================================================================
// Types
// ============================================================================

export type ToolCode =
  | 'knowledge_retrieval'
  | 'api_request'
  | 'endpoint'
  | 'put_on_hold'
  | 'end_of_conversation';

export interface ToolDefinition {
  name: string;
  description: string;
  parameters: string;
}

export interface BuildinToolConfig {
  code: string;
  parameters: Metadata[];
}

interface ConfigureToolProps {
  toolDefinition: ToolDefinition;
  onChangeToolDefinition: (value: ToolDefinition) => void;
  parameters: Metadata[] | null;
  inputClass?: string;
  onParameterChange: (params: Metadata[]) => void;
}

// ============================================================================
// Tool Registry - Single source of truth for tool configurations
// ============================================================================

const TOOL_REGISTRY: Record<
  ToolCode,
  {
    definition: ToolDefinition;
    getDefaultOptions: (params: Metadata[]) => Metadata[];
    validateOptions: (params: Metadata[]) => string | undefined;
    Component: FC<ConfigureToolProps>;
  }
> = {
  knowledge_retrieval: {
    definition: KnowledgeRetrievalToolDefintion,
    getDefaultOptions: GetKnowledgeRetrievalDefaultOptions,
    validateOptions: ValidateKnowledgeRetrievalDefaultOptions,
    Component: ConfigureKnowledgeRetrieval,
  },
  api_request: {
    definition: APIRequestToolDefintion,
    getDefaultOptions: GetAPIRequestDefaultOptions,
    validateOptions: ValidateAPIRequestDefaultOptions,
    Component: ConfigureAPIRequest,
  },
  endpoint: {
    definition: EndpointToolDefintion,
    getDefaultOptions: GetEndpointDefaultOptions,
    validateOptions: ValidateEndpointDefaultOptions,
    Component: ConfigureEndpoint,
  },
  put_on_hold: {
    definition: PutOnHoldToolDefintion,
    getDefaultOptions: GetPutOnHoldDefaultOptions,
    validateOptions: ValidatePutOnHoldDefaultOptions,
    Component: ConfigurePutOnHold,
  },
  end_of_conversation: {
    definition: EndOfConverstaionToolDefintion,
    getDefaultOptions: GetEndOfConversationDefaultOptions,
    validateOptions: ValidateEndOfConversationDefaultOptions,
    Component: ConfigureEndOfConversation,
  },
};

const DEFAULT_TOOL_CODE: ToolCode = 'endpoint';

// ============================================================================
// Helper Functions
// ============================================================================

const isValidToolCode = (code: string): code is ToolCode => {
  return code in TOOL_REGISTRY;
};

const getToolConfig = (code: string) => {
  return isValidToolCode(code)
    ? TOOL_REGISTRY[code]
    : TOOL_REGISTRY[DEFAULT_TOOL_CODE];
};

/**
 * Returns the default tool definition for a given tool code.
 * If an existing definition has all required fields, it returns the existing one.
 * This should only be called during initialization, not on every render.
 */
export const GetDefaultToolDefintion = (
  code: string,
  existing?: Partial<ToolDefinition>,
): ToolDefinition => {
  const hasValidExisting =
    existing?.name && existing?.description && existing?.parameters;

  if (hasValidExisting) {
    return existing as ToolDefinition;
  }

  return getToolConfig(code).definition;
};

/**
 * Returns default tool config parameters, merging with existing if valid.
 */
export const GetDefaultToolConfigIfInvalid = (
  code: string,
  parameters: Metadata[],
): Metadata[] => {
  const config = getToolConfig(code);
  return config.getDefaultOptions(parameters);
};

/**
 * Validates tool parameters and returns an error message if invalid.
 */
export const ValidateToolDefaultOptions = (
  code: string,
  parameters: Metadata[],
): string | undefined => {
  if (!isValidToolCode(code)) {
    return undefined;
  }
  return TOOL_REGISTRY[code].validateOptions(parameters);
};

// ============================================================================
// Components
// ============================================================================

const ToolOptionRenderer: FC<{ icon: string; name: string }> = ({
  icon,
  name,
}) => (
  <span className="inline-flex items-center gap-2 sm:gap-2.5 max-w-full text-sm font-medium">
    <img
      alt=""
      loading="lazy"
      width={16}
      height={16}
      className="w-4 h-4 align-middle block shrink-0"
      src={icon}
    />
    <span className="truncate capitalize">{name}</span>
  </span>
);

const ConfigureBuildinTool: FC<{
  toolDefinition: ToolDefinition;
  onChangeToolDefinition: (value: ToolDefinition) => void;
  config: BuildinToolConfig;
  onParameterChange: (params: Metadata[]) => void;
  inputClass?: string;
}> = ({
  config,
  inputClass,
  toolDefinition,
  onChangeToolDefinition,
  onParameterChange,
}) => {
  if (!isValidToolCode(config.code)) {
    return null;
  }

  const { Component } = TOOL_REGISTRY[config.code];

  return (
    <Component
      toolDefinition={toolDefinition}
      onChangeToolDefinition={onChangeToolDefinition}
      parameters={config.parameters}
      inputClass={inputClass}
      onParameterChange={onParameterChange}
    />
  );
};

export const BuildinTool: FC<{
  toolDefinition: ToolDefinition;
  onChangeToolDefinition: (value: ToolDefinition) => void;
  onChangeBuildinTool: (code: string) => void;
  onChangeConfig: (config: BuildinToolConfig) => void;
  inputClass?: string;
  config: BuildinToolConfig;
}> = ({
  toolDefinition,
  onChangeToolDefinition,
  onChangeBuildinTool,
  onChangeConfig,
  config,
  inputClass,
}) => {
  const handleParameterChange = useCallback(
    (params: Metadata[]) => {
      onChangeConfig({ ...config, parameters: params });
    },
    [config, onChangeConfig],
  );

  const currentTool = useMemo(
    () => BUILDIN_TOOLS.find(tool => tool.code === config.code),
    [config.code],
  );

  const renderOption = useCallback(
    (tool: { icon: string; name: string }) => (
      <ToolOptionRenderer icon={tool.icon} name={tool.name} />
    ),
    [],
  );

  return (
    <>
      <div className="p-6">
        <FieldSet>
          <FormLabel>Action</FormLabel>
          <Dropdown
            className={cn('bg-light-background dark:bg-gray-950', inputClass)}
            currentValue={currentTool}
            setValue={tool => onChangeBuildinTool(tool.code)}
            allValue={BUILDIN_TOOLS}
            placeholder="Select provider"
            option={renderOption}
            label={renderOption}
          />
        </FieldSet>
      </div>

      <ConfigureBuildinTool
        toolDefinition={toolDefinition}
        onChangeToolDefinition={onChangeToolDefinition}
        config={config}
        onParameterChange={handleParameterChange}
        inputClass={inputClass}
      />
    </>
  );
};
