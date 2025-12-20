import Link from "next/link";
import { Controller, ControllerProps } from "react-hook-form";

import { Anchor } from "src/components/site/Anchor";

import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CloseIcon } from "@/components/ui/icons/Close";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { css } from "@/styled-system/css";
import { HStack, LStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox } from "@/styled-system/patterns";
import { timestamp } from "@/utils/date";

import { useReplyContext } from "../ReplyContext";

import { Form, Props, useReplyBox } from "./useReplyBox";

export function ReplyBox(props: Props) {
  const { replyTo, clearReplyTo } = useReplyContext();
  const {
    isLoggedIn,
    isEmpty,
    isLoading,
    form,
    resetKey,
    postedReply,
    handlers,
  } = useReplyBox(props);

  if (!isLoggedIn) {
    return <LoginToReply />;
  }

  return (
    <VStack w="full" pb="12" gap="2" alignItems="stretch">
      <Admonition
        value={!!postedReply}
        onChange={handlers.handleReplyPostedAdmonitionClose}
      >
        {postedReply && (
          <LStack h="full" justifyContent="center">
            <styled.p fontSize="sm" color="fg.muted">
              Your reply has been posted on{" "}
              <Link
                className={css({
                  color: "fg.emphasized",
                  _hover: { textDecoration: "underline" },
                })}
                href={postedReply.permalink}
                onClick={handlers.handleReplyNavigation}
              >
                page {postedReply.pageNumber}
              </Link>
              .
            </styled.p>
          </LStack>
        )}
      </Admonition>

      <styled.form
        className={CardBox()}
        display="flex"
        flexDirection="column"
        gap="1"
        width="full"
        onSubmit={handlers.handleSubmit}
      >
        {replyTo && (
          <WStack py="1" px="2" borderRadius="md" bgColor="bg.muted">
            <HStack gap="1" fontSize="sm" color="fg.muted">
              <styled.span>Replying&nbsp;to</styled.span>
              <MemberIdent
                profile={replyTo.reply.author}
                name="handle"
                size="xs"
              />
              <styled.a href={`#${replyTo.reply.id}`}>
                {timestamp(replyTo.reply.createdAt)}
              </styled.a>
            </HStack>

            <IconButton
              type="button"
              size="xs"
              variant="ghost"
              aria-label="Clear reply-to"
              onClick={clearReplyTo}
            >
              <CloseIcon />
            </IconButton>
          </WStack>
        )}

        <HStack justifyContent="space-between">
          <HStack gap="1">
            <styled.span textWrap="nowrap">Reply to</styled.span>
            <MemberIdent
              profile={props.thread.author}
              name="handle"
              avatar="hidden"
            />
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
    </VStack>
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
      bgColor="border.muted"
      justifyContent="center"
    >
      <DiscussionIcon width="4" />

      <p>
        Please <Anchor href="/register">sign up or log in</Anchor> to reply
      </p>
    </HStack>
  );
}
