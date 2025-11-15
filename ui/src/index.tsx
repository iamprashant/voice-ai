import 'react-app-polyfill/ie11';
import 'react-app-polyfill/stable';
import * as React from 'react';
import ReactDOM from 'react-dom/client';
import { App } from '@/app';
import { HelmetProvider } from 'react-helmet-async';
import { DarkModeProvider } from '@/context/dark-mode-context';
import { WorkspaceProvider } from '@/workspace';
import { initializeAnalytics } from '@/react-web-analytics';
const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement,
);
initializeAnalytics();
root.render(
  <HelmetProvider>
    <React.StrictMode>
      <DarkModeProvider>
        <WorkspaceProvider>
          <App />
        </WorkspaceProvider>
      </DarkModeProvider>
    </React.StrictMode>
  </HelmetProvider>,
);
