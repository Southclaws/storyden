import { Thread } from "src/api/openapi-schema";
import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { ProfilePill } from "@/components/site/ProfilePill/ProfilePill";
import { Button } from "@/components/ui/button";
import { CardBox, HStack, styled } from "@/styled-system/jsx";

import { useReplyBox } from "./useReplyBox";

export function ReplyBox(props: Thread) {
  const { onChange, onReply, onKeyDown, isLoading, isEmpty, resetKey } =
    useReplyBox(props);

  return (
    <CardBox>
      <styled.form
        display="flex"
        flexDirection="column"
        gap="1"
        width="full"
        borderRadius="2xl"
        onKeyDown={onKeyDown}
      >
        <HStack justifyContent="space-between">
          <HStack gap="1">
            Reply to{" "}
            <ProfilePill profileReference={props.author} showAvatar={false} />
          </HStack>

          <Button
            type="submit"
            size="xs"
            onClick={onReply}
            disabled={isLoading || isEmpty}
          >
            Post
          </Button>
        </HStack>

        <ContentComposer
          onChange={onChange}
          disabled={isLoading}
          resetKey={resetKey}
        />
      </styled.form>
    </CardBox>
  );
}
