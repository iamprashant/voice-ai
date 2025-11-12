import { TableCell } from '@/app/components/base/tables/table-cell';
import { formatFileSize } from '@/utils/format';
import { FC } from 'react';

export const SizeCell: FC<{ size?: string | number }> = ({ size }) => {
  if (size) return <TableCell>{formatFileSize(+size)}</TableCell>;
  return <TableCell>{size}</TableCell>;
};
