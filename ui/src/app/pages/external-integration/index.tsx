import { lazyLoad } from '@/utils/loadable';
import { PageLoader } from '@/app/components/loader/page-loader';

export const IntegrationModelPage = lazyLoad(
  () => import('./provider-models'),
  module => module.ProviderModelPage,
  {
    fallback: <PageLoader />,
  },
);

export const IntegrationModelInformationPage = lazyLoad(
  () => import('./provider-models/information'),
  module => module.ProviderModelInformationPage,
  {
    fallback: <PageLoader />,
  },
);

export const IntegrationProjectCredentialPage = lazyLoad(
  () => import('./rapida-credentials'),
  module => module.ProjectCredentialPage,
  {
    fallback: <PageLoader />,
  },
);

export const IntegrationPersonalCredentialPage = lazyLoad(
  () => import('./rapida-credentials'),
  module => module.PersonalCredentialPage,
  {
    fallback: <PageLoader />,
  },
);
