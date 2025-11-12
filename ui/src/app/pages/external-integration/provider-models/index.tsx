import { Helmet } from '@/app/components/helmet';
import { ProviderCard } from '@/app/components/base/cards/provider-card';
import { BluredWrapper } from '@/app/components/wrapper/blured-wrapper';
import { SearchIconInput } from '@/app/components/form/input/IconInput';
import { PageHeaderBlock } from '@/app/components/blocks/page-header-block';
import { PageTitleBlock } from '@/app/components/blocks/page-title-block';
import { cn } from '@/utils';
import { COMPLETE_PROVIDER } from '@/app/components/providers';
/**
 *
 * @returns
 */
export function ProviderModelPage() {
  return (
    <div className={cn('flex flex-col h-full flex-1 overflow-auto')}>
      <Helmet title="Providers and Models" />
      <PageHeaderBlock>
        <div className="flex items-center gap-3">
          <PageTitleBlock>Providers and Models</PageTitleBlock>
          <div className="text-xs opacity-75">
            ({`${COMPLETE_PROVIDER.length}/${COMPLETE_PROVIDER.length}`})
          </div>
        </div>
      </PageHeaderBlock>
      <BluredWrapper className="p-0">
        <SearchIconInput className="bg-light-background" />
      </BluredWrapper>
      <div className="sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 grid p-4">
        {COMPLETE_PROVIDER.map((mp, idx) => {
          return (
            <ProviderCard
              key={`spc-${idx}`}
              provider={mp}
              className="col-span-1"
            />
          );
        })}
      </div>
    </div>
  );
}
