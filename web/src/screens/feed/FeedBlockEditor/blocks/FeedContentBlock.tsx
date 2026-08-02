import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";

import { useFeedBlockEditor } from "../Context";
import { useDebouncedSiteUpdate } from "../useDebouncedSiteUpdate";

export function FeedContentBlock() {
  const { initialSettings, isEditing } = useFeedBlockEditor();
  const update = useDebouncedSiteUpdate();

  return (
    <ContentComposer
      disabled={!isEditing}
      initialValue={initialSettings?.content ?? ""}
      placeholder="About this community..."
      onChange={(content) => update({ content })}
    />
  );
}
