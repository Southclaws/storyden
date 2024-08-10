import { PostReference } from "src/api/openapi/schemas";

import { styled } from "@/styled-system/jsx";

import { Empty } from "../../../site/Empty";

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
