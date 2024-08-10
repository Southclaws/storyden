import { SWRResponse } from "swr";

import { useAccountAuthProviderList } from "src/api/openapi-client/accounts";
import { APIError, AccountAuthMethods } from "src/api/openapi-schema";
import {
  groupAuthMethods,
  groupAuthProviders,
} from "src/components/settings/utils";

export function useAuthMethodSettings() {
  const response: SWRResponse<AccountAuthMethods, APIError> =
    useAccountAuthProviderList();

  if (!response.data) {
    return { ready: false as const, error: response.error };
  }

  const { active, available } = response.data;

  const { password, phone, webauthn, providers } =
    groupAuthProviders(available);

  const {
    password: passwordActive,
    phone: phoneActive,
    webauthn: webauthnActive,
    methods,
  } = groupAuthMethods(active);

  const sorted = providers.sort((a, b) => a.name.localeCompare(b.name));

  return {
    ready: true as const,
    available: {
      password,
      phone,
      webauthn,
      oauth: sorted,
    },
    active: {
      password: passwordActive ?? [],
      phone: phoneActive ?? [],
      webauthn: webauthnActive ?? [],
      methods,
    },
  };
}
