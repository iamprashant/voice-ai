import React from 'react';
import { OptionMenu } from '@/app/components/menu';
import { OptionMenuItem } from '@/app/components/menu/index';

/**
 *
 * @param props
 * @returns
 */
export function UserOption(props: { id: string }) {
  return (
    <OptionMenu
      options={[
        {
          option: 'Edit user',
          onActionClick: () => {},
        },
        {
          option: <OptionMenuItem type="danger">Delete</OptionMenuItem>,
          onActionClick: () => {},
        },
      ]}
    />
  );
}
