"use client";

import { useTagList } from "@/api/openapi-client/tags";
import { TagListResult } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Breadcrumbs } from "@/components/ui/Breadcrumbs";
import { Text } from "@/components/ui/text";
import { LStack } from "@/styled-system/jsx";

type Props = {
  initialTagList: TagListResult;
};

export function TagsIndexScreen(props: Props) {
  const { data, error } = useTagList(
    {},
    { swr: { fallbackData: props.initialTagList } },
  );
  if (!data) {
    return <Unready error={error} />;
  }

  const tags = data.tags
    .filter((t) => t.item_count > 0)
    .sort((a, b) => b.item_count - a.item_count);

  return (
    <LStack>
      <LStack gap="1">
        <Breadcrumbs
          index={{
            label: "Tags",
            href: "/tags",
          }}
          crumbs={[]}
        />

        <Text textStyle="sm">
          Threads and library pages can be tagged with related topics.
        </Text>
      </LStack>

      <TagBadgeList tags={tags} showItemCount />
    </LStack>
  );
}
