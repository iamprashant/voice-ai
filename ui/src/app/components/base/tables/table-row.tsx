import { cn } from '@/utils';

export function TableRow(props: React.HTMLAttributes<HTMLTableRowElement>) {
  return (
    <tr
      {...props}
      className={cn(
        'dark:divide-gray-800 divide-gray-200 dark:border-gray-800 border-gray-200 border-b-[0.5px] divide-x hover:bg-gray-50 dark:hover:bg-gray-950/20 text-sm/6 text-pretty hover:dark:text-gray-400',
        props.className,
      )}
    >
      {props.children}
    </tr>
  );
}
