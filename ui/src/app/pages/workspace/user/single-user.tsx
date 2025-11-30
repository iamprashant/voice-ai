import React from 'react';
import { UserOption } from './user-options';
import { TextImage } from '@/app/components/text-image';
import { User } from '@rapidaai/react';
import { StatusIndicator } from '@/app/components/indicators/status';
import { RoleIndicator } from '@/app/components/indicators/role';
import { toHumanReadableDate } from '@/utils/date';
import { TableRow } from '@/app/components/base/tables/table-row';
import { TableCell } from '@/app/components/base/tables/table-cell';

/**
 * all the user row
 * @param props
 * @returns
 */
export function SingleUser(props: { user: User }) {
  return (
    <TableRow>
      <TableCell>{props.user.getId()}</TableCell>
      <TableCell>
        <div className="flex items-center">
          <div className="shrink-0 mr-3">
            <TextImage size={7} name={props.user.getName()}></TextImage>
          </div>
          <div className="">{props.user.getName()}</div>
        </div>
      </TableCell>
      <TableCell>
        <div className="font-normal text-left max-w-[20rem] truncate">
          {props.user.getEmail()}
        </div>
      </TableCell>
      <TableCell>
        <RoleIndicator role={'SUPER_ADMIN'} />
      </TableCell>
      <TableCell>
        <div className="font-normal text-left underline decoration-dotted">
          {toHumanReadableDate(props.user.getCreateddate()!)}
        </div>
      </TableCell>
      <TableCell>
        <StatusIndicator state={props.user.getStatus()} />
      </TableCell>
      <TableCell>
        <UserOption id={props.user.getId()}></UserOption>
      </TableCell>
    </TableRow>
  );
}
