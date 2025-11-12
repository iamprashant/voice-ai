import React, { HTMLAttributes } from 'react';
import { cn } from '@/utils';
import { TableCell } from '@/app/components/base/tables/table-cell';

/**
 *
 * @param props
 * @returns
 */
export function TextCell(props: HTMLAttributes<HTMLDivElement>) {
  return (
    <TableCell>
      <span
        className={cn(
          'font-normal text-left max-w-[20rem] truncate',
          props.className,
        )}
      >
        {props.children}
      </span>
    </TableCell>
  );
}
