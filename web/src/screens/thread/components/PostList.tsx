import { OrderedList, ListItem, Divider } from "@chakra-ui/react";
import { Post } from "src/api/openapi/schemas";
import { PostView } from "./Post";

type Props = {
  slug: string;
  posts: Post[];
};

export function PostListView(props: Props) {
  return (
    <OrderedList gap={4} display="flex" flexDir="column" width="full">
      {props.posts.map((p) => (
        <>
          <Divider />
          <ListItem key={p.id} listStyleType="none" m={0}>
            <PostView slug={props.slug} {...p} />
          </ListItem>
        </>
      ))}
    </OrderedList>
  );
}
