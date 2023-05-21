import {
  Button,
  Code,
  Flex,
  HStack,
  Heading,
  Link,
  ListItem,
  OrderedList,
  Text,
  UnorderedList,
} from "@chakra-ui/react";
import ReactMarkdown from "react-markdown";
import { SpecialComponents } from "react-markdown/lib/ast-to-react";
import { NormalComponents } from "react-markdown/lib/complex-types";
import { PostProps } from "src/api/openapi/schemas";
import { Byline } from "../Byline";
import { ReactList } from "../ReactList/ReactList";
import { PostMenu } from "../PostMenu/PostMenu";
import { usePostView } from "./usePostView";
import { Editor } from "src/components/Editor";

const components: Partial<
  Omit<NormalComponents, keyof SpecialComponents> & SpecialComponents
> = {
  h1: (props) => (
    <Heading as="h1" variant="h1">
      {props.children}
    </Heading>
  ),
  h2: (props) => (
    <Heading as="h2" variant="h2">
      {props.children}
    </Heading>
  ),
  h3: (props) => (
    <Heading as="h3" variant="h3">
      {props.children}
    </Heading>
  ),
  h4: (props) => (
    <Heading as="h4" variant="h4">
      {props.children}
    </Heading>
  ),
  h5: (props) => (
    <Heading as="h5" variant="h5">
      {props.children}
    </Heading>
  ),

  // Typography
  p: (props) => (
    <Text overflowWrap="break-word" wordBreak="break-word" overflowX="clip">
      {props.children}
    </Text>
  ),
  a: ({ href, children }) => <Link href={href ?? "#"}>{children}</Link>,

  // Lists
  ul: (props) => <UnorderedList ml="2em">{props.children}</UnorderedList>,
  ol: (props) => <OrderedList ml="2em">{props.children}</OrderedList>,
  li: (props) => <ListItem>{props.children}</ListItem>,

  // Code
  pre: (props) => (
    <Code
      display="block"
      whiteSpace="pre"
      overflowX="scroll"
      padding={2}
      borderRadius="md"
    >
      {props.children}
    </Code>
  ),
  code: (props) => <Code>{props.children}</Code>,

  // Headings
  td: (props) => <td>{props.children}</td>,
  th: (props) => <th>{props.children}</th>,
  tr: (props) => <tr>{props.children}</tr>,
};

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
        <ReactMarkdown components={components}>{props.body}</ReactMarkdown>
      )}
      <ReactList {...props} />
    </Flex>
  );
}
