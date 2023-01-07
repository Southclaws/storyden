import { Account } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

type UseProfileData =
  | {
      authenticated: false;
    }
  | {
      authenticated: true;
      account: Account;
    };

export function useProfile(): UseProfileData {
  const account = useSession();
  if (!account) return { authenticated: false };

  return {
    authenticated: true,
    account,
  };
}
