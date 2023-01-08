import { Account } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

type UseOptions =
  | {
      authenticated: false;
    }
  | {
      authenticated: true;
      account: Account;
    };

export function useOptions(): UseOptions {
  const account = useSession();
  if (!account) return { authenticated: false };

  return {
    authenticated: true,
    account,
  };
}
