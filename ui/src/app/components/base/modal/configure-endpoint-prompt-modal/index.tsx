import React, { FC, useState } from 'react';
import { IBlueBGButton, ICancelButton } from '@/app/components/form/button';
import { GenericModal, ModalProps } from '@/app/components/base/modal';
import { ModalFitHeightBlock } from '@/app/components/blocks/modal-fit-height-block';
import { ModalHeader } from '@/app/components/base/modal/modal-header';
import { ModalTitleBlock } from '@/app/components/blocks/modal-title-block';
import { ModalBody } from '@/app/components/base/modal/modal-body';
import { ModalFooter } from '@/app/components/base/modal/modal-footer';
import { cn } from '@/utils';
import endpointTemplates from '@/prompts/endpoints/index.json';
import { ChartColumn, Check } from 'lucide-react';

interface EndpointTemplate {
  name: string;
  description: string;
  provider: string;
  model: string;
  parameters: {
    temperature: number;
    response_format: string;
  };
  instruction: {
    role: string;
    content: string;
  }[];
}

interface ConfigureEndpointPromptDialogProps extends ModalProps {
  onSelectTemplate?: (template: EndpointTemplate) => void;
}

export const ConfigureEndpointPromptDialog: FC<
  ConfigureEndpointPromptDialogProps
> = props => {
  const [selectedTemplate, setSelectedTemplate] =
    useState<EndpointTemplate | null>(null);

  const handleContinue = () => {
    if (selectedTemplate && props.onSelectTemplate) {
      props.onSelectTemplate(selectedTemplate);
    }
    props.setModalOpen(false);
  };

  return (
    <GenericModal
      className="flex"
      modalOpen={props.modalOpen}
      setModalOpen={props.setModalOpen}
    >
      <ModalFitHeightBlock className="w-[1000px]">
        <ModalHeader
          onClose={() => {
            props.setModalOpen(false);
          }}
          title={'Select a usecase template'}
        >
          <ModalTitleBlock>Select a usecase template</ModalTitleBlock>
        </ModalHeader>
        <ModalBody className="overflow-auto max-h-[80dvh] px-4 py-4">
          <div className="grid grid-cols-2 gap-4">
            {(endpointTemplates as EndpointTemplate[]).map(
              (template, index) => (
                <div
                  key={index}
                  onClick={() => setSelectedTemplate(template)}
                  className={cn(
                    'relative p-4 border rounded-lg cursor-pointer transition-all hover:shadow-md bg-white dark:bg-gray-950',
                    selectedTemplate?.name === template.name &&
                      'border-blue-500 bg-blue-50 ring-1 ring-blue-500',
                  )}
                >
                  {selectedTemplate?.name === template.name && (
                    <div className="absolute top-3 right-3 w-5 h-5 bg-blue-500 rounded-full flex items-center justify-center">
                      <Check className="w-3 h-3 text-white" strokeWidth={3} />
                    </div>
                  )}
                  <div className="p-2 bg-gray-200 dark:bg-gray-800 w-fit rounded-md mb-3">
                    <ChartColumn className="w-6 h-6" strokeWidth={1.5} />
                  </div>
                  <h3 className="font-semibold mb-2 pr-6 text-base">
                    {template.name}
                  </h3>
                  <p className="line-clamp-2 text-sm text-muted mb-3">
                    {template.description}
                  </p>
                  <div className="flex flex-wrap gap-2 mt-auto">
                    <span className="inline-flex items-center px-2 py-1 text-xs font-medium bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300 rounded capitalize">
                      {template.provider}
                    </span>
                    <span className="inline-flex items-center px-2 py-1 text-xs font-medium bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300 rounded">
                      {template.model}
                    </span>
                    <span className="inline-flex items-center px-2 py-1 text-xs font-medium bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded">
                      Temp: {template.parameters.temperature}
                    </span>
                    {template.parameters.response_format && (
                      <span className="inline-flex items-center px-2 py-1 text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded">
                        JSON Schema
                      </span>
                    )}
                  </div>
                </div>
              ),
            )}
          </div>
        </ModalBody>
        <ModalFooter>
          <ICancelButton
            className="px-4 rounded-[2px]"
            onClick={() => {
              props.setModalOpen(false);
            }}
          >
            Cancel
          </ICancelButton>
          <IBlueBGButton
            className="px-4 rounded-[2px]"
            type="button"
            disabled={!selectedTemplate}
            onClick={handleContinue}
          >
            Continue
          </IBlueBGButton>
        </ModalFooter>
      </ModalFitHeightBlock>
    </GenericModal>
  );
};
