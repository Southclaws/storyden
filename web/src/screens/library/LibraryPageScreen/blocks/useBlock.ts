import { LibraryPageBlock, LibraryPageBlockType } from "@/lib/library/metadata";

import { useWatch } from "../store";

export function useBlock<T extends LibraryPageBlockType>(
  type: T,
): Extract<LibraryPageBlock, { type: T }> | undefined {
  const block = useWatch((s) => {
    const layout = s.draft.meta.layout;
    const block = layout?.blocks.find(
      (b): b is Extract<LibraryPageBlock, { type: T }> => b.type === type,
    );
    return block;
  });

  return block;
}
