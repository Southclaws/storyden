import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";

import { Thread } from "src/api/openapi-schema";
import { Anchor } from "src/components/site/Anchor";

import { HStack } from "@/styled-system/jsx";

import { ReplyBox } from "./ReplyBox/ReplyBox";
import { useReply } from "./useReply";

export function Reply(props: Thread) {
  const { loggedIn } = useReply();
  // NOTE: isLoading is a hack to easily reset the ReplyBox + provide feedback.

  if (loggedIn) {
    return <ReplyBox {...props} />;
  }

  return (
    <HStack
      w="full"
      p="8"
      borderRadius="xl"
      bgColor="blackAlpha.50"
      justifyContent="center"
    >
      <ChatBubbleLeftRightIcon width="1.5em" />

      <p>
        Please <Anchor href="/register">sign up or log in</Anchor> to reply
      </p>
    </HStack>
  );
}
