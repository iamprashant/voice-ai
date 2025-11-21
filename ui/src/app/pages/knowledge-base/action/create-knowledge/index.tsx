import { FieldSet } from '@/app/components/form/fieldset';
import { Helmet } from '@/app/components/helmet';
import { Input } from '@/app/components/form/input';
import { TagInput } from '@/app/components/form/tag-input';
import { useCreateKnowledgePageStore } from '@/hooks/use-create-knowledge-page-store';
import { useCallback, useEffect, useState } from 'react';
import { useCredential, useRapidaStore } from '@/hooks';
import {
  IBlueBGArrowButton,
  ICancelButton,
} from '@/app/components/form/button';
import { KnowledgeDocument } from '@rapidaai/react';
import { TabForm } from '@/app/components/form/tab-form';
import { useCreateKnowledgeDocumentPageStore } from '@/hooks/use-create-knowledge-document-page-store';
import { Knowledge } from '@rapidaai/react';
import { KnowledgeTags } from '@/app/components/form/tag-input/knowledge-tags';
import toast from 'react-hot-toast/headless';
import { create_knowledge_success_message } from '@/utils/messages';
import { useGlobalNavigation } from '@/hooks/use-global-navigator';
import { ManualFile } from '@/app/pages/knowledge-base/action/components/datasource-uploader/manual-file';
import { FormLabel } from '@/app/components/form-label/index';
import ConfirmDialog from '@/app/components/base/modal/confirm-ui';
import { useNavigate } from 'react-router-dom';
import { Textarea } from '@/app/components/form/textarea';
import { EmbeddingProvider } from '@/app/components/providers/embedding';
import { YellowNoticeBlock } from '@/app/components/container/message/notice-block';
import { ExternalLink, Info } from 'lucide-react';

/**
 *
 * @returns
 */
export function CreateKnowledgePage() {
  /**
   * tab management for multistep form
   */
  const [activeTab, setActiveTab] = useState('create-knowledge');
  const [errorMessage, setErrorMessage] = useState('');
  const { goToKnowledge } = useGlobalNavigation();
  const {
    name,
    clear,
    onChangeName,
    description,
    onChangeDescription,
    tags,
    onAddTag,
    onRemoveTag,
    onCreateKnowledge,
    provider,
    onChangeProvider,
    providerParamters,
    onChangeProviderParameter,
  } = useCreateKnowledgePageStore();

  /**
   *
   */
  const knowledgeDocumentAction = useCreateKnowledgeDocumentPageStore();
  /**
   * all the credentials you will need do things
   */
  const [userId, token, projectId] = useCredential();

  /**
   * show and hide loaders
   */
  const { loading, showLoader, hideLoader } = useRapidaStore();

  /**
   *
   */
  const [knowledge, setKnowledge] = useState<Knowledge | null>(null);

  //
  //
  useEffect(() => {
    clear();
  }, []);

  /**
   *
   */
  const onsuccessofcreateknowledge = (kn: Knowledge) => {
    toast.success(create_knowledge_success_message(kn?.getName()));
    setKnowledge(kn);
    hideLoader();
    setActiveTab('upload-document');
  };

  /**
   *
   */
  const onerror = useCallback((err: string) => {
    hideLoader();
    setErrorMessage(err);
  }, []);
  /**
   *
   */
  const createKnowledge = () => {
    showLoader();
    setErrorMessage('');
    onCreateKnowledge(
      projectId,
      token,
      userId,
      onsuccessofcreateknowledge,
      onerror,
    );
  };

  //   on success
  const onsuccessknowledgedocument = useCallback(
    (d: KnowledgeDocument[]) => {
      hideLoader();
      if (knowledge?.getId()) goToKnowledge(knowledge?.getId());
    },
    [JSON.stringify(knowledge)],
  );

  //   on create document
  const createKnowledgeDocument = () => {
    if (knowledge) {
      showLoader();
      knowledgeDocumentAction.onCreateKnowledgeDocument(
        knowledge.getId(),
        projectId,
        token,
        userId,
        onsuccessknowledgedocument,
        onerror,
      );
    }
  };
  let navigator = useNavigate();
  const [isShow, setIsShow] = useState(false);
  return (
    <>
      <Helmet title="Create a knowledge"></Helmet>

      <ConfirmDialog
        showing={isShow}
        type="warning"
        title={'Are you sure?'}
        content={
          'You want to cancel creating the knowledge? Any unsaved changes will be lost.'
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
        activeTab={activeTab}
        formHeading="Please complete steps to create knowledge"
        onChangeActiveTab={(tabName: string) => {}}
        errorMessage={errorMessage}
        form={[
          {
            name: 'Create Knowledge',
            description:
              'By creating a knowledge you combine multiple documents from different source.',
            code: 'create-knowledge',
            body: (
              <>
                <YellowNoticeBlock className="flex items-center">
                  <Info className="shrink-0 w-4 h-4" strokeWidth={1.5} />
                  <div className="ms-3 text-sm font-medium">
                    A collection in the knowledge base is a curated set of
                    documents grouped by a similar domain, topic, or purpose.
                  </div>
                  <a
                    target="_blank"
                    href="https://doc.rapida.ai/knowledge/overview"
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
                <div className="px-8 pb-8 space-y-6 max-w-3xl">
                  <EmbeddingProvider
                    onChangeParameter={onChangeProviderParameter}
                    onChangeProvider={onChangeProvider}
                    parameters={providerParamters}
                    provider={provider}
                  />
                  <FieldSet>
                    <FormLabel>Name</FormLabel>
                    <Input
                      name="knowledge_name"
                      onChange={e => {
                        onChangeName(e.target.value);
                      }}
                      className="form-input"
                      value={name}
                      placeholder="Name of your knowledge base eg: our website database"
                    ></Input>
                  </FieldSet>
                  <FieldSet>
                    <FormLabel>Description</FormLabel>
                    <Textarea
                      value={description}
                      placeholder={
                        "What's the purpose of the knowledge or what it contains?"
                      }
                      onChange={t => onChangeDescription(t.target.value)}
                    />
                  </FieldSet>
                  <TagInput
                    tags={tags}
                    addTag={onAddTag}
                    allTags={KnowledgeTags}
                    removeTag={onRemoveTag}
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
                className="px-4 rounded-[2px]"
                onClick={createKnowledge}
                isLoading={loading}
              >
                Create knowledge
              </IBlueBGArrowButton>,
            ],
          },
          {
            name: 'Add new document',
            description: 'Adding more documents to your knowledge base',
            code: 'upload-document',
            body: <ManualFile />,
            actions: [
              <ICancelButton
                className="px-4 rounded-[2px]"
                onClick={() => {
                  if (knowledge?.getId()) goToKnowledge(knowledge?.getId());
                }}
              >
                Skip
              </ICancelButton>,
              <IBlueBGArrowButton
                type="button"
                className="px-4 rounded-[2px]"
                onClick={createKnowledgeDocument}
                isLoading={loading}
              >
                Upload new document
              </IBlueBGArrowButton>,
            ],
          },
        ]}
      />
    </>
  );
}
