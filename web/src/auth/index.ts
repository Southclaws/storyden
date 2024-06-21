import { useAccountGet } from "src/api/openapi/accounts";

import { Account } from "@/api/openapi/schemas";

export function useSession(initial?: Account) {
  const { data } = useAccountGet({
    swr: {
      fallbackData: initial,
    },
  });
  return data;
}
