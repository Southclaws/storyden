import { PostProps, ThreadReference } from "src/api/openapi/schemas";
import { EmptyState } from "src/components/feed/EmptyState";

import { styled } from "@/styled-system/jsx";

import { PostRef } from "./PostRef";

type Either = PostProps | ThreadReference;

type Props = {
  items: Either[];
};

export function PostRefList({ items }: Props) {
  if (items.length === 0) {
    return <EmptyState />;
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
