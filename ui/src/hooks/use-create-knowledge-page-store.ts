import { create } from 'zustand';
import {
  ConnectionConfig,
  CreateKnowledgeRequest,
  Knowledge,
  Metadata,
} from '@rapidaai/react';
import { CreateKnowledgeType } from '@/types/types.create-knowledge';
import { CreateKnowledgeTypeProperty } from '../types/types.create-knowledge';
import { CreateKnowledge } from '@rapidaai/react';
import { CreateKnowledgeResponse } from '@rapidaai/react';
import {
  GetDefaultEmbeddingConfigIfInvalid,
  ValidateEmbeddingDefaultOptions,
} from '@/app/components/providers/embedding';
import { randomMeaningfullName } from '@/utils';
import { connectionConfig } from '@/configs';

/**
 *
 */

const initialState: CreateKnowledgeTypeProperty = {
  /**
   * description of endpoint provider model
   */
  description: '',

  /**
   * name of the endpoint
   */
  name: randomMeaningfullName('knowledge'),

  /**
   *
   */
  provider: 'openai',

  /**
   *
   */
  providerParamters: GetDefaultEmbeddingConfigIfInvalid('openai', []),

  /**
   * list of tags for endpoint
   */
  tags: [],
};

export const useCreateKnowledgePageStore = create<CreateKnowledgeType>(
  (set, get) => ({
    ...initialState,

    onChangeDescription: (st: string) => {
      set({
        description: st,
      });
    },

    /**
     *
     * on change of name of endpoint
     * @param st
     */
    onChangeName: (st: string) => {
      set({
        name: st,
      });
    },

    /**
     *
     */
    onChangeProvider: (v: string) => {
      set({ provider: v });
    },

    /**
     *
     * @param parameters
     */
    onChangeProviderParameter: (parameters: Metadata[]) => {
      set({ providerParamters: parameters });
    },

    /**
     * adding new tags of endpoint
     * @param s
     */
    onAddTag: (s: string) => {
      if (s.trim() !== '') {
        let oldTags = get().tags;
        const index = oldTags.indexOf(s, 0);
        if (index > -1) {
          return;
        }
        const all = [...oldTags, s];
        set({ tags: all });
      }
    },

    /**
     * remove the tags for endpoint
     */
    onRemoveTag: (s: string) => {
      set({
        tags: get().tags.filter(x => x !== s),
      });
    },

    /**
     *
     * @param projectId
     * @param token
     * @param userId
     * @param onSuccess
     * @param onError
     */
    onCreateKnowledge: async (
      projectId: string,
      token: string,
      userId: string,
      onSuccess: (knowledge: Knowledge) => void,
      onError: (err: string) => void,
    ) => {
      // validations
      let _providerModel = get().provider;
      if (!_providerModel) {
        onError('Please select the embedding models.');
        return;
      }

      let err = ValidateEmbeddingDefaultOptions(
        _providerModel,
        get().providerParamters,
      );
      if (err) {
        onError(err);
        return;
      }
      let _name = get().name;
      if (!_name) {
        onError('Please enter the name of knowledge base.');
        return;
      }
      let _description = get().description;
      if (!_description) {
        onError('Please enter the description of knowledge base.');
        return;
      }

      let _tags = get().tags;

      const req = new CreateKnowledgeRequest();
      req.setEmbeddingmodelprovidername(_providerModel);
      req.setKnowledgeembeddingmodeloptionsList(get().providerParamters);
      req.setName(_name);
      req.setDescription(_description);
      req.setTagsList(_tags);
      CreateKnowledge(
        connectionConfig,
        req,
        ConnectionConfig.WithDebugger({
          authorization: token,
          userId: userId,
          projectId: projectId,
        }),
      )
        .then((car: CreateKnowledgeResponse | null) => {
          if (car?.getSuccess()) {
            let assistant = car.getData();
            if (assistant) onSuccess(assistant);
          } else {
            const errorMessage =
              'Unable to create knowledge. please try again later.';
            const error = car?.getError();
            if (error) {
              onError(error.getHumanmessage());
              return;
            }
            onError(errorMessage);
            return;
          }
        })
        .catch(error => {
          onError('Unable to create knowledge. please try again later.');
        });
    },

    /**
     * clear everything from the context
     * @returns
     */
    clear: () => {
      set({ ...initialState }, false);
    },
  }),
);
