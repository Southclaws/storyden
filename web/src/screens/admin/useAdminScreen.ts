import { useAccountGet } from "src/api/openapi/accounts";

export function useAdminScreen() {
  const { data, error } = useAccountGet();

  return {
    data,
    error,
  };
}
