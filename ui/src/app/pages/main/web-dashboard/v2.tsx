import { useState, useEffect } from 'react';
import { Datepicker } from '@/app/components/Datepicker';
import { Helmet } from '@/app/components/Helmet';
import { ProjectUserGroupAvatar } from '@/app/components/Avatar/ProjectUserGroupAvatar';
import { GetProject } from '@rapidaai/react';
import { useCredential } from '@/hooks/use-credential';
import { GetProjectResponse } from '@rapidaai/react';
import { useRapidaStore } from '@/hooks';
import { AnimatedLinkButton } from '@/app/components/Form/Button/AnimateLinkButton';
import { User } from '@rapidaai/react';
import toast from 'react-hot-toast/headless';
import {
  TokenOutputSpeed,
  PeriodParams,
  TokenUsages,
  AverageResponseTime,
} from '@/app/pages/main/web-dashboard/components/app-chart';
import moment from 'moment';
import { useActivityLogPage } from '@/hooks/use-activity-log-page-store';
import { toDateString } from '@/utils';
import { ServiceError } from '@rapidaai/react';
import { connectionConfig } from '@/configs';
export function HomePage() {
  const [userId, token, projectId] = useCredential();
  const [members, setMembers] = useState<User[]>([]);

  const {
    getActivities,
    addCriteria,
    activities,
    onChangeActivities,
    setPageSize,
  } = useActivityLogPage();

  //
  const initializePeriodParams = (): PeriodParams => {
    const start = moment().subtract(30, 'days').toDate();
    const end = moment().add(1, 'days').toDate();
    return { start, end };
  };

  //
  const [selectedDates, setSelectedDates] = useState<PeriodParams>(
    initializePeriodParams,
  );

  const { showLoader, hideLoader, loading } = useRapidaStore();
  //   const [activities, setActivities] = useState<AuditLog[]>([]);

  const afterGetProject = (
    err: ServiceError | null,
    gur: GetProjectResponse | null,
  ) => {
    if (err) {
      return;
    }
    if (gur?.getSuccess()) {
      let members = gur.getData()?.getMembersList();
      if (members) setMembers(members);
    }
  };

  useEffect(() => {
    if (projectId) {
      GetProject(connectionConfig, projectId, afterGetProject, {
        authorization: token,
        'x-auth-id': userId,
        'x-project-id': projectId,
      });
    }
  }, [projectId, token, userId]);

  const onDateSelect = (to: Date, from: Date) =>
    setSelectedDates({
      start: from,
      end: to,
    });

  useEffect(() => {
    if (selectedDates) {
      setPageSize(500);
      showLoader();
      addCriteria('created_date', toDateString(selectedDates.start), '>=');
      addCriteria('created_date', toDateString(selectedDates.end), '<=');

      getActivities(
        projectId,
        token,
        userId,
        err => {
          hideLoader();
          toast.error(err);
        },
        logs => {
          hideLoader();
          onChangeActivities(logs);
        },
      );
    }
  }, [projectId, token, userId, JSON.stringify(selectedDates)]);

  return (
    <>
      <Helmet title="Dashboard"></Helmet>
      <div className="px-4 sm:px-6 lg:px-8 py-8 w-full max-w-9xl mx-auto">
        <div className="sm:flex sm:justify-between sm:items-center mb-8">
          {/* Left: Avatars */}
          {projectId && (
            <ProjectUserGroupAvatar
              projectId={projectId}
              members={members.map(m => ({ name: m.getName() }))}
            />
          )}

          <div className="grid grid-flow-col sm:auto-cols-max items-center justify-center sm:justify-end gap-2">
            <Datepicker onDateSelect={onDateSelect} />
          </div>
        </div>
        <div className="grid gap-6 grid-cols-1 xl:grid-cols-2 w-full mb-6">
          <TokenUsages
            {...selectedDates}
            loading={loading}
            activities={activities}
            className="col-span-1"
          />
          <TokenOutputSpeed
            {...selectedDates}
            className="col-span-1"
            loading={loading}
            activities={activities}
          />
          <AverageResponseTime
            {...selectedDates}
            className="col-span-2"
            loading={loading}
            activities={activities}
          />
        </div>
      </div>
    </>
  );
}

function OnboardingCard() {
  return (
    <div className="flex flex-col items-center px-8 min-h-[95vh]">
      <div className="mt-10 sm:max-w-3xl w-full dark:bg-gray-950 rounded-[2px] ring-1 ring-gray-900/5 bg-white">
        <div className="bg-white dark:bg-gray-950 flex justify-between items-center p-8 rounded-t-md border-b dark:border-gray-800">
          <div>
            <h3 className="font-semibold text-lg">
              Your dashboard is not ready yet.
            </h3>
            <p className="text-sm mt-2">
              You don't have enough data to build the dashboard and provide the
              insight, some of helpful links to get started with rapidaAI.
            </p>
          </div>
        </div>
        <div className="bg-gray-50 shadow-sm border dark:border-gray-950 dark:bg-gray-900 flex m-2 p-4 pt-8 rounded-[2px] flex-col ">
          <h3 className="font-semibold text-lg">Create a provider key</h3>
          <p className="text-sm mt-4 opacity-80">
            The providers keys are secured with AWS credential manager, It will
            be used to authenticate with your provider on your behalf.
          </p>
          <div className="flex justify-end">
            <AnimatedLinkButton
              className="text-sm mt-4"
              target="_self"
              iconsize="w-4 h-4"
              to="/integration/vault"
              text="Create Provider Key"
            />
          </div>
        </div>
        <div className="bg-gray-50 shadow-sm border dark:border-gray-950 dark:bg-gray-900 flex m-2 p-4 pt-8 rounded-[2px] flex-col ">
          <h3 className="font-semibold text-lg">
            Tryout your first prompt experiment
          </h3>
          <p className="text-sm mt-4 opacity-80">
            A utility-first experiments framework packed with model providers
            like openAI, cohere, anthropic and other open-source models.
          </p>
          <div className="flex justify-end">
            <AnimatedLinkButton
              target="_blank"
              className="text-sm mt-4"
              iconsize="w-4 h-4"
              to="https://experiments.rapida.ai/"
              text="Go to experimentation"
            />
          </div>
        </div>
        <div className="bg-gray-50 shadow-sm border dark:border-gray-950 dark:bg-gray-900 flex m-2 p-4 pt-8 rounded-[2px] flex-col ">
          <h3 className="font-semibold text-lg">
            Invite and manage your team members
          </h3>
          <p className="text-sm mt-4 opacity-80">
            You can invite users users to an organization or to a project with
            specific role permissions.
          </p>
          <div className="flex justify-end">
            <AnimatedLinkButton
              target="_self"
              className="text-sm mt-4"
              iconsize="w-4 h-4"
              to="/organization/users"
              text="Go to users and teams"
            />
          </div>
        </div>
        <p className="px-2 md:px-5 py-3 pb-6 text-sm">
          Please{' '}
          <a className="underline" href="mailto:mukesh@rapida.ai">
            click here
          </a>{' '}
          to reachout to us incase any help.
        </p>
      </div>
    </div>
  );
}
