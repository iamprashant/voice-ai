import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { toHumanReadableDate, toHumanReadableRelativeTime } from '@/utils/date';
import { TableCell } from '@/app/components/base/tables/table-cell';

/**
 *
 * @param props
 * @returns
 */
export function RelativeDateCell(props: { date?: Timestamp }) {
  return (
    <TableCell>
      <div className="font-normal text-left underline decoration-dotted">
        {props.date && toHumanReadableRelativeTime(props.date)}
      </div>
    </TableCell>
  );
}

/**
 *
 * @param props
 * @returns
 */
export function DateCell(props: { date?: Timestamp }) {
  return (
    <TableCell>
      <div className="font-normal text-left underline decoration-dotted">
        {props.date && toHumanReadableDate(props.date)}
      </div>
    </TableCell>
  );
}
