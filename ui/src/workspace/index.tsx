import { CONFIG, WorkspaceConfig } from '@/configs';

import React, { createContext, useContext } from 'react';

/**
 *
 */
const WorkspaceContext = createContext<WorkspaceConfig>(CONFIG.workspace);

/**
 *
 * @returns
 */
export const useWorkspace = () => useContext(WorkspaceContext);

/**
 *
 * @param param0
 * @returns
 */
export const WorkspaceProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <WorkspaceContext.Provider value={CONFIG.workspace}>
      {children}
    </WorkspaceContext.Provider>
  );
};
