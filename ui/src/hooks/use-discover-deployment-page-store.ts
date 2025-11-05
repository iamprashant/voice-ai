import {
  GetAllDeployment,
  GetAllDeploymentResponse,
  SearchableDeployment,
} from '@rapidaai/react';
import { ServiceError } from '@rapidaai/react';
import { EndpointTag } from '@/app/components/Form/tag-input/endpoint-tags';
import { DiscoverDeploymentType } from '@/types';
import { initialPaginated } from '@/types/types.paginated';
import { create } from 'zustand';
import { connectionConfig } from '@/configs';

/**
 *
 */
export const useDiscoverDeploymentPageStore = create<DiscoverDeploymentType>(
  (set, get) => ({
    ...initialPaginated,

    /**
     *
     * @param number
     * @returns
     */
    setPageSize: (pageSize: number) => {
      // when someone change pagesize change the page to zero
      set({
        page: 1,
        pageSize: pageSize,
      });
    },

    /**
     *
     * @param number
     * @returns
     */
    setPage: (pg: number) => {
      set({
        page: pg,
      });
    },

    /**
     *
     * @param number
     * @returns
     */
    setTotalCount: (tc: number) => {
      set({
        totalCount: tc,
      });
    },
    /**
     *
     * @param k
     * @param v
     */
    addCriteria: (k: string, v: string, condition: string) => {
      let current = get().criteria.filter(x => x.key !== k);
      if (v) current.push({ key: k, value: v, logic: condition });
      set({
        criteria: current,
      });
    },

    /**
     *
     */
    deployments: [],

    /**
     *
     * @param ep
     * @returns
     */
    setAllDeployments: (ep: SearchableDeployment[]) => {
      set({
        deployments: ep,
      });
    },
    /**
     *
     * @param projectId
     * @param token
     * @param userId
     */
    getAllDeployments: (
      token: string,
      userId: string,
      onError: (err: string) => void,
      onSuccess: (e: SearchableDeployment[]) => void,
    ) => {
      const afterGetAllDeployment = (
        err: ServiceError | null,
        gur: GetAllDeploymentResponse | null,
      ) => {
        if (gur?.getSuccess()) {
          get().setAllDeployments(gur.getDataList());
          let paginated = gur.getPaginated();
          if (paginated) {
            get().setTotalCount(paginated.getTotalitem());
          }
          onSuccess(gur.getDataList());
        } else {
          let errorMessage = gur?.getError();
          if (errorMessage) {
            onError(errorMessage.getHumanmessage());
            return;
          }
          onError('Unable to get deployments, please try again later.');
        }
      };

      GetAllDeployment(
        connectionConfig,
        get().page,
        get().pageSize,
        get().criteria,
        afterGetAllDeployment,
        {
          authorization: token,
          'x-auth-id': userId,
        },
      );
    },

    /**
     * all the models
     */
    allModel: [],
    allLanguage: ['english', 'russian', 'spanish'],
    allUsecase: EndpointTag,

    /**
     *
     */
    clearCriteria: () => {
      set({
        criteria: [],
      });
    },

    clear: () => set({}, true),
  }),
);
