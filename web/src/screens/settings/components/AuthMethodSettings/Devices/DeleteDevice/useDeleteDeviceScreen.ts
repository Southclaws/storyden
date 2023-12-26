import {
  accountAuthMethodDelete,
  useAccountAuthProviderList,
} from "src/api/openapi/accounts";
import { APIError } from "src/api/openapi/schemas";
import { handleError } from "src/components/site/ErrorBanner";
import { UseDisclosureProps } from "src/utils/useDisclosure";

export type Props = {
  id: string;
};

export type WithDisclosure<T> = UseDisclosureProps & T;

export function useDeleteDeviceScreen(props: WithDisclosure<Props>) {
  const { mutate } = useAccountAuthProviderList();

  const handleConfirm = async () => {
    try {
      await accountAuthMethodDelete(props.id);

      mutate();

      props.onClose?.();
    } catch (e: unknown) {
      handleError(e as APIError);
    }
  };

  return {
    handleConfirm,
  };
}
