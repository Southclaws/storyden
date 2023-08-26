import { PropsWithChildren } from "react";

import { useAuthProvider } from "./useAuthProvider";

export function AuthProvider({ children }: PropsWithChildren) {
  useAuthProvider();
  return <>{children}</>;
}
