"use client";

import { Controller, ControllerProps } from "react-hook-form";

import { Unready } from "src/components/site/Unready";

import { Thread, Visibility } from "@/api/openapi-schema";
import { CategoryBadge } from "@/components/category/CategoryBadge";
import { Byline } from "@/components/content/Byline";
import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { LinkCard } from "@/components/library/links/LinkCard";
import { CancelAction } from "@/components/site/Action/Cancel";
import { SaveAction } from "@/components/site/Action/Save";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Breadcrumbs } from "@/components/thread/Breadcrumbs";
import { ReplyBox } from "@/components/thread/ReplyBox/ReplyBox";
import { ReplyProvider } from "@/components/thread/ReplyContext";
import { ReplyList } from "@/components/thread/ReplyList/ReplyList";
import { ThreadDeletedAlert } from "@/components/thread/ThreadDeletedAlert";
import { ThreadMenu } from "@/components/thread/ThreadMenu/ThreadMenu";
import { TagListField } from "@/components/thread/ThreadTagList";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Heading } from "@/components/ui/heading";
import { HeadingInput } from "@/components/ui/heading-input";
import {
  DiscussionIcon,
  DiscussionParticipatingIcon,
} from "@/components/ui/icons/Discussion";
import { LikeIcon, LikeSavedIcon } from "@/components/ui/icons/Like";
import { VisibilityBadge } from "@/components/visibility/VisibilityBadge";
import { HStack, LStack, VStack, WStack, styled } from "@/styled-system/jsx";

import { Form, Props, useThreadScreen } from "./useThreadScreen";

export function ThreadScreen(props: Props) {
  const { ready, error, form, isEditing, isEmpty, resetKey, data, handlers } =
    useThreadScreen(props);

  if (!ready) {
    return <Unready error={error} />;
  }

  const { thread } = data;

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

            {thread.visibility !== Visibility.published && (
              <VisibilityBadge visibility={thread.visibility} />
            )}
          </WStack>
          <FormErrorText>{form.formState.errors.root?.message}</FormErrorText>

          {isEditing ? (
            <TitleInput name="title" control={form.control} />
          ) : (
            <Heading fontSize="heading.variable.1" fontWeight="bold">
              {thread.title}
            </Heading>
          )}

          {isEditing ? (
            <TagListField
              name="tags"
              control={form.control}
              initialTags={thread.tags}
            />
          ) : (
            <TagBadgeList tags={thread.tags} />
          )}

          {thread.link && <LinkCard link={thread.link} />}

          <ThreadBodyInput
            control={form.control}
            name="body"
            initialValue={thread.body}
            resetKey={resetKey}
            disabled={!isEditing}
            handleEmptyStateChange={handlers.handleEmptyStateChange}
          />
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

        <ReplyBox initialSession={props.initialSession} thread={thread} />
      </LStack>
    </ReplyProvider>
  );
}

type TitleInputProps = Omit<ControllerProps<Form>, "render">;

export function TitleInput({ control }: TitleInputProps) {
  return (
    <Controller<Form>
      render={({ field: { onChange, ...field }, formState, fieldState }) => {
        return (
          <>
            <HeadingInput
              id="title-input"
              placeholder="Thread title..."
              onValueChange={onChange}
              defaultValue={formState.defaultValues?.["title"]}
              {...field}
            />

            <FormErrorText>{fieldState.error?.message}</FormErrorText>
          </>
        );
      }}
      control={control}
      name="title"
    />
  );
}

type ThreadBodyInputProps = Omit<ControllerProps<Form>, "render"> & {
  initialValue: string;
  resetKey: string;
  handleEmptyStateChange: (isEmpty: boolean) => void;
};

function ThreadBodyInput({
  control,
  name,
  initialValue,
  resetKey,
  disabled,
  handleEmptyStateChange,
}: ThreadBodyInputProps) {
  return (
    <Controller<Form>
      render={({ field: { onChange } }) => {
        function handleChange(value: string, isEmpty: boolean) {
          handleEmptyStateChange(isEmpty);
          onChange(value);
        }

        return (
          <ContentComposer
            initialValue={initialValue}
            onChange={handleChange}
            resetKey={resetKey}
            disabled={disabled}
          />
        );
      }}
      control={control}
      name={name}
    />
  );
}

function ThreadStats({ thread }: { thread: Thread }) {
  const likeCount = thread.likes.likes;
  const likeLabel = likeCount === 1 ? "like" : "likes";
  const replyCount = thread.reply_status.replies;
  const replyLabel = replyCount === 1 ? "reply" : "replies";

  return (
    <HStack gap="4" color="fg.muted">
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
