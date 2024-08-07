import { PostProps } from "src/api/openapi/schemas";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { Byline } from "../../content/Byline";
import { PostMenu } from "../PostMenu/PostMenu";

import { Button } from "@/components/ui/button";
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
            <Button variant="ghost" onClick={onCancelEdit}>
              Cancel
            </Button>
          </HStack>
        </>
      ) : (
        <>
          <ContentComposer initialValue={props.body} disabled />
        </>
      )}
    </Flex>
  );
}
