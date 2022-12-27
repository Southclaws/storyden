import { Flex, Text } from "@chakra-ui/react";
import { Post } from "src/api/openapi/schemas";
import { Byline } from "./Byline";

export function Post(props: Post) {
  return (
    <Flex flexDir="column">
      <Text>{props.body}</Text>
      <Byline author={props.author.handle} time={new Date(props.createdAt)} />
    </Flex>
  );
}
