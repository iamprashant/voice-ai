import { FC } from 'react';
import { ExternalLink, Info } from 'lucide-react';
import { cn } from '@/utils';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { Input } from '@/app/components/form/input';
import { Select } from '@/app/components/form/select';
import { Textarea } from '@/app/components/form/textarea';
import { CodeEditor } from '@/app/components/form/editor/code-editor';
import { InputGroup } from '@/app/components/input-group';
import { YellowNoticeBlock } from '@/app/components/container/message/notice-block';
import {
  ToolDefinition,
  ParameterType,
  ASSISTANT_KEY_OPTIONS,
  CONVERSATION_KEY_OPTIONS,
  TOOL_KEY_OPTIONS,
} from './types';

// ============================================================================
// Documentation Notice Block
// ============================================================================

interface DocumentationNoticeProps {
  title?: string;
  documentationUrl: string;
}

export const DocumentationNotice: FC<DocumentationNoticeProps> = ({
  title = 'Know more about knowledge tool definition that can be supported by rapida',
  documentationUrl,
}) => (
  <YellowNoticeBlock className="flex items-center -mx-6 -mt-6">
    <Info className="shrink-0 w-4 h-4" />
    <div className="ms-3 text-sm font-medium">{title}</div>
    <a
      target="_blank"
      href={documentationUrl}
      className="h-7 flex items-center font-medium hover:underline ml-auto text-yellow-600"
      rel="noreferrer"
    >
      Read documentation
      <ExternalLink className="shrink-0 w-4 h-4 ml-1.5" strokeWidth={1.5} />
    </a>
  </YellowNoticeBlock>
);

// ============================================================================
// Tool Definition Form
// ============================================================================

interface ToolDefinitionFormProps {
  toolDefinition: ToolDefinition;
  onChangeToolDefinition: (value: ToolDefinition) => void;
  inputClass?: string;
  documentationUrl?: string;
  documentationTitle?: string;
}

export const ToolDefinitionForm: FC<ToolDefinitionFormProps> = ({
  toolDefinition,
  onChangeToolDefinition,
  inputClass,
  documentationUrl = 'https://doc.rapida.ai/assistants/overview',
  documentationTitle,
}) => {
  const handleChange = <K extends keyof ToolDefinition>(
    field: K,
    value: ToolDefinition[K],
  ) => {
    onChangeToolDefinition({ ...toolDefinition, [field]: value });
  };

  return (
    <InputGroup title="Tool Definition">
      <DocumentationNotice
        title={documentationTitle}
        documentationUrl={documentationUrl}
      />
      <div className={cn('mt-4 flex flex-col gap-8 max-w-6xl')}>
        <FieldSet className="relative w-full">
          <FormLabel>Name</FormLabel>
          <Input
            value={toolDefinition.name}
            onChange={e => handleChange('name', e.target.value)}
            placeholder="Enter tool name"
            className={cn('bg-light-background', inputClass)}
          />
        </FieldSet>

        <FieldSet className="relative w-full">
          <FormLabel>Description</FormLabel>
          <Textarea
            value={toolDefinition.description}
            onChange={e => handleChange('description', e.target.value)}
            className={cn('bg-light-background', inputClass)}
            placeholder="A tool description or definition of when this tool will get triggered."
            rows={2}
          />
        </FieldSet>

        <FieldSet className="relative w-full">
          <FormLabel>Parameters</FormLabel>
          <CodeEditor
            placeholder="Provide a tool parameters that will be passed to llm"
            value={toolDefinition.parameters}
            onChange={value => handleChange('parameters', value)}
            className={cn(
              'min-h-40 max-h-dvh bg-light-background dark:bg-gray-950',
              inputClass,
            )}
          />
        </FieldSet>
      </div>
    </InputGroup>
  );
};

// ============================================================================
// Type Key Selector
// ============================================================================

interface TypeKeySelectorProps {
  type: ParameterType;
  value: string;
  onChange: (newValue: string) => void;
  inputClass?: string;
}

export const TypeKeySelector: FC<TypeKeySelectorProps> = ({
  type,
  value,
  onChange,
  inputClass,
}) => {
  const selectClassName = cn('bg-light-background border-none', inputClass);

  switch (type) {
    case 'assistant':
      return (
        <Select
          value={value}
          onChange={e => onChange(e.target.value)}
          className={selectClassName}
          options={[...ASSISTANT_KEY_OPTIONS]}
        />
      );
    case 'conversation':
      return (
        <Select
          value={value}
          onChange={e => onChange(e.target.value)}
          className={selectClassName}
          options={[...CONVERSATION_KEY_OPTIONS]}
        />
      );
    case 'tool':
      return (
        <Select
          value={value}
          onChange={e => onChange(e.target.value)}
          className={selectClassName}
          options={[...TOOL_KEY_OPTIONS]}
        />
      );
    default:
      return (
        <Input
          value={value}
          onChange={e => onChange(e.target.value)}
          placeholder="Key"
          className={cn('bg-light-background w-full border-none', inputClass)}
        />
      );
  }
};
