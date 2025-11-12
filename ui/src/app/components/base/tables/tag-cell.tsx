import { TableCell } from '@/app/components/base/tables/table-cell';
import { MultiplePills } from '@/app/components/pill';
import React from 'react';

export function TagCell(props: { tags?: string[] }) {
  return (
    <TableCell className="flex">
      {props.tags && props.tags.length > 0 ? (
        <MultiplePills tags={props.tags} />
      ) : (
        <>no tags</>
      )}
    </TableCell>
  );
}
