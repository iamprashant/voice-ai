import React, { useCallback, useState } from 'react';
import { useRapidaStore } from '@/hooks';
import { useCurrentCredential } from '@/hooks/use-credential';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast/headless';
import { Helmet } from '@/app/components/helmet';
import {
  IBlueBGArrowButton,
  ICancelButton,
  SimpleButton,
} from '@/app/components/form/button';
import { TabForm } from '@/app/components/form/tab-form';
import {
  ConnectionConfig,
  CreateEndpointResponse,
  EndpointAttribute,
  EndpointProviderModelAttribute,
  Metadata,
} from '@rapidaai/react';
import ConfirmDialog from '@/app/components/base/modal/confirm-ui';
import { create_endpoint_success_message } from '@/utils/messages';
import {
  GetDefaultTextProviderConfigIfInvalid,
  TextProvider,
  ValidateTextProviderDefaultOptions,
} from '@/app/components/providers/text';
import { ConfigPrompt } from '@/app/components/configuration/config-prompt';
import { randomMeaningfullName, randomString } from '@/utils';
import { FieldSet } from '@/app/components/form/fieldset';
import { FormLabel } from '@/app/components/form-label';
import { Input } from '@/app/components/form/input';
import { TagInput } from '@/app/components/form/tag-input';
import { EndpointTag } from '@/app/components/form/tag-input/endpoint-tags';
import { Textarea } from '@/app/components/form/textarea';
import { CreateEndpoint } from '@rapidaai/react';
import { ServiceError } from '@rapidaai/react';
import { ChatCompletePrompt } from '@/utils/prompt';
import { connectionConfig } from '@/configs';
import {
  BlueNoticeBlock,
  YellowNoticeBlock,
} from '@/app/components/container/message/notice-block';
import { ArrowUpRight, ExternalLink, Info, MoveRight } from 'lucide-react';
import { ConfigureEndpointPromptDialog } from '@/app/components/base/modal/configure-endpoint-prompt-modal/index';

/**
 *
 * @param props
 * @returns
 */
