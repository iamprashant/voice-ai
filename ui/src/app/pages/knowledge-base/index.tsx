import { lazyLoad } from '@/utils/loadable';
import { LineLoader } from '@/app/components/loader/line-loader';

export const KnowledgePage = lazyLoad(
  () => import('./listing'),
  module => module.KnowledgePage,
  {
    fallback: <LineLoader />,
  },
);

// create-knowledge
export const KnowledgeCreateKnowledgePage = lazyLoad(
  () => import('./action/create-knowledge'),
  module => module.CreateKnowledgePage,
  {
    fallback: <LineLoader />,
  },
);

export const KnowledgeViewKnowledgePage = lazyLoad(
  () => import('./view'),
  module => module.ViewKnowledgePage,
  {
    fallback: <LineLoader />,
  },
);

export const KnowledgeAddNewKnowledgeFilePage = lazyLoad(
  () => import('./action/create-knowledge-document'),
  module => module.CreateKnowledgeDocumentPage,
  {
    fallback: <LineLoader />,
  },
);

export const KnowledgeAddNewStructureDocumentPage = lazyLoad(
  () => import('./action/create-knowledge-document'),
  module => module.CreateKnowledgeStructureDocumentPage,
  {
    fallback: <LineLoader />,
  },
);
