import { SearchableDeployment } from '@rapidaai/react';
import { InputCheckbox } from '@/app/components/Form/Checkbox';
import { SearchIconInput } from '@/app/components/Form/Input/IconInput';
import { Label } from '@/app/components/Form/Label';
import { DeploymentIcon } from '@/app/components/Icon/Deployment';
import { LanguageIcon } from '@/app/components/Icon/Language';
import { TablePagination } from '@/app/components/base/tables/table-pagination';
import { BluredWrapper } from '@/app/components/Wrapper/BluredWrapper';
import { useDiscoverDeploymentPageStore } from '@/hooks';
import { useRapidaStore } from '@/hooks';
import { useCredential } from '@/hooks/use-credential';
import React, { useCallback, useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import { HubEndpointCard } from '@/app/components/base/cards/endpoint-card';
import { HubAssistantCard } from '@/app/components/base/cards/assistant-card';
import { Helmet } from '@/app/components/Helmet';
import { ActionableEmptyMessage } from '@/app/components/container/message/actionable-empty-message';
import { Spinner } from '@/app/components/Loader/Spinner';

export const DiscoverPage = () => {
  const {
    deployments,
    page,
    setPage,
    totalCount,
    pageSize,
    setPageSize,
    getAllDeployments,
    addCriteria,
    criteria,
    allLanguage,
    allUsecase,
  } = useDiscoverDeploymentPageStore();

  /**
   * all the parameters
   */
  const [userId, token, projectId] = useCredential();
  const { loading, showLoader, hideLoader } = useRapidaStore();
  /**
   * getting all the models
   */

  const onError = useCallback((err: string) => {
    hideLoader();
    toast.error(err);
  }, []);

  const onSuccess = useCallback((data: SearchableDeployment[]) => {
    hideLoader();
  }, []);
  /**
   * call the api
   */
  const getDeployments = useCallback((projectId, token, userId) => {
    showLoader();
    getAllDeployments(token, userId, onError, onSuccess);
  }, []);

  useEffect(() => {
    getDeployments(projectId, token, userId);
  }, [page, pageSize]);

  /**
   * working with models
   */
  const [selectedModels, setSelectedModels] = useState<string[]>([]);
  useEffect(() => {
    addCriteria(
      'EndpointProviderModel.provider_model_id',
      selectedModels.toString(),
      'oneOf',
    );
    getDeployments(projectId, token, userId);
  }, [selectedModels]);

  /**
   * working with usecase
   */
  const [selectedUsecases, setSelectedUsecase] = useState<string[]>([]);
  useEffect(() => {
    addCriteria('tag.keyword', selectedUsecases.toString(), 'oneOf');
    getDeployments(projectId, token, userId);
  }, [selectedUsecases]);

  /**
   * working with language
   */
  const [selectedLanguages, setSelectedLanguages] = useState<string[]>([]);
  useEffect(() => {
    addCriteria('language.keyword', selectedLanguages.toString(), 'oneOf');
    getDeployments(projectId, token, userId);
  }, [selectedLanguages]);

  return (
    <>
      <Helmet title="Rapida Hub" />
      <div className="grid grid-cols-2 md:grid-cols-5 h-full overflow-hidden">
        <div className="flex-col overflow-auto no-scrollbar col-span-1 items-start gap-6 p-4 border-r dark:border-gray-700/50 bg-white dark:bg-gray-900 hidden md:flex transition-none delay-700">
          <div className="flex flex-col items-start justify-start space-y-8 shrink-0  h-full">
            {/* Usecase */}
            <div className="space-y-4">
              <div className="flex flex-col gap-1 self-stretch">
                <div className="inline-flex cursor-pointer items-center justify-between gap-3 self-stretch">
                  <div className="flex items-center justify-start gap-2">
                    <DeploymentIcon className="w-4 h-4" />
                    <Label>Usecases</Label>
                  </div>
                </div>
              </div>
              <div className="flex flex-col items-start justify-start gap-3 self-stretch">
                {allUsecase.map((ix, id) => {
                  return (
                    <div
                      className="flex flex-1 flex-row items-center justify-stretch gap-2 self-stretch text-sm"
                      key={id}
                    >
                      <InputCheckbox
                        className="w-3.5 h-3.5"
                        id={ix}
                        onChange={e => {
                          const isChecked = e.target.checked;
                          const updatedSelection = isChecked
                            ? [...selectedUsecases, ix]
                            : selectedUsecases.filter(x => x !== ix);
                          setSelectedUsecase(updatedSelection);
                        }}
                        checked={selectedUsecases.includes(ix)}
                      />
                      <Label for={ix} text={ix} />
                    </div>
                  );
                })}
              </div>
            </div>

            {/* language */}
            <div className="space-y-4">
              <div className="flex flex-col gap-1 self-stretch">
                <div className="inline-flex cursor-pointer items-center justify-between gap-3 self-stretch">
                  <div className="flex items-center justify-start gap-2">
                    <LanguageIcon className="w-4 h-4" />
                    <Label>Language</Label>
                  </div>
                </div>
              </div>
              <div className="flex flex-col items-start justify-start gap-3 self-stretch">
                {allLanguage.map((ix, id) => {
                  return (
                    <div
                      className="flex flex-1 flex-row items-center justify-stretch gap-2 self-stretch text-sm"
                      key={id}
                    >
                      <InputCheckbox
                        id={ix}
                        className="w-3.5 h-3.5"
                        onChange={e => {
                          const isChecked = e.target.checked;
                          const updatedSelection = isChecked
                            ? [...selectedLanguages, ix]
                            : selectedLanguages.filter(x => x !== ix);
                          setSelectedLanguages(updatedSelection);
                        }}
                        checked={selectedLanguages.includes(ix)}
                      />
                      <Label for={ix} text={ix} />
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        </div>
        <div className="flex flex-col md:col-span-4 col-span-2 gap-8 w-full overflow-auto">
          <BluredWrapper className="border-t-0">
            <SearchIconInput
              className="pl-8"
              iconClassName="w-4 h-4"
              placeholder="Search across all the usecases"
            />
            <TablePagination
              currentPage={page}
              onChangeCurrentPage={setPage}
              totalItem={totalCount}
              pageSize={pageSize}
              onChangePageSize={setPageSize}
            />
          </BluredWrapper>

          {deployments && deployments.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-3 px-8">
              {deployments.map((deployment, id) => {
                return deployment.getType() === 'endpoint' ? (
                  <HubEndpointCard deployment={deployment} key={id} />
                ) : (
                  <HubAssistantCard deployment={deployment} key={id} />
                );
              })}
            </div>
          ) : criteria.length > 0 ? (
            <ActionableEmptyMessage
              title="No public template "
              subtitle="There are no template matching with your criteria."
            />
          ) : !loading ? (
            <ActionableEmptyMessage
              title="No public template "
              subtitle="There are no template matching with your criteria."
            />
          ) : (
            <div className="h-full flex justify-center items-center">
              <Spinner size="md" />
            </div>
          )}
        </div>
      </div>
    </>
  );
};
