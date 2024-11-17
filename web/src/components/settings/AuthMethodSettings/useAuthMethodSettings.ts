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

  const sorted = oauth.sort((a, b) => a.name.localeCompare(b.name));

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
