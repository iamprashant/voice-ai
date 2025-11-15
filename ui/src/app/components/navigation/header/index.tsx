import React, { FC, Fragment, HTMLAttributes, memo, useContext } from 'react';
import { cn } from '@/utils';
import {
  Menu,
  MenuButton,
  MenuItem,
  MenuItems,
  Transition,
} from '@headlessui/react';
import { RightArrowIcon } from '@/app/components/Icon/RightArrow';
import { AuthContext } from '@/context/auth-context';
import { useWorkspace } from '@/workspace';
import { RapidaIcon } from '@/app/components/Icon/Rapida';
import { RapidaTextIcon } from '@/app/components/Icon/RapidaText';

interface HeaderProps extends HTMLAttributes<HTMLDivElement> {}
export const Header: FC<HeaderProps> = memo(({ ...attr }) => {
  return (
    <nav
      className={cn(
        'py-3 sticky top-0 z-10 w-full backdrop-blur-sm flex-none transition-colors duration-500 border-b bg-white dark:bg-gray-950 dark:border-gray-900  px-4 sm:px-0',
        attr.className,
      )}
    >
      <div
        className={cn(
          'mx-auto container',
          'text-gray-700 transition data-disabled:text-gray-400 dark:data-disabled:text-gray-700 duration-200 dark:text-gray-300',
        )}
      >
        <HeaderContent />
      </div>
    </nav>
  );
});

function HeaderContent() {
  const workspace = useWorkspace();
  return (
    <div className="flex justify-between items-center">
      <div className="flex items-center justify-center space-x-1 w-fit">
        {workspace.logo ? (
          <>
            <img
              src={workspace.logo.light}
              alt={workspace.title}
              className="h-8 block dark:hidden"
            />
            <img
              src={workspace.logo.dark}
              alt={workspace.title}
              className="h-8 hidden dark:block"
            />
          </>
        ) : (
          <div className="flex items-center shrink-0 space-x-1.5 ml-1 text-blue-600 dark:text-blue-500">
            <RapidaIcon className="h-8 w-8" />
            <RapidaTextIcon className="h-6" />
          </div>
        )}
      </div>
      <MobileMenu />
    </div>
  );
}

function MobileMenu() {
  const { isAuthenticated } = useContext(AuthContext);
  return (
    <Menu as="div" className="inline-block text-left md:hidden">
      {({ open }) => (
        <>
          <MenuButton
            className={cn(
              'block md:hidden',
              'relative ml-auto flex w-fit items-center',
              'focus:outline-hidden',
              'rounded-[2px]',
            )}
          >
            {!open ? (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                />
              </svg>
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M6 18 18 6M6 6l12 12"
                />
              </svg>
            )}
          </MenuButton>
          <Transition
            as={Fragment}
            enter="transition ease-out duration-100"
            enterFrom="transform opacity-0 scale-95"
            enterTo="transform opacity-100 scale-100"
            leave="transition ease-in duration-75"
            leaveFrom="transform opacity-100 scale-100"
            leaveTo="transform opacity-0 scale-95"
          >
            <MenuItems
              className={cn(
                'focus:outline-hidden',
                'absolute right-0 top-full w-full origin-top-right bg-white shadow-lg',
                'z-20 backdrop-blur-3xl bg-white/90 dark:bg-gray-950/90',
              )}
            >
              <div className="px-4 py-4 space-y-2">
                <MenuItem>
                  {({ active }) => (
                    <a
                      href="/"
                      className="hover:bg-gray-950/30 font-medium hover:text-white dark:hover:bg-white/30 focus:ring-1 focus:ring-gray-950/10 dark:focus:ring-white/20 dark:data-[state=open]:text-white data-[state=open]:bg-gray-950/5 dark:data-[state=open]:bg-white/5 data-[state=open]:text-gray-950 group flex select-none items-center justify-between gap-1 rounded-[2px] px-3 py-2.5 leading-none outline-hidden"
                    >
                      Product
                    </a>
                  )}
                </MenuItem>
                <MenuItem>
                  {({ active }) => (
                    <a
                      target="_blank"
                      href="https://blog.rapida.ai"
                      className="hover:bg-gray-950/30 font-medium hover:text-white dark:hover:bg-white/30 focus:ring-1 focus:ring-gray-950/10 dark:focus:ring-white/20 dark:data-[state=open]:text-white data-[state=open]:bg-gray-950/5 dark:data-[state=open]:bg-white/5 data-[state=open]:text-gray-950 group flex select-none items-center justify-between gap-1 rounded-[2px] px-3 py-2.5 leading-none outline-hidden"
                      rel="noreferrer"
                    >
                      Blog
                    </a>
                  )}
                </MenuItem>
                <MenuItem>
                  {({ active }) => (
                    <a
                      href={
                        isAuthenticated && isAuthenticated()
                          ? '/dashboard'
                          : '/auth/signin'
                      }
                      className="hover:bg-gray-950/30 font-medium hover:text-white dark:hover:bg-white/30 focus:ring-1 focus:ring-gray-950/10 dark:focus:ring-white/20 dark:data-[state=open]:text-white data-[state=open]:bg-gray-950/5 dark:data-[state=open]:bg-white/5 data-[state=open]:text-gray-950 group flex select-none items-center justify-between gap-1 rounded-[2px] px-3 py-2.5 leading-none outline-hidden"
                    >
                      Sign in
                    </a>
                  )}
                </MenuItem>
                <MenuItem>
                  {({ active }) => (
                    <a
                      className={cn(
                        'relative px-4 flex h-9 w-fit items-center justify-center before:absolute before:inset-0 before:rounded-[2px] before:transition-transform before:duration-300 hover:before:scale-105 active:duration-75 active:before:scale-95 dark:before:border-gray-600 sm:px-4 before:border before:border-blue-600 before:bg-gray-100 dark:before:bg-gray-800',
                      )}
                      target="_blank"
                      href="https://calendly.com/rapida-ai/30min"
                      rel="noreferrer"
                    >
                      <span className="relative font-medium text-blue-600 dark:text-white">
                        Book a Demo
                      </span>
                      <span className="relative ml-2">
                        <RightArrowIcon />
                      </span>
                    </a>
                  )}
                </MenuItem>
              </div>
            </MenuItems>
          </Transition>
        </>
      )}
    </Menu>
  );
}
