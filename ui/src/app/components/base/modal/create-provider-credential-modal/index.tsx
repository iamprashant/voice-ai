import { useEffect, useState } from 'react';
import {
  ConnectionConfig,
  CreateProviderCredentialRequest,
} from '@rapidaai/react';
import { useCurrentCredential } from '@/hooks/use-credential';
import { CreateProviderKey } from '@rapidaai/react';
import { Input } from '@/app/components/form/input';
import { ProviderDropdown } from '@/app/components/dropdown/provider-dropdown';
import { ErrorMessage } from '@/app/components/form/error-message';
import { useRapidaStore } from '@/hooks';
import toast from 'react-hot-toast/headless';
import { GenericModal, ModalProps } from '@/app/components/base/modal';
import { FieldSet } from '@/app/components/form/fieldset';
import {
  IBlueBGArrowButton,
  ICancelButton,
} from '@/app/components/form/button';
import { ModalFooter } from '@/app/components/base/modal/modal-footer';
import { ModalBody } from '@/app/components/base/modal/modal-body';
import { FormLabel } from '@/app/components/form-label';
import { ModalHeader } from '@/app/components/base/modal/modal-header';
import { ModalTitleBlock } from '@/app/components/blocks/modal-title-block';
import { connectionConfig } from '@/configs';
import { useProviderContext } from '@/context/provider-context';
import { Textarea } from '@/app/components/form/textarea';
import { Struct } from 'google-protobuf/google/protobuf/struct_pb';
import { INTEGRATION_PROVIDER, RapidaProvider } from '@/providers';
import { ModalFitHeightBlock } from '@/app/components/blocks/modal-fit-height-block';

/**
 * creation provider key dialog props that gives ability for opening and closing modal props
 */
interface CreateProviderCredentialDialogProps extends ModalProps {
  /**
   * exiting provider if there
   */
  currentProvider?: string | null;
}
/**
 *
 * to create a provider key for given model
 * @param props
 * @returns
 */

export function CreateProviderCredentialDialog(
  props: CreateProviderCredentialDialogProps,
) {
  const { authId, projectId, token } = useCurrentCredential();
  const [provider, setProvider] = useState<RapidaProvider | null>();
  const providerCtx = useProviderContext();

  const { loading, showLoader, hideLoader } = useRapidaStore();
  const [error, setError] = useState('');
  const [keyName, setKeyName] = useState('');
  const [config, setConfig] = useState<Record<string, string>>({});

  useEffect(() => {
    setProvider(
      INTEGRATION_PROVIDER.slice()
        .reverse()
        .find(x => x.code === props.currentProvider),
    );
  }, [props.currentProvider]);

  const handleInputChange = (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = event.target;
    if (name.startsWith('config.')) {
      setConfig(prev => ({
        ...prev,
        [name.replace('config.', '')]: value,
      }));
    } else if (name === 'keyName') {
      setKeyName(value);
    }
  };

  const validateAndSubmit = () => {
    if (!props.currentProvider) {
      setError('Please select the provider which you want to create the key.');
      return;
    }

    if (!provider) {
      setError('Please select the provider which you want to create the key.');
      return;
    }

    if (!keyName.trim()) {
      setError('Please provide a valid key name for the credential.');
      return;
    }

    const missingFields = provider.configurations?.filter(
      configOption => !config[configOption.name]?.trim(),
    );

    if (missingFields && missingFields.length > 0) {
      setError(
        `Please fill out the following fields: ${missingFields
          .map(field => field.label)
          .join(', ')}`,
      );
      return;
    }

    // Proceed with creating the provider key
    showLoader();
    const requestObject = new CreateProviderCredentialRequest();
    requestObject.setProvider(provider.code);
    requestObject.setCredential(Struct.fromJavaScript(config));
    requestObject.setName(keyName);

    CreateProviderKey(
      connectionConfig,
      requestObject,
      ConnectionConfig.WithDebugger({
        authorization: token,
        userId: authId,
        projectId: projectId,
      }),
    )
      .then(cpkr => {
        hideLoader();
        if (cpkr?.getSuccess()) {
          toast.success(
            'Provider credential have been successfully added to the vault.',
          );
          providerCtx.reloadProviderCredentials();
          props.setModalOpen(false);
          setError('');
          setKeyName('');
          setConfig({});
        } else {
          let errorMessage = cpkr?.getError();
          setError(
            errorMessage?.getHumanmessage() ??
              'Unable to process your request. Please try again later.',
          );
        }
      })
      .catch(() => {
        hideLoader();
        toast.error(
          'Unable to create provider credential, please try again later.',
        );
      });
  };

  return (
    <GenericModal modalOpen={props.modalOpen} setModalOpen={props.setModalOpen}>
      <ModalFitHeightBlock>
        <ModalHeader
          onClose={() => {
            props.setModalOpen(false);
          }}
        >
          <ModalTitleBlock>Create provider credential</ModalTitleBlock>
        </ModalHeader>
        <ModalBody>
          <FieldSet>
            <FormLabel>Select your provider</FormLabel>
            <ProviderDropdown
              currentProvider={provider ? provider : undefined}
              setCurrentProvider={p => {
                setError('');
                setKeyName('');
                setConfig({});
                setProvider(p);
              }}
            />
          </FieldSet>

          <FieldSet>
            <FormLabel>Key Name</FormLabel>
            <Input
              type="text"
              name="keyName"
              value={keyName}
              required
              placeholder="Assign a unique name to this provider key for easy identification."
              onChange={handleInputChange}
            />
          </FieldSet>

          {provider &&
            provider.configurations?.map((x, idx) => (
              <FieldSet key={idx}>
                <FormLabel htmlFor={`config.${x.name}`}>{x.label}</FormLabel>
                {x.type === 'text' ? (
                  <Textarea
                    required
                    name={`config.${x.name}`}
                    placeholder={x.label}
                    value={config[x.name] || ''}
                    onChange={handleInputChange}
                  />
                ) : (
                  <Input
                    type="text"
                    required
                    name={`config.${x.name}`}
                    placeholder={x.label}
                    value={config[x.name] || ''}
                    onChange={handleInputChange}
                  />
                )}
              </FieldSet>
            ))}

          <ErrorMessage message={error} />
        </ModalBody>
        <ModalFooter>
          <ICancelButton
            className="px-4 rounded-[2px]"
            onClick={() => {
              props.setModalOpen(false);
            }}
          >
            Cancel
          </ICancelButton>
          <IBlueBGArrowButton
            className="px-4 rounded-[2px]"
            onClick={validateAndSubmit}
            isLoading={loading}
          >
            Configure
          </IBlueBGArrowButton>
        </ModalFooter>
      </ModalFitHeightBlock>
    </GenericModal>
  );
}
