import { Button } from "@chakra-ui/react";

import { Thread } from "src/api/openapi/schemas";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";
import { Bold } from "src/components/content/ContentComposer/controls/Bold";
import { Italic } from "src/components/content/ContentComposer/controls/Italic";

import { HStack, styled } from "@/styled-system/jsx";

import { useReplyBox } from "./useReplyBox";

export function ReplyBox(props: Thread) {
  const { onChange, onReply, onKeyDown, isLoading, resetKey } =
    useReplyBox(props);

  return (
    <styled.form
      display="flex"
      flexDirection="column"
      width="full"
      borderRadius="2xl"
      p={2}
      alignItems="end"
      onKeyDown={onKeyDown}
    >
      <ContentComposer
        onChange={onChange}
        minHeight="8em"
        disabled={isLoading}
        resetKey={resetKey}
      >
        <HStack>
          <Bold />
          <Italic />
        </HStack>
      </ContentComposer>

      <HStack mt={4} justifyContent="end">
        <Button type="submit" onClick={onReply} isLoading={isLoading}>
          Post
        </Button>
      </HStack>
    </styled.form>
  );
}
