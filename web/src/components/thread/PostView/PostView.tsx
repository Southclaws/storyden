import { Button } from "@chakra-ui/react";

import { PostProps } from "src/api/openapi/schemas";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";

import { Byline } from "../../content/Byline";
import { PostMenu } from "../PostMenu/PostMenu";
import { ReactList } from "../ReactList/ReactList";

import { Flex, HStack } from "@/styled-system/jsx";

import { usePostView } from "./usePostView";

type Props = PostProps & {
  slug?: string;
};

export function PostView(props: Props) {
  const {
    isEditing,
    editingContent,
    onContentChange,
    onPublishEdit,
    onCancelEdit,
  } = usePostView(props);

  return (
    <Flex id={props.id} flexDir="column" gap="2">
      <Byline
        href={`#${props.id}`}
        author={props.author}
        time={new Date(props.createdAt)}
        updated={new Date(props.updatedAt)}
        more={<PostMenu {...props} />}
      />
      {isEditing ? (
        <>
          <ContentComposer
            onChange={onContentChange}
            initialValue={editingContent}
          />
          <HStack>
            <Button onClick={onPublishEdit}>Update</Button>
            <Button variant="outline" onClick={onCancelEdit}>
              Cancel
            </Button>
          </HStack>
        </>
      ) : (
        <>
          <ContentViewer value={props.body} />
        </>
      )}
      <ReactList {...props} />
    </Flex>
  );
}
