import Link from "next/link";

import { Reply as ReplyType, Thread } from "@/api/openapi-schema";
import { ContentComposerField } from "@/components/content/ContentComposer";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { CancelAction } from "@/components/site/Action/Cancel";
import { SaveAction } from "@/components/site/Action/Save";
import { Timestamp } from "@/components/site/Timestamp";
import { CardBox } from "@/components/ui/card-box";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { HStack, WStack, styled } from "@/styled-system/jsx";
import { hstack } from "@/styled-system/patterns";

import { Byline } from "../../content/Byline";
import { PostReviewBadge } from "../PostReviewBadge";
import { ThreadReactList } from "../ReactList/ThreadReactList";
import { ReplyMenu } from "../ReplyMenu/ReplyMenu";
import { Signature } from "../Signature";

import { ReplyToButton } from "./ReplyToButton";
import { useFragmentScroll } from "./useFragmentScroll";
import { Form, Props, useReply } from "./useReply";

export function Reply(props: Props) {
  const {
    isEmpty,
    isEditing,
    isEditingInReview,
    canManageReplies,
    resetKey,
    form,
    isConfirmingDelete,
    handlers,
  } = useReply(props);
  const isTargeted = useFragmentScroll(props.reply.id);

  const { initialSession, thread, reply, currentPage, initialSignatureConfig } =
    props;

  const isInReview = reply.visibility === "review";

  return (
    <CardBox
      id={reply.id}
      data-targeted={isTargeted || undefined}
      _target={{
        scrollMarginTop: {
          base: "0",
          md: "20",
        },
        animation: "target-pulse",
      }}
      css={{
        "&[data-targeted]": {
          animation: "target-pulse",
        },
        backgroundColor: isInReview ? "status.warning.surface/30" : undefined,
      }}
    >
      <styled.form
        display="flex"
        flexDirection="column"
        gap="2"
        onSubmit={handlers.handleSave}
      >
        <WStack>
          <Byline
            href={`#${reply.id}`}
            author={reply.author}
            time={new Date(reply.createdAt)}
            updated={new Date(reply.updatedAt)}
          />

          {isEditing ? (
            <HStack>
              <>
                <CancelAction
                  type="button"
                  onClick={handlers.handleDiscardChanges}
                >
                  Discard
                </CancelAction>
                <SaveAction type="submit" disabled={isEmpty}>
                  {isEditingInReview ? "Accept" : "Save"}
                </SaveAction>
              </>
            </HStack>
          ) : (
            <HStack>
              <ReplyToButton thread={thread} reply={reply} />
              <ReplyMenu
                thread={thread}
                reply={reply}
                currentPage={currentPage}
                onEdit={handlers.handleSetEditing}
              />
            </HStack>
          )}
        </WStack>

        {reply.reply_to && <InReplyTo to={reply.reply_to} thread={thread} />}

        <ContentComposerField
          control={form.control}
          name="body"
          initialValue={reply.body}
          resetKey={resetKey}
          disabled={!isEditing}
          onEmptyStateChange={handlers.handleEmptyStateChange}
        />

        {initialSignatureConfig.enabled && (
          <Signature
            signature={reply.author.signature}
            maxHeight={initialSignatureConfig.maxHeight}
          />
        )}
      </styled.form>

      <WStack>
        <ThreadReactList
          initialSession={initialSession}
          thread={thread}
          reply={reply}
          currentPage={currentPage}
        />

        {isInReview && (
          <PostReviewBadge
            isModerator={canManageReplies}
            postId={reply.id}
            onAccept={handlers.handleAcceptReply}
            onEditAndAccept={handlers.handleSetEditingInReview}
            onDelete={handlers.handleConfirmDelete}
            isConfirmingDelete={isConfirmingDelete}
            onCancelDelete={handlers.handleCancelDelete}
          />
        )}
      </WStack>
    </CardBox>
  );
}

function InReplyTo({ to, thread }: { to: ReplyType; thread: Thread }) {
  // figure out if the reply-to is on the current page, then  do a fragment link
  // if on same page, otherwise use /t/locate to navigate to the right page.
  const isOnCurrentPage = thread.replies.replies.some((r) => r.id === to.id);
  const permalink = isOnCurrentPage ? `#${to.id}` : `/t/locate/${to.id}`;

  // NOTE: because nextjs does some weird shit, we gotta use a normal anchor
  // for fragment navigation, otherwise it breaks :target etc for some reason.
  const AnchorComponent = isOnCurrentPage ? styled.a : Link;

  return (
    <WStack
      gap="1"
      fontSize="xs"
      color="text.subtle"
      px="2"
      py="1"
      borderRadius="md"
      bgColor="background.inset"
      w="full"
      minW="0"
    >
      <AnchorComponent
        href={permalink}
        className={hstack({
          minW: "0",
          flexShrink: "1",
        })}
      >
        <ReplyIcon w="4" minW="4" />
        <styled.span
          minW="0"
          overflow="hidden"
          textOverflow="ellipsis"
          whiteSpace="nowrap"
          lineClamp="1"
        >
          “{to.description}”
        </styled.span>
      </AnchorComponent>

      <HStack flexShrink="0" minW="0">
        <MemberBadge
          profile={to.author}
          size="xs"
          name="handle"
          avatar="visible"
        />
        <AnchorComponent href={permalink}>
          <Timestamp created={to.createdAt} />
        </AnchorComponent>
      </HStack>
    </WStack>
  );
}
