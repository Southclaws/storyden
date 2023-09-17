import { ThreadReference } from "src/api/openapi/schemas";
import { EmptyState } from "src/components/feed/EmptyState";

import { styled } from "@/styled-system/jsx";

import { TextPost } from "./TextPost";

type Props = {
  posts: ThreadReference[];
  showEmptyState?: boolean | undefined;
};

export function TextPostList(props: Props) {
  if (props.showEmptyState && props.posts.length === 0) {
    return <EmptyState />;
  }

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap={2}>
      {props.posts.map((t) => (
        <TextPost key={t.id} thread={t} />
      ))}
    </styled.ol>
  );
}
