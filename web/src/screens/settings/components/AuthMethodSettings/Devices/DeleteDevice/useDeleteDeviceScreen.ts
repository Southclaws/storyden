import { UseDisclosureProps, useToast } from "@chakra-ui/react";

import {
  accountAuthMethodDelete,
  useAccountAuthProviderList,
} from "src/api/openapi/accounts";
import { APIError } from "src/api/openapi/schemas";
import { errorToast } from "src/components/site/ErrorBanner";

export type Props = {
  id: string;
};

export type WithDisclosure<T> = UseDisclosureProps & T;

export function useDeleteDeviceScreen(props: WithDisclosure<Props>) {
  const toast = useToast();
  const { mutate } = useAccountAuthProviderList();

  const handleConfirm = async () => {
    try {
      await accountAuthMethodDelete(props.id);

      mutate();

      props.onClose?.();
    } catch (e: unknown) {
      errorToast(toast)(e as APIError);
    }
  };

  return {
    handleConfirm,
  };
}
