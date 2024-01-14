import { PostProps, ThreadReference } from "src/api/openapi/schemas";

import { Empty } from "../../../site/Empty";

import { styled } from "@/styled-system/jsx";

import { PostRef } from "./PostRef";

type Either = PostProps | ThreadReference;

type Props = {
  items: Either[];
  emptyText?: string;
};

export function PostRefList({ items, emptyText }: Props) {
  if (items.length === 0) {
    return <Empty>{emptyText}</Empty>;
  }

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="4">
      {items.map((t) =>
        isThread(t) ? (
          <PostRef key={t.id} kind="thread" item={t as ThreadReference} />
        ) : (
          <PostRef key={t.id} kind="post" item={t as PostProps} />
        ),
      )}
    </styled.ol>
  );
}

function isThread(e: Either) {
  return "title" in (e as ThreadReference);
}
