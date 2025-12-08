import { GenericModal, ModalProps } from '@/app/components/base/modal';
import { ModalBody } from '@/app/components/base/modal/modal-body';
import { ModalFooter } from '@/app/components/base/modal/modal-footer';
import { ModalHeader } from '@/app/components/base/modal/modal-header';
import { ModalFitHeightBlock } from '@/app/components/blocks/modal-fit-height-block';
import { ModalTitleBlock } from '@/app/components/blocks/modal-title-block';
import { IBlueBGButton, ICancelButton } from '@/app/components/form/button';
import { CheckCircle, ChevronRight, ExternalLink, Globe } from 'lucide-react';
import { FC } from 'react';
import { Bug, Code, PhoneCall } from 'lucide-react';
import { useGlobalNavigation } from '@/hooks/use-global-navigator';
import { Assistant } from '@rapidaai/react';

/**
 *
 */
interface ConfigureAssistantNextDialogProps extends ModalProps {
  assistant: Assistant;
}

/**
 *
 * @param param0
 * @returns
 */
export const ConfigureAssistantNextDialog: FC<
  ConfigureAssistantNextDialogProps
> = ({ assistant, modalOpen, setModalOpen }) => {
  /**
   * navigation
   */
  const {
    goBack,
    goToAssistant,
    goToAssistantPreview,
    goToConfigureDebugger,
    goToConfigureWeb,
    goToConfigureCall,
    goToConfigureApi,
    goToCreateAssistantAnalysis,
    goToCreateAssistantWebhook,
  } = useGlobalNavigation();

  /**
   *
   */
  const deploymentOptions = [
    {
      icon: PhoneCall,
      title: 'Phone call',
      description: 'Enable voice conversations over phone call',
      action: 'Enable phone call',
      onclick: () => {
        goToConfigureCall(assistant.getId());
      },
    },
    {
      icon: Code,
      title: 'API',
      description: 'Integrate into your application using sdks',
      action: 'Enable Api',
      onclick: () => {
        goToConfigureApi(assistant.getId());
      },
    },
    {
      icon: Globe,
      title: 'Web Widget',
      description:
        'Embed on your website to handle text and voice customer query.',
      action: 'Deploy to Web Widget',
      onclick: () => {
        goToConfigureWeb(assistant.getId());
      },
    },
    {
      icon: Bug,
      title: 'Debugger / Testing',
      description: 'Deploy the agent for testing and debugging.',
      action: 'Deploy to Debugger / Testing',
      onclick: () => {
        goToConfigureDebugger(assistant.getId());
      },
    },
  ];

  return (
    <GenericModal
      className="flex"
      modalOpen={modalOpen}
      setModalOpen={setModalOpen}
    >
      <ModalFitHeightBlock className="w-[1000px]">
        <ModalHeader
          onClose={() => {
            setModalOpen(false);
          }}
        >
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center bg-green-500/10">
              <CheckCircle className="h-6 w-6 text-green-500" />
            </div>
            <div>
              <ModalTitleBlock>Assistant Created Successfully</ModalTitleBlock>
              <p className="text-sm text-muted">
                Configure deployments and integrations below
              </p>
            </div>
          </div>
        </ModalHeader>
        <ModalBody className="overflow-auto max-h-[80dvh]">
          <div className="">
            <div className="px-6 py-3 bg-amber-500/10 border-b-2 border-amber-500">
              <h3 className="text-sm font-medium text-foreground">
                Configure deployments to integrate channels into assistants
              </h3>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 border border-t-0 divide-x">
              {deploymentOptions.map(option => (
                <div key={option.title} className="p-4 flex flex-col">
                  <option.icon
                    className="h-5 w-5 text-muted mb-3"
                    strokeWidth={1.5}
                  />
                  <h4 className="font-medium text-sm text-foreground mb-1">
                    {option.title}
                  </h4>
                  <p className="text-xs text-muted-foreground flex-1 mb-3">
                    {option.description}
                  </p>
                  <button
                    onClick={option.onclick}
                    className="flex items-center gap-1 text-xs text-primary hover:underline font-medium cursor-pointer"
                  >
                    {option.action}
                    <ChevronRight className="h-3 w-3" />
                  </button>
                </div>
              ))}
            </div>
          </div>

          {/*  */}
          <div className="">
            <div className="px-6 py-3 bg-amber-500/10 border-b-2 border-amber-500">
              <h3 className="text-sm font-medium text-foreground">
                Enable Post-Conversation analysis
              </h3>
            </div>
            <div className="px-6 py-4 border border-t-0">
              <p className="text-sm text-muted-foreground">
                Gain insights from every interaction eg: Automatic conversation
                transcripts Quality, sentiment, and SOP adherence analysis
                Custom reporting and dashboards
              </p>
              <button
                onClick={() => {
                  goToCreateAssistantAnalysis(assistant.getId());
                }}
                className="flex items-center gap-1 text-xs text-primary hover:underline font-medium cursor-pointer mt-3"
              >
                <span>Configure analysis</span>
                <ChevronRight className="h-3 w-3" />
              </button>
            </div>
          </div>

          {/* Webhook & Integration Section */}
          <div className="">
            <div className="px-6 py-3 bg-amber-500/10 border-b-2 border-amber-500">
              <h3 className="text-sm font-medium text-foreground">
                Configure deployments to integrate channels into assistants
              </h3>
            </div>
            <div className="px-6 py-4 border border-t-0">
              <p className="text-sm text-muted-foreground">
                Keep your workflows connected by triggering events when key
                actions happen: eg: Conversation started / ended Escalation to a
                human agent Custom events for analytics or CRM sync
              </p>
              <button
                onClick={() => {
                  goToCreateAssistantWebhook(assistant.getId());
                }}
                className="flex items-center gap-1 text-xs text-primary hover:underline font-medium cursor-pointer mt-3"
              >
                <span>Configure webhook</span>
                <ChevronRight className="h-3 w-3" />
              </button>
            </div>
          </div>
        </ModalBody>
        <ModalFooter errorMessage={''}>
          <ICancelButton
            className="px-4 rounded-[2px]"
            onClick={() => {
              setModalOpen(false);
            }}
          >
            Cancel
          </ICancelButton>
          <IBlueBGButton
            className="px-4 rounded-[2px]"
            type="button"
            onClick={() => {
              goToAssistantPreview(assistant.getId());
            }}
          >
            <span>Preview assistant</span>
            <ExternalLink className="w-4 h-4 ml-1" strokeWidth={1.5} />
          </IBlueBGButton>
        </ModalFooter>
      </ModalFitHeightBlock>
    </GenericModal>
  );
};
