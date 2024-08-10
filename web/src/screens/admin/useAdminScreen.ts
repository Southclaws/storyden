import { useAccountGet } from "src/api/openapi-client/accounts";

export function useAdminScreen() {
  const { data, error } = useAccountGet();

  return {
    data,
    error,
  };
}
