import Link from "next/link";
import { memo } from "react";

import {
  Permission,
  ThreadReference,
  Visibility,
} from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { Card } from "@/components/ui/rich-card";
import { Box, HStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";
import { timestamp } from "@/utils/date";
import { hasPermission } from "@/utils/permissions";

import { CategoryBadge } from "../category/CategoryBadge";
import { PostReviewBadge } from "../thread/PostReviewBadge";
import { ThreadMenu } from "../thread/ThreadMenu/ThreadMenu";
import {
  DiscussionIcon,
  DiscussionParticipatingIcon,
} from "../ui/icons/Discussion";

import { LikeButton } from "./LikeButton/LikeButton";
import { useThreadCardModeration } from "./useThreadCardModeration";

type Props = {
  thread: ThreadReference;
  hideCategoryBadge?: boolean;
};

export const ThreadReferenceCard = memo(
  ({ thread, hideCategoryBadge = false }: Props) => {
    const session = useSession();
    const isDraft = thread.visibility === Visibility.draft;
    const permalink = isDraft ? `/new?id=${thread.id}` : `/t/${thread.slug}`;
    const isModerator = hasPermission(
      session,
      Permission.MANAGE_POSTS,
      Permission.ADMINISTRATOR,
    );

    const { isConfirmingDelete, handlers } = useThreadCardModeration(thread);

    const title = thread.title || thread.link?.title || "Untitled post";

    const hasReplied = thread.reply_status.replied > 0;
    const replyCount = thread.reply_status.replies;
    const replyCountLabel =
      replyCount === 1 ? `1 reply` : `${replyCount} replies`;

    const replyStatusLabel = hasReplied
      ? `${replyCountLabel} (including you!)`
      : replyCountLabel;

    const newRepliesCount = thread.read_status?.replies_since ?? 0;
    const lastReadAt = thread.read_status?.last_read_at;
    const newRepliesLabel =
      newRepliesCount > 0 && lastReadAt
        ? `${newRepliesCount} ${newRepliesCount === 1 ? "reply" : "replies"} since you last visited ${timestamp(lastReadAt, false)} ago`
        : undefined;

    const isInReview = thread.visibility === Visibility.review;

    // Don't show images for in-review threads, too noisy.
    const image = isInReview
      ? undefined
      : getAssetURL(
          thread.assets?.[0]?.path ?? thread.link?.primary_image?.path,
        );

    return (
      <Card
        shape="responsive"
        id={thread.id}
        title={title}
        text={thread.description}
        url={permalink}
        image={image}
        controls={
          session && (
            <HStack>
              {!hideCategoryBadge && thread.category && (
                <CategoryBadge category={thread.category} />
              )}
              {isInReview ? (
                <>
                  <PostReviewBadge
                    isModerator={isModerator}
                    postId={thread.id}
                    onAccept={handlers.handleAcceptThread}
                    onEditAndAccept={handlers.handleEditAndAccept}
                    onDelete={handlers.handleConfirmDelete}
                    isConfirmingDelete={isConfirmingDelete}
                    onCancelDelete={handlers.handleCancelDelete}
                  />
                </>
              ) : (
                <>
                  <LikeButton thread={thread} showCount />
                  <CollectionMenu account={session} thread={thread} />
                </>
              )}
              <ThreadMenu thread={thread} />
            </HStack>
          )
        }
      >
        <Byline
          href={permalink}
          author={thread.author}
          time={new Date(thread.createdAt)}
          updated={new Date(thread.updatedAt)}
          more={
            <Box
              className="thread-byline__more"
              flexShrink="0"
              overflow="hidden"
            >
              <Link
                className="thread-byline__anchor"
                href={permalink}
                title={replyStatusLabel}
              >
                <styled.span
                  className="thread-byline__reply-status-label"
                  color="fg.muted"
                  display="flex"
                  gap="0.5"
                  alignItems="center"
                >
                  {hasReplied ? (
                    <DiscussionParticipatingIcon width="4" />
                  ) : (
                    <DiscussionIcon width="4" />
                  )}
                  {replyCount}
                  {newRepliesCount > 0 && (
                    <styled.span
                      color="fg.muted"
                      fontSize="xs"
                      title={newRepliesLabel}
                    >
                      +{newRepliesCount}
                    </styled.span>
                  )}
                </styled.span>
              </Link>
            </Box>
          }
        />
      </Card>
    );
  },
);

ThreadReferenceCard.displayName = "ThreadReferenceCard";
