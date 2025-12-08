import { SidebarIconWrapper } from '@/app/components/navigation/sidebar/sidebar-icon-wrapper';
import { SidebarLabel } from '@/app/components/navigation/sidebar/sidebar-label';
import { SidebarSimpleListItem } from '@/app/components/navigation/sidebar/sidebar-simple-list-item';
import { cn } from '@/utils';
import { FolderGit } from 'lucide-react';
import React from 'react';
import { useLocation } from 'react-router-dom';

export function Project() {
  const location = useLocation();
  const { pathname } = location;
  const currentPath = '/organization/projects';
  return (
    <li>
      <SidebarSimpleListItem
        navigate={currentPath}
        active={pathname.includes(currentPath)}
      >
        <SidebarIconWrapper>
          <FolderGit className={cn('w-5 h-5 opacity-75')} strokeWidth={1.5} />
        </SidebarIconWrapper>
        <SidebarLabel>Projects</SidebarLabel>
      </SidebarSimpleListItem>
    </li>
  );
}
