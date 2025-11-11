import React, { FC, HTMLAttributes } from 'react';
import { cn } from '@/utils';

export const ModalTitleBlock: FC<HTMLAttributes<HTMLDivElement>> = props => {
  return (
    <div className={cn('text-lg font-medium', props.className)}>
      {props.children}
    </div>
  );
};
