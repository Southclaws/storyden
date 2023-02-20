import { useCategoriesList } from "src/api/openapi/categories";
import { useAuthProvider } from "src/auth/useAuthProvider";

export function useMenu() {
  const { account } = useAuthProvider();

  const { data, error } = useCategoriesList();

  return {
    isAuthenticated: !!account,
    data,
    error,
  };
}
