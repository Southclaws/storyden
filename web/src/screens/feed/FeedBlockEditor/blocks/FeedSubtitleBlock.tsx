import { Text } from "@/components/ui/text";

import { useFeedBlockEditor } from "../Context";

export function FeedSubtitleBlock() {
  const { initialSettings } = useFeedBlockEditor();
  const description = initialSettings?.description ?? "";

  return description ? <Text variant="supporting">{description}</Text> : null;
}
