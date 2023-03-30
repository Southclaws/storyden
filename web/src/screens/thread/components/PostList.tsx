import { OrderedList, ListItem, Divider } from "@chakra-ui/react";
import { Post } from "src/api/openapi/schemas";
import { PostView } from "./Post";

export function PostListView(props: { posts: Post[] }) {
  return (
    <OrderedList gap={4} display="flex" flexDir="column" width="full">
      {props.posts.map((p) => (
        <>
          <Divider />
          <ListItem key={p.id} listStyleType="none" m={0}>
            <PostView {...p} />
          </ListItem>
        </>
      ))}
    </OrderedList>
  );
}
