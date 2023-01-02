import { OrderedList, ListItem } from "@chakra-ui/react";
import { Post } from "src/api/openapi/schemas";
import { PostView } from "./Post";

export function PostListView(props: { posts: Post[] }) {
  return (
    <OrderedList gap={6} display="flex" flexDir="column" width="full">
      {props.posts.map((p) => (
        <ListItem key={p.id} listStyleType="none" m={0}>
          <PostView {...p} />
        </ListItem>
      ))}
    </OrderedList>
  );
}
