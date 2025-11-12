import React, { HTMLAttributes } from 'react';
import { cn } from '@/utils';

export function InfoHeading(props: HTMLAttributes<HTMLSpanElement>) {
  return (
    <span className={cn('text-base ml-2', props.className)}>
      {props.children}
    </span>
  );
}
