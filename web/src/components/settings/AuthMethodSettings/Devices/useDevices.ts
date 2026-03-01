import {
  useAccountAuthProviderList,
} from "src/api/openapi-client/accounts";
import { useAccountSession } from "src/auth";
import { AccountAuthMethod } from "src/api/openapi-schema";
import { passkeyRegister } from "src/components/auth/webauthn/utils";
import { deriveError } from "src/utils/error";

export type Props = {
  active: AccountAuthMethod[];
};

export function useDevices() {
  const { mutate } = useAccountAuthProviderList();
  const { data, error } = useAccountSession();
  if (!data) return { ready: false as const, error };

  const { handle } = data;

  async function handleDeviceRegister() {
    try {
      await passkeyRegister(handle);
      mutate();
    } catch (e) {
      deriveError(e);
    }
  }

  return {
    ready: true as const,
    handleDeviceRegister,
  };
}
