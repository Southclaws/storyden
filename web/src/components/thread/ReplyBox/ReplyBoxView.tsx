import Link from "next/link";
import { useId } from "react";

import type { Reply, Thread } from "@/api/openapi-schema";
import { ContentComposerField } from "@/components/content/ContentComposer";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Anchor } from "@/components/site/Anchor";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { IconButton } from "@/components/ui/icon-button";
import { CloseIcon } from "@/components/ui/icons/Close";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { ThreadLockedIcon } from "@/components/ui/icons/ThreadLock";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import {
  Box,
  HStack,
  LStack,
  VStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { timestamp } from "@/utils/date";

import type { useReplyBox } from "./useReplyBox";

export type ReplyBoxViewProps = {
  thread: Pick<Thread, "author" | "locked_at">;
  state: ReturnType<typeof useReplyBox>;
  replyTo: Reply | null;
  registrationAllowed: boolean;
  isActive: boolean;
  focusOnReady: boolean;
  onActivate: () => void;
  onWarmEditor: () => void;
  onClearReplyTo: () => void;
};

export function ReplyBoxView(props: ReplyBoxViewProps) {
  const { replyTo, isActive, focusOnReady } = props;
  const labelID = useId();
  const lockNoticeID = useId();
  const {
    isLoggedIn,
    canReply,
    isEmpty,
    isLoading,
    form,
    draftKey,
    postedReply,
    handlers,
  } = props.state;
  const isReplyingDisabled = Boolean(props.thread.locked_at) && !canReply;

  if (!isLoggedIn && !isReplyingDisabled) {
    return <LoginToReply registrationAllowed={props.registrationAllowed} />;
  }

  return (
    <VStack
      className="reply-composer"
      w="full"
      minW="0"
      pb="12"
      gap="2"
      alignItems="stretch"
    >
      <Admonition
        value={!!postedReply}
        onChange={handlers.handleReplyPostedAdmonitionClose}
      >
        {postedReply && (
          <LStack h="full" justifyContent="center">
            <Text variant="supporting">
              Your reply has been posted on{" "}
              <Link
                className={css({
                  color: "interactive.emphasized.content",
                  _hover: { textDecoration: "underline" },
                })}
                href={postedReply.permalink}
                onClick={handlers.handleReplyNavigation}
              >
                page {postedReply.pageNumber}
              </Link>
              .
            </Text>
          </LStack>
        )}
      </Admonition>

      <styled.form
        className="reply-composer__form"
        width="full"
        minW="0"
        aria-labelledby={labelID}
        aria-disabled={isReplyingDisabled}
        aria-describedby={props.thread.locked_at ? lockNoticeID : undefined}
        onPointerEnter={(event) => {
          if (isReplyingDisabled) {
            return;
          }

          if (event.pointerType === "mouse") {
            props.onWarmEditor();
          }
        }}
        onPointerDown={(event) => {
          // Remember a click on the placeholder while the editor downloads.
          if (
            !isReplyingDisabled &&
            event.target instanceof Element &&
            event.target.closest(".content-composer--readonly")
          ) {
            props.onActivate();
          }
        }}
        onSubmit={(event) => {
          if (isReplyingDisabled) {
            event.preventDefault();
            return;
          }

          props.onWarmEditor();
          return handlers.handleSubmit(event);
        }}
      >
        <FormControl gap="3" minW="0">
          <LStack gap="2" alignItems="stretch" minW="0">
            <FormLabel id={labelID} mb="0">
              Reply to thread
            </FormLabel>
            {replyTo && (
              <HStack gap="2" minW="0">
                <HStack
                  gap="1"
                  minW="0"
                  flex="1"
                  fontSize="sm"
                  color="text.muted"
                >
                  <styled.span flexShrink="0">Replying to</styled.span>
                  <Box minW="0" flex="1">
                    <MemberIdent
                      profile={replyTo.author}
                      name="handle"
                      size="xs"
                    />
                  </Box>
                  <styled.a href={`#${replyTo.id}`} flexShrink="0">
                    {timestamp(replyTo.createdAt)}
                  </styled.a>
                </HStack>
                <IconButton
                  type="button"
                  variant="ghost"
                  aria-label="Clear reply-to"
                  disabled={isReplyingDisabled}
                  flexShrink="0"
                  onClick={props.onClearReplyTo}
                >
                  <CloseIcon />
                </IconButton>
              </HStack>
            )}
          </LStack>

          <styled.div
            className="reply-composer__editor"
            p="3"
            minH="24"
            minW="0"
            bg="background.control"
            borderWidth="thin"
            borderColor="border.default"
            borderRadius="md"
          >
            {isActive ? (
              <ContentComposerField
                className={css({ pe: "8" })}
                name="body"
                control={form.control}
                disabled={isReplyingDisabled}
                focusOnReady={focusOnReady || !!replyTo}
                placeholder="Write a reply…"
                onEmptyStateChange={handlers.handleEmptyStateChange}
                key={draftKey}
              />
            ) : (
              <Button
                className="content-composer content-composer--idle reply-composer__trigger"
                type="button"
                variant="ghost"
                disabled={isReplyingDisabled}
                aria-label="Write a reply"
                w="full"
                minH="16"
                px="0"
                alignItems="start"
                justifyContent="start"
                fontWeight="normal"
                fontSize="md"
                color="text.muted"
                textAlign="start"
                whiteSpace="normal"
                onFocus={props.onActivate}
                onClick={props.onActivate}
              >
                {isReplyingDisabled ? "Replies are locked." : "Write a reply…"}
              </Button>
            )}
          </styled.div>

          <WStack gap="3" alignItems="center" minW="0">
            {props.thread.locked_at && (
              <HStack
                gap="2"
                alignItems="center"
                color="text.muted"
                flex="1"
                minW="0"
              >
                <ThreadLockedIcon
                  w="4"
                  h="4"
                  flexShrink="0"
                  aria-hidden="true"
                />
                <FormHelperText id={lockNoticeID} color="text.muted">
                  Thread locked ·{" "}
                  {canReply
                    ? "You have permission to reply still."
                    : "Only members with permission may reply."}
                </FormHelperText>
              </HStack>
            )}
            <Button
              type="submit"
              ml="auto"
              flexShrink="0"
              disabled={isReplyingDisabled || isEmpty}
              loading={isLoading}
              loadingText="Posting…"
            >
              Post reply
            </Button>
          </WStack>
        </FormControl>
      </styled.form>
    </VStack>
  );
}

function LoginToReply({
  registrationAllowed,
}: {
  registrationAllowed: boolean;
}) {
  const action = registrationAllowed ? "sign up or log in" : "log in";

  return (
    <HStack
      w="full"
      p="8"
      borderRadius="xl"
      bgColor="border.strong"
      justifyContent="center"
    >
      <DiscussionIcon width="4" />

      <p>
        Please{" "}
        <Anchor href={registrationAllowed ? "/register" : "/login"}>
          {action}
        </Anchor>{" "}
        to reply
      </p>
    </HStack>
  );
}
