import { ThreadReference } from "src/api/openapi/schemas";
import { EmptyState } from "src/components/feed/EmptyState";

import { styled } from "@/styled-system/jsx";

import { TextPost } from "./TextPost";

type Props = {
  posts: ThreadReference[];
};

export function TextPostList(props: Props) {
  if (props.posts.length === 0) {
    return <EmptyState />;
  }

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap={4}>
      {props.posts.map((t) => (
        <TextPost key={t.id} thread={t} />
      ))}
    </styled.ol>
  );
}
