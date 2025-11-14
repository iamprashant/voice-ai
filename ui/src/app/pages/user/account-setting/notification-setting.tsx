import { FormLabel } from '@/app/components/form-label';
import { IBlueBGArrowButton } from '@/app/components/form/button';
import { InputCheckbox } from '@/app/components/form/checkbox';
import { FieldSet } from '@/app/components/form/fieldset';
import { InputGroup } from '@/app/components/input-group';
import { InputHelper } from '@/app/components/input-helper';
import { connectionConfig } from '@/configs';
import { RAPIDA_SYSTEM_NOTIFICATION } from '@/models/notification';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { PageActionButtonBlock } from '@/app/components/blocks/page-action-button-block';
import {
  UpdateNotificationSettingRequest,
  NotificationSetting as Setting,
  UpdateNotificationSetting,
  ConnectionConfig,
} from '@rapidaai/react';
import { useRapidaStore } from '@/hooks';
import { useCurrentCredential } from '@/hooks/use-credential';
import toast from 'react-hot-toast/headless';

export const NotificationSetting = () => {
  /**
   * loggedin user
   */
  const { token, authId, projectId } = useCurrentCredential();
  /**
   * page error
   */
  const [error, setError] = useState('');

  /**
   * form handling
   */
  const { register, handleSubmit } = useForm();

  /**
   *
   */
  const { loading, showLoader, hideLoader } = useRapidaStore();

  /**
   *
   * @param data
   */
  const onSubmit = (data: any) => {
    setError('');
    showLoader();
    const notificationSettingRequest = new UpdateNotificationSettingRequest();
    const buildEventNotification = (prefix: string, obj: any) => {
      Object.entries(obj).forEach(([key, value]) => {
        const eventNotification = new Setting();
        eventNotification.setChannel('email'); // Example channel, adjust if needed
        eventNotification.setEventtype(prefix ? `${prefix}.${key}` : key); // Use prefix to build event type

        if (typeof value === 'boolean') {
          eventNotification.setEnabled(value);
          notificationSettingRequest.addSettings(eventNotification);
        } else {
          // Recursive case: handle nested objects
          buildEventNotification(prefix ? `${prefix}.${key}` : key, value);
        }
      });
    };
    buildEventNotification('', data);
    UpdateNotificationSetting(
      connectionConfig,
      notificationSettingRequest,
      ConnectionConfig.WithDebugger({
        authorization: token,
        userId: authId,
        projectId: projectId,
      }),
    )
      .then(rlp => {
        hideLoader();
        if (rlp?.getSuccess()) {
          toast.success(
            'The notification setting has been updated successfully.',
          );
        } else {
          let errorMessage = rlp?.getError();
          if (errorMessage) setError(errorMessage.getHumanmessage());
          else {
            setError('Unable to process your request. please try again later.');
          }
          return;
        }
      })
      .catch(e => {
        setError('Unable to process your request. please try again later.');
        hideLoader();
      });
  };

  return (
    <form
      className="pb-20"
      onSubmit={handleSubmit(onSubmit)} // Use the onSubmit handler
    >
      <InputGroup title="Alert and Notifications" className="border-none">
        <hr className="border-t" />
        {RAPIDA_SYSTEM_NOTIFICATION.map(notificationCategory => (
          <div className="p-5 space-y-6" key={notificationCategory.category}>
            <div>
              <legend className="font-semibold ">
                {notificationCategory.category} notifications
              </legend>
              <p className="mt-1 text-sm/6">
                We'll always let you know about important changes, but you pick
                what else you want to hear about.
              </p>
            </div>
            <div className="mt-6 space-y-6 grid grid-cols-4 gap-4">
              {notificationCategory.items.map(item => (
                <div className="flex gap-3" key={item.id}>
                  <div className="flex h-6 shrink-0 items-center">
                    {/* Bind the checkbox with `register` */}
                    <InputCheckbox
                      {...register(item.id)} // Register the field
                      checked={item.default} // Optional initial value
                    />
                  </div>
                  <FieldSet className="text-sm/6">
                    <FormLabel htmlFor={item.id}>{item.label}</FormLabel>
                    <InputHelper id={`${item.id}-description`}>
                      {item.description}
                    </InputHelper>
                  </FieldSet>
                </div>
              ))}
            </div>
          </div>
        ))}
      </InputGroup>
      <PageActionButtonBlock errorMessage={error}>
        <IBlueBGArrowButton
          type="submit"
          className="px-4 rounded-[2px]"
          isLoading={loading}
        >
          Update Notification
        </IBlueBGArrowButton>
      </PageActionButtonBlock>
    </form>
  );
};
