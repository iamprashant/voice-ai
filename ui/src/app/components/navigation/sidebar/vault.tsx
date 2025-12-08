import { Disclosure } from '@/app/components/disclosure';
import { RapidaIcon } from '@/app/components/Icon/Rapida';
import { SidebarIconWrapper } from '@/app/components/navigation/sidebar/sidebar-icon-wrapper';
import { SidebarLabel } from '@/app/components/navigation/sidebar/sidebar-label';
import { SidebarSimpleListItem } from '@/app/components/navigation/sidebar/sidebar-simple-list-item';
import { useSidebar } from '@/context/sidebar-context';
import { cn } from '@/utils';
import { ChevronDown, FolderKey, Key } from 'lucide-react';
import { useState } from 'react';
import { useLocation } from 'react-router-dom';

export function Vault() {
  const location = useLocation();
  const { open } = useSidebar();
  const { pathname } = location;
  const [opt, setOpt] = useState(
    false ||
      pathname.includes('/project-credential') ||
      pathname.includes('/personal-credential'),
  );

  return (
    <li>
      <SidebarSimpleListItem
        className={cn('justify-between')}
        active={opt}
        onClick={() => {
          setOpt(!opt);
        }}
        navigate="#"
      >
        <div className="flex items-center">
          <SidebarIconWrapper>
            <RapidaIcon />
          </SidebarIconWrapper>
          <SidebarLabel>Credentials</SidebarLabel>
        </div>
        <SidebarIconWrapper className="transition-all duration-100">
          <ChevronDown
            className={cn(
              'w-5 h-5 opacity-70 transition-all duration-200',
              opt && 'rotate-180',
            )}
          />
        </SidebarIconWrapper>
      </SidebarSimpleListItem>
      <Disclosure open={opt}>
        <div
          className={cn(
            'ml-6 dark:border-gray-800 border-l',
            open ? 'block' : 'hidden',
          )}
        >
          <SidebarSimpleListItem
            className="mx-0 mr-2"
            active={pathname.includes('/project-credential')}
            navigate="/integration/project-credential"
          >
            <SidebarIconWrapper>
              <FolderKey
                className={cn('w-5 h-5 opacity-75')}
                strokeWidth={1.5}
              />
            </SidebarIconWrapper>
            <SidebarLabel>Project Credential</SidebarLabel>
          </SidebarSimpleListItem>
          <SidebarSimpleListItem
            className="mx-0 mr-2"
            active={pathname.includes('/personal-credential')}
            navigate="/integration/personal-credential"
          >
            <SidebarIconWrapper>
              <Key className={cn('w-5 h-5 opacity-75')} strokeWidth={1.5} />
            </SidebarIconWrapper>
            <SidebarLabel>Personal Token</SidebarLabel>
          </SidebarSimpleListItem>
        </div>
      </Disclosure>
    </li>
  );
}
