import { useQueryState } from "nuqs";

export function useSelectedPlugin() {
  return useQueryState("plugin");
}
