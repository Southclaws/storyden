"use client";

import { parseAsInteger, useQueryStates } from "nuqs";

import { useTagList } from "@/api/openapi-client/tags";
import { TagListResult } from "@/api/openapi-schema";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Breadcrumbs } from "@/components/ui/Breadcrumbs";
import { Text } from "@/components/ui/text";
import { LStack } from "@/styled-system/jsx";

type Props = {
  initialTagList: TagListResult;
};

export function TagsIndexScreen(props: Props) {
  const [filters] = useQueryStates({
    page: parseAsInteger.withDefault(1),
  });

  const { data, error } = useTagList(
    { page: filters.page.toString() },
    { swr: { fallbackData: props.initialTagList } },
  );
  if (!data) {
    return <Unready error={error} />;
  }

  const tags = data.tags.filter((t) => t.item_count > 0);

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

      <PaginationControls
        path="/tags"
        currentPage={filters.page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />
    </LStack>
  );
}
