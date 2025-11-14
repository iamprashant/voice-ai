import { cn } from '@/utils';
import { HTMLAttributes } from 'react';

export function TabBody(props: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'bg-transparent flex-1 grow w-full overflow-auto',
        props.className,
      )}
    >
      {props.children}
    </div>
  );
}
