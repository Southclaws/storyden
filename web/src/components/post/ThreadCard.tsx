import Link from "next/link";
import { memo } from "react";

import { ThreadReference } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { Card } from "@/components/ui/rich-card";
import { HStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { CategoryBadge } from "../category/CategoryBadge";
import { ThreadMenu } from "../thread/ThreadMenu/ThreadMenu";
import {
  DiscussionIcon,
  DiscussionParticipatingIcon,
} from "../ui/icons/Discussion";

import { LikeButton } from "./LikeButton/LikeButton";

type Props = {
  thread: ThreadReference;
  hideCategoryBadge?: boolean;
};

export const ThreadReferenceCard = memo(({ thread, hideCategoryBadge = false }: Props) => {
  const session = useSession();
  const permalink = `/t/${thread.slug}`;

  const title = thread.title || thread.link?.title || "Untitled post";

  const hasReplied = thread.reply_status.replied > 0;
  const replyCount = thread.reply_status.replies;
  const replyCountLabel =
    replyCount === 1 ? `1 reply` : `${replyCount} replies`;

  const replyStatusLabel = hasReplied
    ? `${replyCountLabel} (including you!)`
    : replyCountLabel;

  return (
    <Card
      shape="responsive"
      id={thread.id}
      title={title}
      text={thread.description}
      url={permalink}
      image={getAssetURL(
        thread.assets?.[0]?.path ?? thread.link?.primary_image?.path,
      )}
      controls={
        session && (
          <HStack>
            {!hideCategoryBadge && thread.category && <CategoryBadge category={thread.category} />}
            <LikeButton thread={thread} />
            <CollectionMenu account={session} thread={thread} />
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
          <HStack justify="space-between">
            <Link href={permalink} title={replyStatusLabel}>
              <styled.span color="fg.subtle" display="flex" gap="1">
                {hasReplied ? (
                  <DiscussionParticipatingIcon width="4" />
                ) : (
                  <DiscussionIcon width="4" />
                )}
                {replyCount}
              </styled.span>
            </Link>
          </HStack>
        }
      />
    </Card>
  );
});

ThreadReferenceCard.displayName = "ThreadReferenceCard";
