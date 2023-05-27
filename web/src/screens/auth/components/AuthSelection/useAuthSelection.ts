import { useAuthProviderList } from "src/api/openapi/auth";

export function useAuthSelection() {
  const { data, error } = useAuthProviderList();

  return { data, error };
}
