import { lazyLoad } from '@/utils/loadable';
import { PageLoader } from '@/app/components/loader/page-loader';

export const DashboardHomePage = lazyLoad(
  () => import('./web-dashboard'),
  module => module.HomePage,
  {
    fallback: <PageLoader />,
  },
);
