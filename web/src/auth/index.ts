import { useAccountGet } from "@/api/openapi-client/accounts";
import { Account } from "@/api/openapi-schema";

export function useSession(initial?: Account) {
  const { data } = useAccountGet({
    swr: {
      fallbackData: initial,
    },
  });
  return data;
}
