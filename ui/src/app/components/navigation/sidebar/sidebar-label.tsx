import React, { HTMLAttributes } from 'react';
import { cn } from '@/utils';
import { useSidebar } from '@/context/sidebar-context';

export function SidebarLabel(props: HTMLAttributes<HTMLSpanElement>) {
  const { open } = useSidebar();
  return (
    <span
      className={cn(
        'text-sm truncate inline-block transition-all duration-200 text-gray-900 dark:text-gray-100',
        open ? 'opacity-100' : 'opacity-0',
        props.className,
      )}
      {...props}
    >
      {props.children}
    </span>
  );
}
