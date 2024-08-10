import { Post } from "src/api/openapi/schemas";
import { Byline } from "src/components/content/Byline";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { PostMenu } from "src/components/thread/PostMenu/PostMenu";

import { Flex } from "@/styled-system/jsx";

export function PostListItem(props: Post) {
  return (
    <Flex id={props.id} flexDir="column" gap="2">
      <ContentComposer disabled initialValue={props.body} />

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
