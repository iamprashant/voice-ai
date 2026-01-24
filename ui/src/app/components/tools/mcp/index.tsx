import { FC } from 'react';
import { cn } from '@/utils';
import { FormLabel } from '@/app/components/form-label';
import { FieldSet } from '@/app/components/form/fieldset';
import { Input } from '@/app/components/form/input';
import { Textarea } from '@/app/components/form/textarea';
import { InputGroup } from '@/app/components/input-group';
import { ConfigureToolProps, useParameterManager } from '../common';
import { BlueNoticeBlock } from '@/app/components/container/message/notice-block';
import { Select } from '@/app/components/form/select';
import { APiStringHeader } from '@/app/components/external-api/api-header';
import { MCP_PROTOCOL_OPTIONS } from './constant';

// ============================================================================
// Main Component
// ============================================================================

export const ConfigureMCP: FC<ConfigureToolProps> = ({
  toolDefinition,
  onChangeToolDefinition,
  onParameterChange,
  parameters,
  inputClass,
}) => {
  const { getParamValue, updateParameter } = useParameterManager(
    parameters,
    onParameterChange,
  );

  const serverUrl = getParamValue('mcp.server_url');
  const protocol = getParamValue('mcp.protocol') || 'sse';
  const timeout = getParamValue('mcp.timeout') || '30';
  const headers = getParamValue('mcp.headers');

  const handleChange = (field: 'name' | 'description', value: string) => {
    if (toolDefinition && onChangeToolDefinition) {
      onChangeToolDefinition({ ...toolDefinition, [field]: value });
    }
  };

  return (
    <>
      <InputGroup title="MCP Tool Configuration">
        <div className="flex flex-col gap-6 max-w-6xl">
          <FieldSet>
            <FormLabel>Name</FormLabel>
            <Input
              className={cn('bg-light-background', inputClass)}
              value={toolDefinition?.name || ''}
              onChange={e => handleChange('name', e.target.value)}
              placeholder="Enter MCP tool name"
            />
          </FieldSet>

          <FieldSet>
            <FormLabel>Description</FormLabel>
            <Textarea
              className={cn('bg-light-background', inputClass)}
              value={toolDefinition?.description || ''}
              onChange={e => handleChange('description', e.target.value)}
              placeholder="A tool description or definition of when this MCP tool will get triggered."
              rows={2}
            />
          </FieldSet>

          <FieldSet>
            <FormLabel>MCP Server URL</FormLabel>
            <Input
              className={cn('bg-light-background', inputClass)}
              value={serverUrl}
              onChange={e => updateParameter('mcp.server_url', e.target.value)}
              placeholder="https://your-mcp-server.com"
              type="url"
            />
          </FieldSet>

          <div className="grid grid-cols-2 gap-4">
            <FieldSet>
              <FormLabel>Protocol</FormLabel>
              <Select
                className={cn('bg-light-background', inputClass)}
                value={protocol}
                onChange={e => updateParameter('mcp.protocol', e.target.value)}
                options={MCP_PROTOCOL_OPTIONS}
                placeholder="Select protocol"
              />
            </FieldSet>

            <FieldSet>
              <FormLabel>Timeout (seconds)</FormLabel>
              <Input
                className={cn('bg-light-background', inputClass)}
                value={timeout}
                onChange={e => updateParameter('mcp.timeout', e.target.value)}
                placeholder="30"
                type="number"
                min="1"
                max="300"
              />
            </FieldSet>
          </div>

          <FieldSet>
            <FormLabel>Headers</FormLabel>
            <APiStringHeader
              inputClass={inputClass}
              headerValue={headers}
              setHeaderValue={value => updateParameter('mcp.headers', value)}
            />
          </FieldSet>

          <BlueNoticeBlock>
            <div className="text-sm text-blue-900 dark:text-blue-100">
              <div className="text-blue-700 dark:text-blue-300">
                This tool will proxy calls to the specified MCP server. If you
                provide a specific MCP Tool Name, it will call that tool on the
                server; otherwise, it will use the tool name specified above.
                The LLM will see the name and description you provide above.
              </div>
            </div>
          </BlueNoticeBlock>
        </div>
      </InputGroup>
    </>
  );
};
