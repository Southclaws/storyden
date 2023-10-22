import { Divider, ListItem, OrderedList } from "@chakra-ui/react";
import { Fragment } from "react";

import { PostProps } from "src/api/openapi/schemas";

import { PostView } from "./PostView/PostView";

type Props = {
  slug?: string;
  posts: PostProps[];
};

export function PostListView(props: Props) {
  return (
    <OrderedList
      listStyleType="none"
      m={0}
      gap={4}
      display="flex"
      flexDir="column"
      width="full"
    >
      {props.posts.map((p) => (
        <Fragment key={p.id}>
          <Divider />
          <ListItem key={p.id} listStyleType="none" m={0}>
            <PostView slug={props.slug} {...p} />
          </ListItem>
        </Fragment>
      ))}
    </OrderedList>
  );
}
