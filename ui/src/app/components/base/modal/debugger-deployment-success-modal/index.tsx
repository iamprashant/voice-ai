import { ICancelButton, ILinkButton } from '@/app/components/form/button';
import { GenericModal } from '@/app/components/base/modal';
import { ModalFooter } from '@/app/components/base/modal/modal-footer';
import type { FC } from 'react';
import { ModalProps } from '../index';
import { ModalFitHeightBlock } from '@/app/components/blocks/modal-fit-height-block';
import { ModalHeader } from '@/app/components/base/modal/modal-header';
import { ModalBody } from '@/app/components/base/modal/modal-body';
import { CheckCircle, ExternalLink } from 'lucide-react';
import { ModalTitleBlock } from '@/app/components/blocks/modal-title-block';

interface DebuggerDeploymentSuccessDialogProps extends ModalProps {
  assistantId: string;
}

export const DebuggerDeploymentSuccessDialog: FC<
  DebuggerDeploymentSuccessDialogProps
> = ({ modalOpen, setModalOpen, assistantId }) => {
  return (
    <GenericModal modalOpen={modalOpen} setModalOpen={setModalOpen}>
      <ModalFitHeightBlock>
        <ModalHeader
          onClose={() => {
            setModalOpen(false);
          }}
        ></ModalHeader>
        <ModalBody>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center bg-green-500/10">
              <CheckCircle className="h-6 w-6 text-green-500" />
            </div>
            <div>
              <ModalTitleBlock>Success!!</ModalTitleBlock>
              <p className="text-sm text-muted">
                New version of debugger deployment has been updated
                successfully.
              </p>
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <ICancelButton
            className="px-4 rounded-[2px]"
            onClick={() => {
              setModalOpen(false);
            }}
          >
            Close
          </ICancelButton>
          <ILinkButton
            href={`/preview/chat/${assistantId}`}
            className="px-4 rounded-[2px] space-x-1"
            type="button"
          >
            <span>Preview assistant</span>
            <ExternalLink className="w-4 h-4" strokeWidth={1.5} />
          </ILinkButton>
        </ModalFooter>
      </ModalFitHeightBlock>
    </GenericModal>
  );
};
