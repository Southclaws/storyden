import { useAccountGet } from "src/api/openapi/accounts";
import { AccountAuthMethod } from "src/api/openapi/schemas";
import { passkeyRegister } from "src/components/auth/webauthn/utils";
import { deriveError } from "src/utils/error";

export type Props = {
  active: AccountAuthMethod[];
};

export function useDevices() {
  const { data, error } = useAccountGet();
  if (!data) return { ready: false as const, error };

  const { handle } = data;

  async function handleDeviceRegister() {
    try {
      await passkeyRegister(handle);
    } catch (e) {
      deriveError(e);
    }
  }

  return {
    ready: true as const,
    handleDeviceRegister,
  };
}
