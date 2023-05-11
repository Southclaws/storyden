import { Divider, ListItem, OrderedList } from "@chakra-ui/react";
import { PostProps } from "src/api/openapi/schemas";
import { PostView } from "./PostView/PostView";

type Props = {
  slug: string;
  posts: PostProps[];
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
