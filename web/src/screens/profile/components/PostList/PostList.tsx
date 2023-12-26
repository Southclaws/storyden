import { PostProps } from "src/api/openapi/schemas";

import { styled } from "@/styled-system/jsx";

import { PostListItem } from "./PostListItem";

type Props = {
  posts: PostProps[];
};

export function PostList(props: Props) {
  return (
    <styled.ol gap="4" display="flex" flexDir="column" width="full" m="0">
      {props.posts.map((p) => (
        <styled.li key={p.id} listStyleType="none">
          <PostListItem {...p} />
        </styled.li>
      ))}
    </styled.ol>
  );
}
