import { Metadata } from '@rapidaai/react';

export type CreateKnowledgeTypeProperty = {
  /**
   *
   * @param string
   * @returns
   */
  description: string;

  /**
   *
   */
  provider: string;

  /**
   *
   */
  providerParamters: Metadata[];

  /**
   * endpoint name
   */
  name: string;

  /**
   * endpoint tags
   */
  tags: string[];
};

export type CreateKnowledgeTypeAction = {
  /**
   *
   * @param s
   * @returns
   */
  onChangeDescription: (s: string) => void;

  /**
   *
   * @param i
   * @param v
   * @returns
   */
  onChangeProvider: (v: string) => void;

  /**
   *
   * @param s
   * @returns
   */
  onChangeProviderParameter: (parameters: Metadata[]) => void;

  /**
   * set name
   */

  onChangeName: (s: string) => void;

  /**
   *
   * @param s
   * @returns
   */
  onRemoveTag: (s: string) => void;

  /**
   *
   * @param s
   * @returns
   */
  onAddTag: (s: string) => void;
};
export type CreateKnowledgeType = {
  /**
   *
   * @param projectId
   * @param token
   * @param userId
   * @param onSuccess
   * @param onError
   * @returns
   */
  onCreateKnowledge: (
    projectId: string,
    token: string,
    userId: string,
    onSuccess: (knowledge: any) => void,
    onError: (err: string) => void,
  ) => void;

  /**
   *
   * @returns
   */
  clear: () => void;
} & CreateKnowledgeTypeProperty &
  CreateKnowledgeTypeAction;
