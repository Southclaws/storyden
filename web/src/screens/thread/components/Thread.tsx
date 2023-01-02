import { Box, Heading, ListItem, OrderedList, VStack } from "@chakra-ui/react";
import { Post, Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";
import { PostView } from "./Post";

export function PostListView(props: { posts: Post[] }) {
  return (
    <OrderedList gap={2} display="flex" flexDir="column" width="full">
      {props.posts.map((p) => (
        <ListItem key={p.id} listStyleType="none" m={0}>
          <PostView {...p} />
        </ListItem>
      ))}
    </OrderedList>
  );
}

export function ThreadView(props: Thread) {
  return (
    <Box as="main">
      <VStack alignItems="start" px={3}>
        <Heading>{props.title}</Heading>
        <CategoryPill category={props.category} />

        <PostListView {...props} />
      </VStack>
    </Box>
  );
}
