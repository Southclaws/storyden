"use client";

import { useEffect } from "react";

import { RequestError } from "@/api/common";
import { useAccountGet } from "@/api/openapi-client/accounts";
import { useGetInfo } from "@/api/openapi-client/misc";
import { getAPIAddress } from "@/config";

import { clearWebPush, reconcileWebPush, supportsWebPush } from "./client";

export function WebPush() {
  const { data: account, error: accountError } = useAccountGet();
  const { data: info } = useGetInfo();
  const accountId = account?.id;
  const publicKey = info?.web_push_public_key;
  const infoReady = Boolean(info);
  const signedOut =
    accountError instanceof RequestError && accountError.status === 401;

  useEffect(() => {
    if (!supportsWebPush()) {
      return;
    }

    if (!signedOut && (!accountId || !infoReady)) {
      return;
    }

    async function reconcile() {
      try {
        if (signedOut || !publicKey) {
          await clearWebPush();
          return;
        }

        if (!accountId) {
          return;
        }

        await reconcileWebPush({
          accountId,
          apiAddress: getAPIAddress(),
          publicKey,
        });
      } catch (error) {
        console.warn("Web Push setup failed", error);
      }
    }

    void reconcile();
    window.addEventListener("online", reconcile);

    return () => window.removeEventListener("online", reconcile);
  }, [accountId, publicKey, infoReady, signedOut]);

  return null;
}
