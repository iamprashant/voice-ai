import { DiscoverIcon } from '@/app/components/Icon/Discover';
import { SidebarIconWrapper } from '@/app/components/navigation/sidebar/sidebar-icon-wrapper';
import { SidebarLabel } from '@/app/components/navigation/sidebar/sidebar-label';
import { SidebarSimpleListItem } from '@/app/components/navigation/sidebar/sidebar-simple-list-item';
import { useLocation } from 'react-router-dom';

export function Discover() {
  const location = useLocation();
  const { pathname } = location;

  /**
   *
   */
  return (
    <li>
      <SidebarSimpleListItem active={pathname.includes('/hub')} navigate="/hub">
        <SidebarIconWrapper>
          <DiscoverIcon />
        </SidebarIconWrapper>
        <SidebarLabel>Hub</SidebarLabel>
      </SidebarSimpleListItem>
    </li>
  );
}
