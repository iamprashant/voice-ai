import { Observability } from '@/app/components/navigation/sidebar/observability';
import { Deployment } from '@/app/components/navigation/sidebar/deployment';
import { Dashboard } from '@/app/components/navigation/sidebar/dashboard';
import { Team } from '@/app/components/navigation/sidebar/team';
import { Project } from '@/app/components/navigation/sidebar/project';
import { Vault } from '@/app/components/navigation/sidebar/vault';
import { SidebarLabel } from '@/app/components/navigation/sidebar/sidebar-label';
import { Knowledge } from '@/app/components/navigation/sidebar/knowledge';
import { Aside } from '@/app/components/aside';
import { ExternalTool } from '@/app/components/navigation/sidebar/external-tools';
import { useWorkspace } from '@/workspace';
import { RapidaIcon } from '@/app/components/Icon/Rapida';
import { RapidaTextIcon } from '@/app/components/Icon/RapidaText';

/**
 *
 * @param props
 * @returns
 */
export function SidebarNavigation(props: {}) {
  const workspace = useWorkspace();
  return (
    <Aside className="space-y-2 ">
      <div className="flex rounded-[2px] my-2 text-blue-600 items-center">
        <div className="pl-[0.8rem] py-2.5 shrink-0">
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
      </div>
      <div className="space-y-4">
        <ul className="space-y-1">
          <Dashboard />
          <Observability />
          <Deployment />
          <Knowledge />
        </ul>
        <div className="space-y-3">
          <SidebarLabel className="uppercase truncate pl-3 text-xs opacity-0 group-hover:opacity-100">
            Integrations
          </SidebarLabel>
          <ul className="space-y-1 mt-2">
            <ExternalTool />
            <Vault />
          </ul>
        </div>
        <div className="space-y-3">
          <SidebarLabel className="uppercase truncate pl-3 text-xs opacity-0 group-hover:opacity-100">
            Organizations
          </SidebarLabel>
          <ul className="space-y-1  mt-2">
            <Team />
            <Project />
          </ul>
        </div>
      </div>
    </Aside>
  );
}
