import { Heading, ListItem, OrderedList, VStack } from "@chakra-ui/react";
import { Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";
import { Post } from "./Post";

export function Thread(props: Thread) {
  return (
    <VStack alignItems="start" px={3}>
      <Heading>{props.title}</Heading>
      <CategoryPill category={props.category} />

      <OrderedList gap={2} display="flex" flexDir="column">
        {props.posts.map((p) => (
          <ListItem key={p.id} listStyleType="none" m={0}>
            <Post {...p} />
          </ListItem>
        ))}
      </OrderedList>
    </VStack>
  );
}
