import { Divider, ListItem, OrderedList } from "@chakra-ui/react";

import { PostProps } from "src/api/openapi/schemas";

import { PostListItem } from "./PostListItem";

type Props = {
  posts: PostProps[];
};

export function PostList(props: Props) {
  return (
    <OrderedList gap={4} display="flex" flexDir="column" width="full" m={0}>
      {props.posts.map((p) => (
        <>
          <Divider />
          <ListItem key={p.id} listStyleType="none">
            <PostListItem {...p} />
          </ListItem>
        </>
      ))}
    </OrderedList>
  );
}
