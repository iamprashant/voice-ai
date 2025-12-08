import { useSidebar } from '@/context/sidebar-context';
import { cn } from '@/utils';
import { FC, HTMLAttributes } from 'react';

interface AsideProps extends HTMLAttributes<HTMLDivElement> {}
export const Aside: FC<AsideProps> = (props: AsideProps) => {
  const { open, setOpen } = useSidebar();
  return (
    <div
      className={cn(
        'flex flex-col z-12',
        //  top-0 bottom-0
        'no-scrollbar',
        'backdrop-blur-2xl',
        'group',
        'overflow-y-scroll',
        open ? 'w-80' : 'w-14',
        'h-full duration-200 pb-10 hover:bg-white dark:hover:bg-gray-900',
        props.className,
      )}
      onMouseEnter={() => {
        setOpen(true);
      }}
      onMouseLeave={() => {
        setOpen(false);
      }}
    >
      {props.children}
    </div>
  );
};
