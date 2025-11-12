import { TableCell } from '@/app/components/base/tables/table-cell';
import { formatHumanReadableNumber } from '@/utils/format';
import React, { useEffect, useState } from 'react';

/**
 *
 * @param props
 * @returns
 */
export function NumberCell(props: { num?: string }) {
  const [num, setNum] = useState('');

  useEffect(() => {
    if (!props.num) {
      setNum('0');
      return;
    }
    setNum(formatHumanReadableNumber(props.num));
  }, [props.num]);
  return (
    <TableCell className="text-center">
      <span className="font-medium max-w-[20rem] truncate">{num}</span>
    </TableCell>
  );
}
