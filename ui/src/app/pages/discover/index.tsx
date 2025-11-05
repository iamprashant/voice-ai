import { PageLoader } from '@/app/components/Loader/page-loader';
import { lazyLoad } from '@/utils/loadable';

export const DiscoverPage = lazyLoad(
  () => import('./marketplace'),
  module => module.DiscoverPage,
  {
    fallback: <PageLoader />,
  },
);
