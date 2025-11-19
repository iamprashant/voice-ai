import { allProvider, RapidaProvider } from '@/providers';
import { cn } from '@/utils';
import { FC, HTMLAttributes, useEffect, useState } from 'react';

/**
 *
 */
interface ProviderPillProps extends HTMLAttributes<HTMLSpanElement> {
  provider?: string;
}

/**
 *
 * @param props
 * @returns
 */
export const ProviderPill: FC<ProviderPillProps> = props => {
  //
  const [currentProvider, setcurrentProvider] = useState<RapidaProvider | null>(
    null,
  );

  useEffect(() => {
    if (props.provider) {
      let cModel = allProvider().find(
        x => x.code.toLowerCase() === props.provider?.toLowerCase(),
      );
      if (cModel) setcurrentProvider(cModel);
    }
  }, [props.provider]);

  return (
    <span
      onClick={props.onClick}
      className={cn(
        'px-2 py-1 truncate',
        'items-center bg-blue-400/10 dark:bg-blue-800/10 rounded-[2px]',
        'ring-blue-500/30 dark:ring-blue-400/20 ring-[0.7px]',
        'flex items-center justify-center w-fit',
        props.className,
      )}
    >
      <img
        alt={currentProvider?.name}
        src={currentProvider?.image}
        className="w-4 h-4 mr-1.5 inline-block"
      />
      <span className="font-medium opacity-80 ">{currentProvider?.name}</span>
    </span>
  );
};
