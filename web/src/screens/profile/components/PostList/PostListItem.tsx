import { Flex } from "@chakra-ui/react";

import { PostProps } from "src/api/openapi/schemas";
import { Byline } from "src/components/content/Byline";
import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";
import { PostMenu } from "src/components/thread/PostMenu/PostMenu";
import { ReactList } from "src/components/thread/ReactList/ReactList";

export function PostListItem(props: PostProps) {
  return (
    <Flex id={props.id} flexDir="column" gap={2}>
      <ContentViewer value={props.body} />
      <ReactList {...props} />

      <Byline
        href={`#${props.id}`}
        author={props.author}
        time={new Date(props.createdAt)}
        updated={new Date(props.updatedAt)}
        more={<PostMenu {...props} />}
      />
    </Flex>
  );
}
