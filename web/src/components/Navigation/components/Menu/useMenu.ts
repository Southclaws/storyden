import { useCategoryList } from "src/api/openapi/categories";
import { useAuthProvider } from "src/auth/useAuthProvider";

export function useMenu() {
  const { account } = useAuthProvider();

  const { data, error } = useCategoryList();

  return {
    isAuthenticated: !!account,
    data,
    error,
  };
}
