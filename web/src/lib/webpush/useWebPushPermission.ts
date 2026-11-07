import { useAccountGet } from "@/api/openapi-client/accounts";
import { useGetInfo } from "@/api/openapi-client/misc";
import { getAPIAddress } from "@/config";

import { enableWebPush } from "./client";

export function useWebPushPermission() {
  const { data: account } = useAccountGet();
  const { data: info } = useGetInfo();

  return async function requestPushPermission() {
    const publicKey = info?.web_push_public_key;

    if (!account || !publicKey) {
      return;
    }

    try {
      await enableWebPush({
        accountId: account.id,
        apiAddress: getAPIAddress(),
        publicKey,
      });
    } catch (error) {
      console.warn("[Web Push] Enabling notifications failed", error);
    }
  };
}
