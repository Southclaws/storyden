import { PageHeading } from "@/components/ui/page-heading";

import { useFeedBlockEditor } from "../Context";

export function FeedTitleBlock() {
  const { initialSettings } = useFeedBlockEditor();
  const title = initialSettings?.title ?? "Storyden";

  return <PageHeading>{title}</PageHeading>;
}
