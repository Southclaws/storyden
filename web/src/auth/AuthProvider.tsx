import { createContext, PropsWithChildren } from "react";
import LoadingBanner from "src/components/LoadingBanner";
import { Account } from "src/api/openapi/schemas";
import { useAuthProvider } from "./useAuthProvider";

export const AuthContext = createContext<Account | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const { firstTime, account } = useAuthProvider();
  return (
    <AuthContext.Provider value={account}>
      {firstTime ? <LoadingBanner /> : children}
    </AuthContext.Provider>
  );
}
