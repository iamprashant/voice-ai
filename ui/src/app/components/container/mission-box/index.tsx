import { ActionableHeader } from '@/app/components/navigation/actionable-header';
import { SidebarNavigation } from '@/app/components/navigation/sidebar';
import { Loader } from '@/app/components/loader';
import { useRapidaStore } from '@/hooks';
import { Toast } from '@/app/components/toasts';
import { ProviderContextProvider } from '@/context/provider-context';
import { SidebarProvider } from '@/context/sidebar-context';

/**
 *
 * @param props
 * @returns
 */
export function MissionBox(props: { children?: any }) {
  useRapidaStore();
  return (
    <ProviderContextProvider>
      <SidebarProvider>
        <div className="flex h-[100dvh] relative w-[100dvw]">
          <SidebarNavigation />
          <main className="antialiased text-base text-gray-700 dark:text-gray-400 relative bg-[linear-gradient(103deg,var(--tw-gradient-stops))] from-custom-gray via-custom-pink to-custom-blue font-sans flex-1 flex w-full overflow-hidden">
            <div className="flex w-full absolute top-0 left-0 right-0 z-10">
              <Loader />
            </div>

            <div className="w-full">
              <ActionableHeader />
              <div className="relative h-[calc(100dvh-3rem)] overflow-hidden dark:bg-gray-900 bg-light-background rounded-tl-md border-t border-l flex-1 flex flex-col">
                <Toast />
                {props.children}
              </div>
            </div>
          </main>
        </div>
      </SidebarProvider>
    </ProviderContextProvider>
  );
}
