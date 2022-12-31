import { Flex } from "@chakra-ui/react";
import { ReactMarkdown } from "react-markdown/lib/react-markdown";
import { Post } from "src/api/openapi/schemas";
import { Byline } from "./Byline";

export function PostView(props: Post) {
  return (
    <Flex flexDir="column">
      <ReactMarkdown>{props.body}</ReactMarkdown>
      <Byline author={props.author.handle} time={new Date(props.createdAt)} />
    </Flex>
  );
}
