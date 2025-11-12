import { TableCell } from '@/app/components/base/tables/table-cell';
import React, { useEffect, useState } from 'react';

/**
 *
 * @param props
 * @returns
 */
export function CostCell(props: { cost?: number }) {
  const [cost, setCost] = useState('');

  useEffect(() => {
    if (!props.cost) {
      setCost('0');
      return;
    }
    setCost(
      new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
      }).format(props.cost),
    );
  }, [props.cost]);
  return (
    <TableCell className="text-center">
      <span className="font-medium max-w-[20rem] truncate">{cost}</span>
    </TableCell>
  );
}
