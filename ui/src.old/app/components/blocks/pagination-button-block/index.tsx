import React, { FC, HTMLAttributes } from 'react';
import { cn } from '@/utils';

export const PaginationButtonBlock: FC<
  HTMLAttributes<HTMLDivElement>
> = props => {
  return (
    <div
      className={cn(
        'flex flex-row divide-x dark:divide-gray-800',
        props.className,
      )}
    >
      {props.children}
    </div>
  );
};
