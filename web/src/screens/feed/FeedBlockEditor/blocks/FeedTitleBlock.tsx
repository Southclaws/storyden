import { HeadingInput } from "@/components/ui/heading-input";
import { PageHeading } from "@/components/ui/page-heading";

import { useFeedBlockEditor } from "../Context";
import { useDebouncedSiteUpdate } from "../useDebouncedSiteUpdate";

export function FeedTitleBlock() {
  const { initialSettings, isEditing } = useFeedBlockEditor();
  const update = useDebouncedSiteUpdate();
  const title = initialSettings?.title ?? "Storyden";

  if (!isEditing) {
    return <PageHeading>{title}</PageHeading>;
  }

  return (
    <HeadingInput
      aria-label="Site title"
      defaultValue={title}
      placeholder="Site title"
      onValueChange={(value) => update({ title: value })}
    />
  );
}
