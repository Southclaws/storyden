import {
  TagNameList,
  TagReference,
  TagReferenceList,
} from "@/api/openapi-schema";
import { HStack } from "@/styled-system/jsx";

import { TagBadge } from "./TagBadge";

export type Props = InteractionProps & {
  tags: TagReferenceList;
  showItemCount?: boolean;
};

type InteractionProps =
  | {
      type?: "link";
      onClick?: never;
      highlightedTags?: never;
    }
  | {
      type: "button";
      onClick: (tr: TagReference) => Promise<void>;
      highlightedTags?: TagNameList;
    };

export function TagBadgeList({
  type,
  onClick,
  tags,
  showItemCount,
  highlightedTags,
}: Props) {
  return (
    <HStack flexWrap="wrap">
      {tags.map((r) => (
        <TagBadge
          key={r.name}
          tag={r}
          showItemCount={showItemCount}
          {...(type === "button"
            ? {
                type: "button",
                onClick: () => onClick?.(r),
                highlighted: highlightedTags?.includes(r.name),
              }
            : {
                type: "link",
              })}
        />
      ))}
    </HStack>
  );
}
