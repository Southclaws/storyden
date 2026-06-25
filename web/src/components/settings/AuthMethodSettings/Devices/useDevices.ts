import {
  useAccountAuthProviderList,
  useAccountGet,
} from "@/api/openapi-client/accounts";
import { AccountAuthMethod } from "@/api/openapi-schema";
import { passkeyRegister } from "@/components/auth/webauthn/utils";
import { deriveError } from "@/utils/error";

export type Props = {
  active: AccountAuthMethod[];
};

export function useDevices() {
  const { mutate } = useAccountAuthProviderList();
  const { data, error } = useAccountGet();
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
