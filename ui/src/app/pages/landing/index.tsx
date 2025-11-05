/**
 * Asynchronously loads the component for HomePage
 */

import { lazyLoad } from '@/utils/loadable';
import { PageLoader } from '@/app/components/Loader/page-loader';

export const HomePage = lazyLoad(
  () => import('./v5'),
  module => module.V5,
  {
    fallback: <PageLoader />,
  },
);

export const DemoPage = lazyLoad(
  () => import('./demo'),
  module => module.LeadGeneration,
  {
    fallback: <PageLoader />,
  },
);
