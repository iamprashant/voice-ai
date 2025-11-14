import { createContext, useContext } from 'use-context-selector';
import { useCallback, useEffect, useState } from 'react';
import { ServiceError } from '@rapidaai/react';
import {
  GetAllOrganizationCredentialResponse,
  VaultCredential,
} from '@rapidaai/react';
import { useCredential } from '@/hooks';
import { GetAllOrganizationCredential } from '@rapidaai/react';

import {
  LOCAL_STORAGE_PROVIDER_CREDENTIALS,
  LOCAL_STORAGE_TOOL_CREDENTIALS,
  serializeProto,
  useLocalStorageSync,
} from '@/hooks/use-storage-sync';
import { connectionConfig } from '@/configs';

const ProviderContext = createContext<{
  providerCredentials: VaultCredential[];
  toolProviderCredentials: VaultCredential[];
  reloadProviderCredentials: () => void;
  reloadToolCredentials: () => void;
}>({
  providerCredentials: [],
  toolProviderCredentials: [],
  reloadProviderCredentials: () => {
    throw new Error('Function not implemented.');
  },
  reloadToolCredentials: () => {
    throw new Error('Function not implemented.');
  },
});

export const useProviderContext = () => useContext(ProviderContext);

type ProviderContextProviderProps = {
  children: React.ReactNode;
};

export const ProviderContextProvider = ({
  children,
}: ProviderContextProviderProps) => {
  const [providerCredentials, setProviderCredentials] = useState<
    VaultCredential[]
  >([]);
  const [toolProviderCredentials, setToolProviderCredentials] = useState<
    VaultCredential[]
  >([]);
  const [userId, token] = useCredential();

  useLocalStorageSync(
    LOCAL_STORAGE_PROVIDER_CREDENTIALS,
    setProviderCredentials,
    VaultCredential,
  );
  useLocalStorageSync(
    LOCAL_STORAGE_TOOL_CREDENTIALS,
    setToolProviderCredentials,
    VaultCredential,
  );

  /**
   *
   */
  useEffect(() => {
    if (token && userId) {
      getAllOrganizationCredential();
    }
  }, [token, userId]);

  /**
   * after getting all the credentials to store in the local storage
   */
  const afterGettingAllCredential = useCallback(
    (
      err: ServiceError | null,
      gapcr: GetAllOrganizationCredentialResponse | null,
    ) => {
      if (gapcr?.getSuccess()) {
        const credentials = gapcr.getDataList();
        setProviderCredentials(
          credentials.filter(
            (x: VaultCredential) => x.getVaulttype() === 'provider-vault',
          ),
        );
        setToolProviderCredentials(
          credentials.filter(
            (x: VaultCredential) => x.getVaulttype() === 'tool-vault',
          ),
        );
        localStorage.setItem(
          LOCAL_STORAGE_PROVIDER_CREDENTIALS,
          JSON.stringify(
            credentials.map((cred: any) => Array.from(serializeProto(cred))),
          ),
        );
        localStorage.setItem(
          LOCAL_STORAGE_TOOL_CREDENTIALS,
          JSON.stringify(
            credentials.map((cred: any) => Array.from(serializeProto(cred))),
          ),
        );
      }
    },
    [],
  );

  /**
   * gettung all the organization
   */
  const getAllOrganizationCredential = () => {
    GetAllOrganizationCredential(
      connectionConfig,
      1,
      50,
      [],
      afterGettingAllCredential,
      {
        authorization: token,
        'x-auth-id': userId,
      },
    );
  };

  /**
   * reload provider credentials
   */
  const reloadProviderCredentials = () => {
    getAllOrganizationCredential();
  };

  /**
   * reloading the tool credentials
   */
  const reloadToolCredentials = () => {
    getAllOrganizationCredential();
  };

  return (
    <ProviderContext.Provider
      value={{
        providerCredentials,
        toolProviderCredentials,
        reloadProviderCredentials,
        reloadToolCredentials,
      }}
    >
      {children}
    </ProviderContext.Provider>
  );
};

export default ProviderContext;
