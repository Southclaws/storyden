import { Fragment } from "react";

import { PostProps } from "src/api/openapi/schemas";

import { Divider, styled } from "@/styled-system/jsx";

import { PostView } from "./PostView/PostView";

type Props = {
  slug?: string;
  posts: PostProps[];
};

export function PostListView(props: Props) {
  return (
    <styled.ol
      listStyleType="none"
      m="0"
      gap="4"
      display="flex"
      flexDir="column"
      width="full"
    >
      {props.posts.map((p) => (
        <Fragment key={p.id}>
          <Divider />
          <styled.li key={p.id} listStyleType="none" m="0">
            <PostView slug={props.slug} {...p} />
          </styled.li>
        </Fragment>
      ))}
    </styled.ol>
  );
}
