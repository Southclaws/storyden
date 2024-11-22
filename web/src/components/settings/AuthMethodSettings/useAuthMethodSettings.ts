import { find } from "lodash";

import { useAccountAuthProviderList } from "src/api/openapi-client/accounts";

import { groupAuthMethods, groupAuthProviders } from "@/lib/auth/utils";

export function useAuthMethodSettings() {
  const { data, error } = useAccountAuthProviderList();
  if (!data) {
    return {
      ready: false as const,
      error: error,
    };
  }

  const { active, available } = data;

  const { password, phone, webauthn, oauth } = groupAuthProviders(available);

  const {
    password: passwordActive,
    phone: phoneActive,
    webauthn: webauthnActive,
    methods,
  } = groupAuthMethods(active);

  // Remove any OAuth providers that are already active
  const availableOAuth = oauth.filter((v) => {
    return !find(methods, (m) => m.provider.provider === v.provider);
  });

  const sorted = availableOAuth.sort((a, b) => a.name.localeCompare(b.name));

  return {
    ready: true as const,
    data: {
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
    },
  };
}
