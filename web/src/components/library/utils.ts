import { last, takeWhile } from "lodash";

export function getTargetSlug(slug: string[]): [string, boolean] {
  const top = last(slug);

  const isNew = top === "new";

  const target = isNew
    ? // If the tail item is "new" then walk back until we find an actual slug.
      last(takeWhile(slug, (s) => s !== "new"))
    : // Otherwise, it's whatever the last item is.
      top;

  return [target ?? "", isNew];
}
