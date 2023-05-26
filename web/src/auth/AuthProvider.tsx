import { createContext, PropsWithChildren } from "react";
import { Account } from "src/api/openapi/schemas";
import { useAuthProvider } from "./useAuthProvider";

export const AuthContext = createContext<Account | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const { account } = useAuthProvider();
  return (
    <AuthContext.Provider value={account}>{children}</AuthContext.Provider>
  );
}
