import { ListItem, OrderedList } from "@chakra-ui/react";
import { Thread } from "src/api/openapi/schemas";
import { Post } from "./Post";

export function Thread(props: Thread) {
  return (
    <OrderedList>
      {props.posts.map((p) => (
        <ListItem key={p.id}>
          <Post {...p} />
        </ListItem>
      ))}
    </OrderedList>
  );
}
