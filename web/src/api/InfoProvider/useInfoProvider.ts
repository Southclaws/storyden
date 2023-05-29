import { createContext, useContext } from "react";
import { Info } from "../openapi/schemas";

export const InfoContext = createContext<Info | undefined>(undefined);

export function useInfoProvider() {
  return useContext(InfoContext);
}