export function CreateEndpointPage() {
  /**
   * authentication
   */
  const { authId, token, projectId } = useCurrentCredential();

  /**
   * multistep form
   */
  const [activeTab, setActiveTab] = useState('choose-model');

  /**
   * error
   */
  const [errorMessage, setErrorMessage] = useState('');

  /**
   * global loader
   */
  const { loading, showLoader, hideLoader } = useRapidaStore();

  /**
   * global navigator
   */
  const navigator = useNavigate();
  /**
   * form element
   */
  const [name, setName] = useState<string>(randomMeaningfullName('endpoint'));
  const [description, setDescription] = useState('');
  const [tags, setTags] = useState<string[]>([]);
  const [promptConfig, setPromptConfig] = useState<{
    prompt: { role: string; content: string }[];
    variables: { name: string; type: string; defaultvalue: string }[];
  }>({
    prompt: [{ role: 'system', content: '' }],
    variables: [],
  });

  const [textProviderModel, setTextProviderModel] = useState<{
    provider: string;
    parameters: Metadata[];
  }>({
    provider: 'azure',
    parameters: GetDefaultTextProviderConfigIfInvalid('azure', []),
  });

  const onChangeTextProvider = (providerName: string) => {
    setTextProviderModel({
      provider: providerName,
      parameters: GetDefaultTextProviderConfigIfInvalid(
        providerName,
        textProviderModel.parameters,
      ),
    });
  };

  const onChangeTextProviderParameter = (parameters: Metadata[]) => {
    setTextProviderModel({ ...textProviderModel, parameters });
  };

  const onAddTag = (newTag: string) => {
    setTags(prevTags => [...prevTags, newTag]);
  };

  const onRemoveTag = (tagToRemove: string) => {
    setTags(prevTags => prevTags.filter(tag => tag !== tagToRemove));
  };
  /**
   *
   */
  const afterCreateEndpoint = useCallback(
    (err: ServiceError | null, response: CreateEndpointResponse | null) => {
      hideLoader();
      if (err) {
        setErrorMessage('Something went wrong, Please try again in sometime.');
        return;
      }
      if (response?.getSuccess() && response.getData()) {
        let ep = response.getData();
        toast.success(create_endpoint_success_message(name));
        navigator(`/deployment/endpoint/${ep?.getId()}`);
        return;
      }
      if (response?.getError()) {
        let err = response.getError();
        const message = err?.getHumanmessage();
        if (message) {
          setErrorMessage(message);
          return;
        }
        setErrorMessage(
          'Unable to create endpoint, please check and try again.',
        );
      }
    },
    [],
  );
  /**
   *
   */

  const onvalidateEndpointInstruction = () => {
    const err = ValidateTextProviderDefaultOptions(
      textProviderModel.provider,
      textProviderModel.parameters,
    );
    if (err) {
      setErrorMessage(err);
      return;
    }

    if (promptConfig.variables.length === 0) {
      setErrorMessage(
        'Please provide a valid prompt template, it should at least have one variable.',
      );
      return;
    }

    // Check if the content is not empty
    const hasNonEmptyContent = promptConfig.prompt.some(
      item => item.content.trim() !== '',
    );
    if (!hasNonEmptyContent) {
      setErrorMessage('Please provide content for at least one prompt item.');
      return;
    }

    // If all validations pass, proceed to the next tab
    setErrorMessage('');
    setActiveTab('define-endpoint');
  };

  const createEndpoint = () => {
    if (name.trim() === '') {
      setErrorMessage(
        'Please a valid name for endpoint, that can help you indentify the endpoint',
      );
      return;
    }
    setErrorMessage('');
    showLoader('overlay');

    const endpointProviderModelAttr = new EndpointProviderModelAttribute();
    endpointProviderModelAttr.setModelprovidername(textProviderModel.provider);
    endpointProviderModelAttr.setEndpointmodeloptionsList(
      textProviderModel.parameters,
    );

    endpointProviderModelAttr.setChatcompleteprompt(
      ChatCompletePrompt(promptConfig),
    );

    const endpointattr = new EndpointAttribute();
    endpointattr.setName(name);
    if (description.trim() === '') {
      endpointattr.setDescription(description);
    }

    CreateEndpoint(
      connectionConfig,
      endpointProviderModelAttr,
      endpointattr,
      tags,
      ConnectionConfig.WithDebugger({
        userId: authId,
        authorization: token,
        projectId: projectId,
      }),
      afterCreateEndpoint,
    );
  };

  const [isShow, setIsShow] = useState(false);
  const [isConfigureEndpointPromptOpen, setIsConfigureEndpointPromptOpen] =
    useState(false);

  const handleSelectTemplate = (template: {
    name: string;
    description: string;
    provider: string;
    model: string;
    parameters: { temperature: number; response_format: string };
    instruction: { role: string; content: string }[];
  }) => {
    // Set the name and description
    setName(template.name);
    setDescription(template.description);

    // Set prompt config from template instruction
    const promptMessages = template.instruction.map(inst => ({
      role: inst.role,
      content: inst.content,
    }));

    // Extract variables from the template content (look for {{ variable }} patterns)
    const variableRegex = /\{\{\s*(\w+)(?:\[.*?\])?\s*\}\}/g;
    const variables: { name: string; type: string; defaultvalue: string }[] =
      [];
    const seenVariables = new Set<string>();

    template.instruction.forEach(inst => {
      let match;
      while ((match = variableRegex.exec(inst.content)) !== null) {
        const varName = match[1];
        if (!seenVariables.has(varName)) {
          seenVariables.add(varName);
          variables.push({
            name: varName,
            type: 'string',
            defaultvalue: '',
          });
        }
      }
    });

    setPromptConfig({
      prompt: promptMessages,
      variables:
        variables.length > 0
          ? variables
          : [{ name: 'messages', type: 'string', defaultvalue: '' }],
    });

    // Get default parameters for the template provider
    const newParams = GetDefaultTextProviderConfigIfInvalid(
      template.provider,
      [],
    );

    // Helper function to set or create a parameter
    const setOrCreateParam = (key: string, value: string) => {
      const existingParam = newParams.find(p => p.getKey() === key);
      if (existingParam) {
        existingParam.setValue(value);
      } else {
        const newParam = new Metadata();
        newParam.setKey(key);
        newParam.setValue(value);
        newParams.push(newParam);
      }
    };

    // Update temperature
    setOrCreateParam(
      'model.temperature',
      String(template.parameters.temperature),
    );

    // Update response_format
    if (template.parameters.response_format) {
      setOrCreateParam(
        'model.response_format',
        template.parameters.response_format,
      );
    }

    // Update model name
    setOrCreateParam('model.name', template.model);

    // Update model id
    setOrCreateParam('model.id', `${template.provider}/${template.model}`);

    // Set the provider and parameters
    setTextProviderModel({
      provider: template.provider,
      parameters: newParams,
    });
  };

  return (
    <>
      <Helmet title="Create an Endpoint"></Helmet>
      <ConfigureEndpointPromptDialog
        modalOpen={isConfigureEndpointPromptOpen}
        setModalOpen={setIsConfigureEndpointPromptOpen}
        onSelectTemplate={handleSelectTemplate}
      />
      <ConfirmDialog
        showing={isShow}
        type="warning"
        title={'Are you sure?'}
        content={
          'You want to cancel creating this endpoint? Any unsaved changes will be lost.'
        }
        confirmText={'Confirm'}
        cancelText="Cancel"
        onConfirm={() => {
          navigator(-1);
        }}
        onCancel={() => {
          setIsShow(false);
        }}
        onClose={() => {
          setIsShow(false);
        }}
      />

      <TabForm
        formHeading="Complete all steps to create new endpoint"
        activeTab={activeTab}
        onChangeActiveTab={() => {}}
        errorMessage={errorMessage}
        form={[
          {
            name: 'Choose Model',
            description: 'The model you want to use for your endpoint.',
            code: 'choose-model',
            body: (
              <>
                <YellowNoticeBlock className="flex items-center">
                  <Info className="shrink-0 w-4 h-4" strokeWidth={1.5} />
                  <div className="ms-3 text-sm font-medium">
                    Endpoints allow you to integrate Large Language Models
                    (LLMs) into your application, providing a powerful interface
                    for AI-driven functionalities.
                  </div>
                  <a
                    target="_blank"
                    href="https://doc.rapida.ai/endpoint/overview"
                    className="h-7 flex items-center font-medium hover:underline ml-auto text-yellow-600"
                    rel="noreferrer"
                  >
                    Read documentation
                    <ExternalLink
                      className="shrink-0 w-4 h-4 ml-1.5"
                      strokeWidth={1.5}
                    />
                  </a>
                </YellowNoticeBlock>
                <div className="space-y-6 px-8 pb-8 max-w-4xl ">
                  <BlueNoticeBlock
                    className="p-4 px-6 rounded-md cursor-pointer border"
                    onClick={() => {
                      setIsConfigureEndpointPromptOpen(true);
                    }}
                  >
                    <h2 className="mb-1 text-[15px] font-semibold text-blue-600">
                      Quick Start with Usecase Template
                    </h2>
                    <div className="flex justify-between space-x-2">
                      <p className="text-sm">
                        Choose from pre-configured templates Click to browse
                        available usecases and auto-fill your form
                      </p>
                      <ArrowUpRight
                        className="text-blue-600 shrink-0"
                        strokeWidth={1.5}
                      />
                    </div>
                  </BlueNoticeBlock>
                  <TextProvider
                    onChangeProvider={onChangeTextProvider}
                    onChangeParameter={onChangeTextProviderParameter}
                    parameters={textProviderModel.parameters}
                    provider={textProviderModel.provider}
                  />
                  <ConfigPrompt
                    instanceId={randomString(10)}
                    existingPrompt={promptConfig}
                    onChange={prompt => setPromptConfig(prompt)}
                  />
                </div>
              </>
            ),
            actions: [
              <ICancelButton
                className="px-4 rounded-[2px]"
                onClick={() => setIsShow(true)}
              >
                Cancel
              </ICancelButton>,
              <IBlueBGArrowButton
                type="button"
                isLoading={loading}
                className="px-4 rounded-[2px]"
                onClick={onvalidateEndpointInstruction}
              >
                Configure instruction
              </IBlueBGArrowButton>,
            ],
          },
          {
            code: 'define-endpoint',
            name: 'Define Endpoint Profile',
            description:
              'Provide the name, a brief description, and relevant tags.',
            actions: [
              <ICancelButton
                className="px-4 rounded-[2px]"
                onClick={() => setIsShow(true)}
              >
                Cancel
              </ICancelButton>,
              <IBlueBGArrowButton
                className="px-4 rounded-[2px]"
                type="button"
                isLoading={loading}
                onClick={() => {
                  createEndpoint();
                }}
              >
                Create endpoint
              </IBlueBGArrowButton>,
            ],
            body: (
              <div className="space-y-6 px-8 py-8 max-w-4xl">
                <FieldSet>
                  <FormLabel>Name</FormLabel>
                  <Input
                    name="name"
                    onChange={e => {
                      setName(e.target.value);
                    }}
                    value={name}
                    className="form-input"
                    placeholder="Enter a name"
                  ></Input>
                </FieldSet>
                <FieldSet>
                  <FormLabel>Description</FormLabel>
                  <Textarea
                    row={5}
                    name="description"
                    value={description}
                    placeholder={"What's the purpose of the endpoint?"}
                    onChange={t => setDescription(t.target.value)}
                  />
                </FieldSet>
                <TagInput
                  tags={tags}
                  addTag={onAddTag}
                  removeTag={onRemoveTag}
                  allTags={EndpointTag}
                />
              </div>
            ),
          },
        ]}
      />
    </>
  );
}
