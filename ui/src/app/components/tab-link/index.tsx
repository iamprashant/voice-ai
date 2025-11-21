import { cn } from '@/utils';
import { FC, HTMLAttributes } from 'react';
import { NavLink } from 'react-router-dom';

interface TabProps extends HTMLAttributes<HTMLDivElement> {
  isActive?: boolean;
}

interface LinkTabProps extends TabProps {
  to: string;
}

export const Tab: FC<TabProps> = ({ isActive, children, ...props }) => {
  return (
    <div
      {...props}
      className={cn(
        'group px-2 border-b-[3px] border-transparent -mb-[0.2rem] cursor-pointer',
        isActive
          ? 'text-blue-500 bg-blue-500/10'
          : 'hover:bg-blue-500/5 hover:text-blue-500',
      )}
    >
      <div className="capitalize px-3 py-3 group-hover:bg-blue-600/5 dark:group-hover:bg-blue-950/50">
        {children}
      </div>
    </div>
  );
};

export const TabLink: FC<LinkTabProps> = ({ to, children }) => {
  return (
    <NavLink
      to={to}
      className={({ isActive }) => {
        return cn(
          'group cursor-pointer hover:bg-gray-500/10',
          isActive
            ? 'text-blue-500 bg-blue-500/10'
            : 'hover:bg-blue-500/5 hover:text-blue-500',
        );
      }}
    >
      <div className="px-6 py-2 font-medium whitespace-nowrap tracking-wide text-pretty">
        {children}
      </div>
    </NavLink>
  );
};
export const SideTabLink: FC<LinkTabProps> = props => {
  return (
    <NavLink
      to={props.to}
      className={({ isActive }) =>
        cn(
          'group px-2 border-r-[3px] border-transparent -ms-[0.1rem] cursor-pointer font-medium text-[14.5px] whitespace-nowrap tracking-wide text-pretty',
          'flex items-center px-5 py-2 relative',

          isActive
            ? 'text-blue-500 bg-blue-500/10'
            : 'hover:bg-blue-500/5 hover:text-blue-500',
          props.className,
        )
      }
    >
      {props.children}
    </NavLink>
  );
};

export const SideTab: FC<LinkTabProps> = props => {
  return (
    <div
      onClick={props.onClick}
      className={cn(
        'group px-2 border-r-[3px] border-transparent -ms-[0.1rem] cursor-pointer font-medium text-[14.5px] whitespace-nowrap tracking-wide text-pretty',
        'flex items-center px-5 py-2 relative',
        props.className,
        props.isActive === true
          ? 'text-blue-500 bg-blue-500/10'
          : 'hover:bg-blue-500/5 hover:text-blue-500',
      )}
    >
      {props.children}
    </div>
  );
};
