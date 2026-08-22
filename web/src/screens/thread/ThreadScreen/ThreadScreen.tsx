"use client";

import { match } from "ts-pattern";

import { Thread, Visibility } from "@/api/openapi-schema";
import { CategoryBadge } from "@/components/category/CategoryBadge";
import { Byline } from "@/components/content/Byline";
import { ContentComposerField } from "@/components/content/ContentComposer";
import { LinkCard } from "@/components/library/links/LinkCard";
import { CancelAction } from "@/components/site/Action/Cancel";
import { SaveAction } from "@/components/site/Action/Save";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Breadcrumbs } from "@/components/thread/Breadcrumbs";
import { PostReviewBadge } from "@/components/thread/PostReviewBadge";
import { ReplyBox } from "@/components/thread/ReplyBox/ReplyBox";
import { ReplyProvider } from "@/components/thread/ReplyContext";
import { ReplyList } from "@/components/thread/ReplyList/ReplyList";
import { Signature } from "@/components/thread/Signature";
import { ThreadDeletedAlert } from "@/components/thread/ThreadDeletedAlert";
import { ThreadMenu } from "@/components/thread/ThreadMenu/ThreadMenu";
import { ThreadTagListField } from "@/components/thread/ThreadTagList.field";
import { FormErrorText } from "@/components/ui/form-error-text";
import { HeadingInputField } from "@/components/ui/heading-input";
import {
  DiscussionIcon,
  DiscussionParticipatingIcon,
} from "@/components/ui/icons/Discussion";
import { LikeIcon, LikeSavedIcon } from "@/components/ui/icons/Like";
import { PageHeading } from "@/components/ui/page-heading";
import { VisibilityBadge } from "@/components/visibility/VisibilityBadge";
import { HStack, LStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

import { Form, Props, useThreadScreen } from "./useThreadScreen";

export function ThreadScreen(props: Props) {
  const {
    ready,
    error,
    form,
    isEditing,
    isEmpty,
    resetKey,
    isModerator,
    isConfirmingDelete,
    data,
    handlers,
  } = useThreadScreen(props);

  if (!ready) {
    return <Unready error={error} />;
  }

  const { thread } = data;
  const signatureConfig = props.initialSignatureConfig ?? {
    enabled: true,
    maxHeight: 160,
  };

  return (
    <ReplyProvider>
      <LStack gap="4">
        <styled.form
          display="flex"
          flexDirection="column"
          alignItems="start"
          gap="1"
          width="full"
          onSubmit={handlers.handleSave}
        >
          <WStack alignItems="start">
            <Breadcrumbs thread={thread} />

            <HStack>
              {isEditing && (
                <>
                  <CancelAction
                    type="button"
                    onClick={handlers.handleDiscardChanges}
                  >
                    Discard
                  </CancelAction>
                  <SaveAction type="submit" disabled={isEmpty}>
                    Save
                  </SaveAction>
                </>
              )}

              <ThreadMenu thread={thread} editingEnabled movingEnabled />
            </HStack>
          </WStack>

          {thread.deletedAt !== undefined && (
            <ThreadDeletedAlert thread={thread} />
          )}
          <WStack justifyContent="space-between">
            <HStack>
              <Byline
                href={`#${thread.id}`}
                author={thread.author}
                time={new Date(thread.createdAt)}
                updated={new Date(thread.updatedAt)}
              />

              {thread.category && <CategoryBadge category={thread.category} />}
            </HStack>

            {match(thread.visibility)
              .with(Visibility.published, () => null)
              .with(Visibility.review, () => (
                <PostReviewBadge
                  isModerator={isModerator}
                  postId={thread.id}
                  onAccept={handlers.handleAcceptThread}
                  onEditAndAccept={handlers.handleEditAndAccept}
                  onDelete={handlers.handleConfirmDelete}
                  isConfirmingDelete={isConfirmingDelete}
                  onCancelDelete={handlers.handleCancelDelete}
                />
              ))
              .otherwise(() => (
                <VisibilityBadge visibility={thread.visibility} />
              ))}
          </WStack>
          <FormErrorText>{form.formState.errors.root?.message}</FormErrorText>

          {isEditing ? (
            <>
              <HeadingInputField
                id="title-input"
                placeholder="Thread title..."
                name="title"
                control={form.control}
              />
              <FormErrorText>
                {form.formState.errors.title?.message}
              </FormErrorText>
            </>
          ) : (
            <PageHeading>{thread.title}</PageHeading>
          )}

          {isEditing ? (
            <ThreadTagListField
              name="tags"
              control={form.control}
              initialTags={thread.tags}
            />
          ) : (
            <TagBadgeList tags={thread.tags} />
          )}

          {thread.link && <LinkCard link={thread.link} />}

          <ContentComposerField
            control={form.control}
            name="body"
            initialValue={thread.body}
            resetKey={resetKey}
            disabled={!isEditing}
            onEmptyStateChange={handlers.handleEmptyStateChange}
          />

          {signatureConfig.enabled && (
            <Signature
              signature={thread.author.signature}
              maxHeight={signatureConfig.maxHeight}
            />
          )}
        </styled.form>

        <ThreadStats thread={thread} />

        <VStack w="full">
          {data.thread.replies.total_pages > 1 && (
            <PaginationControls
              path={`/t/${thread.slug}`}
              currentPage={data.thread.replies.current_page ?? 1}
              totalPages={data.thread.replies.total_pages}
              pageSize={data.thread.replies.page_size}
            />
          )}

          <ReplyList
            initialSession={props.initialSession}
            thread={thread}
            currentPage={data.thread.replies.current_page}
            initialSignatureConfig={signatureConfig}
          />

          {data.thread.replies.total_pages > 1 && (
            <PaginationControls
              path={`/t/${thread.slug}`}
              currentPage={data.thread.replies.current_page ?? 1}
              totalPages={data.thread.replies.total_pages}
              pageSize={data.thread.replies.page_size}
            />
          )}
        </VStack>

        <ReplyBox
          initialSession={props.initialSession}
          initialSettings={props.initialSettings}
          thread={thread}
        />
      </LStack>
    </ReplyProvider>
  );
}

function ThreadStats({ thread }: { thread: Thread }) {
  const likeCount = thread.likes.likes;
  const likeLabel = pluralise(likeCount, "like");
  const replyCount = thread.reply_status.replies;
  const replyLabel = pluralise(replyCount, "reply", "replies");

  return (
    <HStack gap="4" color="text.subtle">
      <styled.span
        display="flex"
        gap="1"
        alignItems="center"
        title={thread.likes.liked ? "You liked this thread" : undefined}
      >
        <span>
          {thread.likes.liked ? (
            <LikeSavedIcon width="4" />
          ) : (
            <LikeIcon width="4" />
          )}
        </span>
        <span>
          {likeCount} {likeLabel}
        </span>
      </styled.span>

      <styled.span
        display="flex"
        gap="1"
        alignItems="center"
        title={
          thread.reply_status.replied
            ? "You have replied to this thread"
            : undefined
        }
      >
        <span>
          {thread.reply_status.replied ? (
            <DiscussionParticipatingIcon width="4" />
          ) : (
            <DiscussionIcon width="4" />
          )}
        </span>
        <span>
          {replyCount} {replyLabel}
        </span>
      </styled.span>
    </HStack>
  );
}
