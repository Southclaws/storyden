import { useAuthProviderList } from "src/api/openapi/auth";

export function useAuthSelection() {
  const authProviderList = useAuthProviderList();

  return authProviderList;
}
