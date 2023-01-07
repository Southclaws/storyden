import { useContext } from "react";
import { AuthContext } from "./AuthProvider";

export function useSession() {
  return useContext(AuthContext);
}
