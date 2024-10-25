import { TagReferenceList } from "@/api/openapi-schema";
import { HStack } from "@/styled-system/jsx";

import { TagBadge } from "./TagBadge";

export type Props = {
  tags: TagReferenceList;
  showItemCount?: boolean;
};

export function TagBadgeList({ tags, showItemCount }: Props) {
  return (
    <HStack flexWrap="wrap">
      {tags.map((r) => (
        <TagBadge key={r.name} tag={r} showItemCount={showItemCount} />
      ))}
    </HStack>
  );
}
