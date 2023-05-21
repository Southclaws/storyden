import { Button, Flex, HStack } from "@chakra-ui/react";
import { PostProps } from "src/api/openapi/schemas";
import { Markdown } from "src/components/Markdown";
import { Editor } from "src/components/Editor";
import { Byline } from "../Byline";
import { PostMenu } from "../PostMenu/PostMenu";
import { ReactList } from "../ReactList/ReactList";
import { usePostView } from "./usePostView";

type Props = PostProps & {
  slug?: string;
};

export function PostView(props: Props) {
  const {
    isEditing,
    editingContent,
    setEditingContent,
    onPublishEdit,
    onCancelEdit,
  } = usePostView(props);

  return (
    <Flex id={props.id} flexDir="column" gap={2}>
      <Byline
        href={`#${props.id}`}
        author={props.author.handle}
        time={new Date(props.createdAt)}
        updated={new Date(props.updatedAt)}
        more={<PostMenu {...props} />}
      />
      {isEditing ? (
        <>
          <Editor onChange={setEditingContent} value={editingContent} />
          <HStack>
            <Button onClick={onPublishEdit}>Update</Button>
            <Button variant="outline" onClick={onCancelEdit}>
              Cancel
            </Button>
          </HStack>
        </>
      ) : (
        <Markdown>{props.body}</Markdown>
      )}
      <ReactList {...props} />
    </Flex>
  );
}
