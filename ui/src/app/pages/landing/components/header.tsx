import { RapidaIcon } from '@/app/components/Icon/Rapida';
import { RapidaTextIcon } from '@/app/components/Icon/RapidaText';
import { AuthContext } from '@/context/auth-context';
import { ChevronRight } from 'lucide-react';
import { useContext } from 'react';

export const Header = () => {
  //
  const { isAuthenticated } = useContext(AuthContext);

  return (
    <div className="border-b fixed inset-x-0 top-0 z-20 flex h-14 items-center bg-white dark:bg-gray-900">
      <header className="container mx-auto flex items-center px-4 after:-bottom-px sm:px-0">
        <a
          href="/"
          className="flex items-center shrink-0 space-x-1.5 ml-1 text-blue-600 dark:text-blue-500"
        >
          <RapidaIcon className="h-8 w-8" />
          <RapidaTextIcon className="h-6 shrink-0 ml-1" />
        </a>
        <div className="@container flex flex-1 justify-start pl-8"></div>
        <div className="flex items-center gap-5 max-md:hidden lg:gap-6">
          <a
            className="text-sm/6 text-gray-950 dark:text-white hover:text-blue-600"
            href="https://blog.rapida.ai/"
          >
            Blog
          </a>
          <a
            className="text-sm/6 text-gray-950 dark:text-white hover:text-blue-600"
            href="https://doc.rapida.ai/"
          >
            Documentation
          </a>

          <a
            href={
              isAuthenticated && isAuthenticated()
                ? '/dashboard'
                : '/auth/signin'
            }
            className="group relative inline-flex items-center justify-center overflow-hidden border border-blue-600 pl-8 pr-3 font-medium  text-white transition duration-300 ease-out rounded-[2px] h-8 text-sm/6"
          >
            <ChevronRight
              className="absolute w-4 h-4 mr-2 z-10 left-2 my-auto text-white group-hover:text-blue-600 dark:group-hover:text-blue-500"
              strokeWidth={1.5}
            />
            <span className="ease absolute inset-0 flex h-full w-full -translate-x-full items-center justify-center duration-500 group-hover:translate-x-0 text-white group-hover:text-blue-600 dark:group-hover:text-blue-500">
              {isAuthenticated && isAuthenticated() ? 'Dashboard' : 'Sign in'}
            </span>
            <span className="ease absolute left-0 pl-3 right-0 flex h-full w-full transform items-center justify-center transition-all duration-500 group-hover:translate-x-full bg-blue-600">
              Start building
            </span>
            <span className="invisible relative">
              <ChevronRight className="w-4 h-4" strokeWidth={1.5} />
              Start building
            </span>
          </a>
        </div>
      </header>
    </div>
  );
};
