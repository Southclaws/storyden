import { last } from "lodash";

export function getTargetSlug(slug: string[]): string {
  return last(slug) ?? "";
}
