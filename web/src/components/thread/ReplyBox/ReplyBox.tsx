import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";
import { Controller, ControllerProps } from "react-hook-form";

import { Thread } from "src/api/openapi-schema";
import { Anchor } from "src/components/site/Anchor";

import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { ProfilePill } from "@/components/site/ProfilePill/ProfilePill";
import { Button } from "@/components/ui/button";
import { Box, HStack, styled } from "@/styled-system/jsx";
import { CardBox } from "@/styled-system/patterns";

import { Form, useReplyBox } from "./useReplyBox";

export type Props = {
  thread: Thread;
};

export function ReplyBox(props: Thread) {
  const { isLoggedIn, isEmpty, isLoading, form, resetKey, handlers } =
    useReplyBox(props);

  if (!isLoggedIn) {
    return <LoginToReply />;
  }

  return (
    <Box
      w="full"
      pb="12" // Provide spacing at the bottom for the editor's menu + navbar.
    >
      <styled.form
        className={CardBox()}
        display="flex"
        flexDirection="column"
        gap="1"
        width="full"
        onSubmit={handlers.handleSubmit}
      >
        <HStack justifyContent="space-between">
          <HStack gap="1">
            Reply to{" "}
            <ProfilePill profileReference={props.author} showAvatar={false} />
          </HStack>

          <Button type="submit" size="xs" disabled={isLoading || isEmpty}>
            Post
          </Button>
        </HStack>

        <ReplyBodyInput
          name="body"
          control={form.control}
          handleEmptyStateChange={handlers.handleEmptyStateChange}
          resetKey={resetKey}
        />
      </styled.form>
    </Box>
  );
}

type ReplyBodyInputProps = Omit<ControllerProps<Form>, "render"> & {
  handleEmptyStateChange: (isEmpty: boolean) => void;
  resetKey: string;
};

function ReplyBodyInput({
  control,
  name,
  handleEmptyStateChange,
  resetKey,
}: ReplyBodyInputProps) {
  return (
    <Controller<Form>
      render={({ field: { onChange } }) => {
        function handleChange(value: string, isEmpty: boolean) {
          handleEmptyStateChange(isEmpty);
          onChange(value);
        }

        return <ContentComposer onChange={handleChange} resetKey={resetKey} />;
      }}
      control={control}
      name={name}
    />
  );
}

function LoginToReply() {
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
