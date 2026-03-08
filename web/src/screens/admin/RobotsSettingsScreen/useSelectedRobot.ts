import { useQueryState } from "nuqs";

export function useSelectedRobot() {
  return useQueryState("robot");
}
