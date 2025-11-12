import { cn } from '@/utils';
import { FC, HTMLAttributes } from 'react';

export const TableBody: FC<HTMLAttributes<HTMLTableSectionElement>> = props => {
  return (
    <tbody {...props} className={cn('text-[15px]', props.className)}>
      {props.children}
    </tbody>
  );
};
