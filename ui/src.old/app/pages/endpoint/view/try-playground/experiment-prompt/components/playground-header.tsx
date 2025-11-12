import { PageHeaderBlock } from '@/app/components/blocks/page-header-block';
import { PageTitleBlock } from '@/app/components/blocks/page-title-block';
import { IBlueButton } from '@/app/components/form/button';
import { PlayIcon } from '@/app/components/Icon/Play';
import { Spinner } from '@/app/components/loader/spinner';
import { FC } from 'react';

export const PlaygroundHeader: FC<{
  isValid: boolean;
  loading: boolean;
}> = ({ isValid, loading }) => {
  return (
    <PageHeaderBlock className="border-b h-11">
      <PageTitleBlock>Playground</PageTitleBlock>
      <IBlueButton type="submit">
        Try execute
        {!loading && <PlayIcon className="w-4 h-4 ml-1" strokeWidth={1.5} />}
        {loading && (
          <Spinner className="w-4 h-4 ml-1 border-white flex items-center" />
        )}
      </IBlueButton>
    </PageHeaderBlock>
  );
};
