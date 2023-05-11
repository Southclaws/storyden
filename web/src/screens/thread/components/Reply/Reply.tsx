import { HStack, Spinner, Text, VStack } from "@chakra-ui/react";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";
import { Thread } from "src/api/openapi/schemas";
import { ReplyBox } from "src/components/ReplyBox";
import { Anchor } from "src/components/site/Anchor";
import { useReply } from "./useReply";

export function Reply(props: Thread) {
  const { loggedIn, onReply, isLoading } = useReply(props);
  // NOTE: isLoading is a hack to easily reset the ReplyBox + provide feedback.

  if (loggedIn) {
    return isLoading ? (
      <VStack width="full" py={6}>
        <Spinner />
      </VStack>
    ) : (
      <ReplyBox onSave={onReply} />
    );
  }

  return (
    <HStack
      w="full"
      p={8}
      borderRadius="xl"
      bgColor="blackAlpha.50"
      justifyContent="center"
    >
      <ChatBubbleLeftRightIcon width="1.5em" />

      <Text>
        Please <Anchor href="/auth">sign up or log in</Anchor> to reply
      </Text>
    </HStack>
  );
}
