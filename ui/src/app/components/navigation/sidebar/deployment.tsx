import { memo } from 'react';
import { cn } from '@/utils';
import { SidebarIconWrapper } from '@/app/components/navigation/sidebar/sidebar-icon-wrapper';
import { SidebarLabel } from '@/app/components/navigation/sidebar/sidebar-label';
import { SidebarSimpleListItem } from '@/app/components/navigation/sidebar/sidebar-simple-list-item';
import { useLocation } from 'react-router-dom';
import { Tooltip } from '@/app/components/tooltip';
import { BetaIcon } from '@/app/components/Icon/Beta';
import { BotMessageSquare, Route } from 'lucide-react';

export const Deployment = memo(() => {
  const location = useLocation();
  const { pathname } = location;

  return (
    <li>
      <SidebarSimpleListItem
        active={pathname.includes('/deployment/assistant')}
        navigate="/deployment/assistant"
      >
        <SidebarIconWrapper>
          <BotMessageSquare
            className={cn('w-5 h-5 opacity-75')}
            strokeWidth={1.5}
          />
        </SidebarIconWrapper>
        <SidebarLabel>
          Assistants
          <Tooltip
            children={
              <p className="text-xs">
                We are working very hard <br />
                to get you best experience.
                <br />
              </p>
            }
            icon={<BetaIcon />}
          ></Tooltip>
        </SidebarLabel>
      </SidebarSimpleListItem>
      <SidebarSimpleListItem
        active={pathname.includes('/deployment/endpoint')}
        navigate="/deployment/endpoint"
      >
        <SidebarIconWrapper>
          <Route className={cn('w-5 h-5 opacity-75')} strokeWidth={1.5} />
        </SidebarIconWrapper>
        <SidebarLabel>Endpoints</SidebarLabel>
      </SidebarSimpleListItem>
    </li>
  );
});
