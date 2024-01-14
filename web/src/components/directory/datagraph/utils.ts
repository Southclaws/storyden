import { last, nth, takeWhile } from "lodash";

export function getTargetSlug(slug: string[]): [string, string, boolean] {
  const top = last(slug);

  const isNew = top === "new";

  const target = isNew
    ? // If the tail item is "new" then walk back until we find an actual slug.
      last(takeWhile(slug, (s) => s !== "new"))
    : // Otherwise, it's whatever the last item is.
      top;

  // The fallback is for when the target is actually an item. Items cannot have
  // children, so we can assume the parent is the left-most slug of the target,
  // which is the third index from the end of the slug array. This is an edge
  // case though as the only way to hit this would be to manually add "/new".
  const fallback = slug.length > 2 ? nth(slug, -3) : "";

  return [target ?? "", fallback ?? "", isNew];
}
