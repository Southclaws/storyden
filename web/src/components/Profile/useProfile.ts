import { Account } from "src/api/openapi/schemas";

type UseProfileData =
  | {
      authenticated: false;
    }
  | {
      authenticated: true;
      account: Account;
    };

export function useProfile(): UseProfileData {
  // TODO: useSession and get account.

  return {
    authenticated: false,
  };
}
