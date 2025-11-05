import { SearchableDeployment } from '@rapidaai/react';
import { PaginatedType } from './types.paginated';

/**
 * endpoint context
 */
export type DiscoverDeploymentType = {
  /**
   *
   */
  deployments: SearchableDeployment[];
  allLanguage: string[];
  allUsecase: string[];

  /**
   *
   * @param ep
   * @returns
   */
  setAllDeployments: (deployments: SearchableDeployment[]) => void;

  /**
   *
   * @param token
   * @param userId
   * @param onError
   * @param onSuccess
   * @returns
   */
  getAllDeployments: (
    token: string,
    userId: string,
    onError: (err: string) => void,
    onSuccess: (e: SearchableDeployment[]) => void,
  ) => void;

  /**
   * clear everything
   * @returns
   */
  clear: () => void;
} & PaginatedType;
