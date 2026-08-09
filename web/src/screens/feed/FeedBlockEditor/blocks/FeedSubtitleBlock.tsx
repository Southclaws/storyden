import { Input } from "@/components/ui/input";
import { Text } from "@/components/ui/text";

import { useFeedBlockEditor } from "../Context";
import { useDebouncedSiteUpdate } from "../useDebouncedSiteUpdate";

export function FeedSubtitleBlock() {
  const { initialSettings, isEditing } = useFeedBlockEditor();
  const update = useDebouncedSiteUpdate();
  const description = initialSettings?.description ?? "";

  if (!isEditing) {
    return description ? <Text variant="supporting">{description}</Text> : null;
  }

  return (
    <Input
      aria-label="Site subtitle"
      defaultValue={description}
      placeholder="Site subtitle"
      variant="ghost"
      onChange={(event) => update({ description: event.currentTarget.value })}
    />
  );
}
