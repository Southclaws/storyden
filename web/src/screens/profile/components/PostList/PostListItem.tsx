import { PostProps } from "src/api/openapi/schemas";
import { Byline } from "src/components/content/Byline";
import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";
import { PostMenu } from "src/components/thread/PostMenu/PostMenu";

import { Flex } from "@/styled-system/jsx";

export function PostListItem(props: PostProps) {
  return (
    <Flex id={props.id} flexDir="column" gap="2">
      <ContentViewer value={props.body} />

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
