import { indexOf } from "lodash";
import { uniq } from "lodash/fp";
import { z } from "zod";

export const ParamsSchema = z.object({
  slug: z
    .string()
    .array()
    .min(1)
    .transform((slugs) => slugs.map((s) => decodeURIComponent(s))),
});
export type Params = z.infer<typeof ParamsSchema>;

export type LibraryPath = string[];

/**
 * Paths within the /l route group are very flexible. This is mostly to simplify
 * the backend so we don't need to walk up the entire tree to find all the
 * parents of a given node or item.
 *
 * What this means is that there's actually zero validation against the path. If
 * you visit `/l/does/not/exist/actually-does-exist` then as long as the
 * `actually-does-exist` is a valid node or item slug, it will render fine.
 *
 * But, we basically generate a tree-like path to give the impression that there
 * is a hierarchy. This function is the main helper for that. Basically, when
 * you're turning a `slug` into a href, you must know the current URL path in
 * order to create the new path. This function takes the current path and the
 * slug of the target item and returns the new path.
 *
 * The new path will simply be the slug appended to the end unless the slug is
 * already in the current path. If it is, it slices the path at that point.
 *
 * For example, say you're at: `/l/node-1/node-2/item-1` and the page lists
 * `node-2` as a parent of `item-1`, this function will be given:
 * `["node-1", "node-2", "item-1"]` and return `node-1/node-2` since
 * `node-2` is passed in as `end` and is already in the path.
 *
 * @param onto the libraryPath, basically a list of slugs from `[...slug]`.
 * @param end the slug of the target library node.
 * @returns a string of the path to append to `/l/`.
 */
export function joinLibraryPath(onto: LibraryPath, end: string): string {
  const inPath = indexOf(onto, end);

  const list = inPath === -1 ? [...onto, end] : onto.slice(0, inPath + 1);

  return uniq(list).join("/");
}

/**
 * replaceLibraryPath is for when the slug of a library node changes. It's
 * similar to joinLibraryPath but it'll replace the old slug with the new one.
 */
export function replaceLibraryPath(
  onto: LibraryPath,
  oldSlug: string,
  newSlug: string,
): string {
  const inPath = indexOf(onto, oldSlug);

  const list =
    inPath === -1 ? [...onto, newSlug] : [...onto.slice(0, inPath), newSlug];

  return list.join("/");
}
