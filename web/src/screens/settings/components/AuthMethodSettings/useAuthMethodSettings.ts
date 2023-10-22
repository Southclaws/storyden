import { partition } from "lodash/fp";
import { SWRResponse } from "swr";

import { useAccountAuthProviderList } from "src/api/openapi/accounts";
import {
  APIError,
  AccountAuthMethod,
  AccountAuthMethods,
} from "src/api/openapi/schemas";

const partitionByActive = partition<AccountAuthMethod>("active");

type AuthMethodSettings =
  | {
      ready: false;
      error?: APIError;
    }
  | {
      ready: true;
      active: AccountAuthMethod[];
      rest: AccountAuthMethod[];
    };

export function useAuthMethodSettings(): AuthMethodSettings {
  const response: SWRResponse<AccountAuthMethods, APIError> =
    useAccountAuthProviderList();

  if (!response.data) {
    return { ready: false, error: response.error };
  }

  const [active, rest] = partitionByActive(
    response.data.auth_methods.sort((a, b) => a.name.localeCompare(b.name)),
  );

  return {
    ready: true,
    active,
    rest,
  };
}
