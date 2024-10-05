import { PostReference } from "src/api/openapi-schema";

import { Empty } from "@/components/site/Empty";
import { styled } from "@/styled-system/jsx";

import { PostRef } from "./PostRef";

type Props = {
  items: PostReference[];
  emptyText?: string;
};

export function PostRefList({ items, emptyText }: Props) {
  if (items.length === 0) {
    return <Empty>{emptyText}</Empty>;
  }

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="4">
      {items.map((t) => (
        <PostRef key={t.id} item={t} />
      ))}
    </styled.ol>
  );
}
