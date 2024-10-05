import {
  accountAuthMethodDelete,
  useAccountAuthProviderList,
} from "src/api/openapi-client/accounts";
import { APIError } from "src/api/openapi-schema";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";

export type Props = {
  id: string;
};

export type WithDisclosure<T> = UseDisclosureProps & T;

export function useDeleteDeviceScreen(props: WithDisclosure<Props>) {
  const { mutate } = useAccountAuthProviderList();

  const handleConfirm = async () => {
    handle(async () => {
      await accountAuthMethodDelete(props.id);

      mutate();

      props.onClose?.();
    });
  };

  return {
    handleConfirm,
  };
}
