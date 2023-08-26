import { useAccountGet } from "src/api/openapi/accounts";

export function useSession() {
  const { data } = useAccountGet();
  return data;
}
