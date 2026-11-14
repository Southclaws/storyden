"use client";

import { useReplyGet } from "@/api/openapi-client/replies";
import type { Reply, Thread } from "@/api/openapi-schema";
import { Byline } from "@/components/content/Byline";
import { ContentRenderer } from "@/components/content/ContentRenderer/ContentRenderer";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { ArrowDownIcon } from "@/components/ui/icons/Arrow";
import { ThreadConclusionIcon } from "@/components/ui/icons/ThreadConclusion";
import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";
import { Box, HStack, WStack } from "@/styled-system/jsx";

type Props = {
  thread: Thread;
  initialConclusionReply?: Reply;
};

export function ThreadConclusion({ thread, initialConclusionReply }: Props) {
  const id = thread.conclusion_reply_id;
  const current = thread.replies.replies.find((reply) => reply.id === id);
  const initial =
    initialConclusionReply?.id === id ? initialConclusionReply : undefined;
  const { data, error, mutate } = useReplyGet(thread.slug, id ?? "", {
    swr: {
      enabled: Boolean(id) && !current,
      fallbackData: current ?? initial,
    },
  });

  if (!id) {
    return null;
  }

  const reply = current ?? data;
  if (reply?.id !== id || reply.deletedAt || reply.visibility !== "published") {
    return (
      <CardBox role="status" tone="success">
        <Text variant="supporting">
          {error
            ? "The highlighted reply could not be loaded."
            : "Loading highlighted reply…"}
        </Text>
        {Boolean(error) && (
          <Button variant="ghost" onClick={() => void mutate()}>
            Try again
          </Button>
        )}
      </CardBox>
    );
  }

  return <HighlightedReply reply={reply} onCurrentPage={Boolean(current)} />;
}

type HighlightedReplyProps = {
  reply: Reply;
  onCurrentPage: boolean;
};

export function HighlightedReply({
  reply,
  onCurrentPage,
}: HighlightedReplyProps) {
  const permalink = onCurrentPage ? `#${reply.id}` : `/t/locate/${reply.id}`;
  const byline = (
    <Byline
      author={reply.author}
      time={new Date(reply.createdAt)}
      updated={new Date(reply.updatedAt)}
      href={permalink}
    />
  );

  if (onCurrentPage) {
    return (
      <CardBox as="nav" aria-label="Highlighted reply" tone="success" w="full">
        <WStack gap="3">
          <HStack gap="3" flexWrap="wrap" minW="0">
            <HStack gap="2">
              <ThreadConclusionIcon
                w="4"
                h="4"
                color="status.success.content"
                aria-hidden="true"
              />
              <Text
                as="span"
                variant="supporting"
                color="text.default"
                fontWeight="semibold"
              >
                Highlighted reply below
              </Text>
            </HStack>
            <Text variant="supporting">•</Text>
            <Box w="fit-content" maxW="full" minW="0">
              {byline}
            </Box>
          </HStack>

          <LinkButton href={permalink} variant="ghost">
            Jump to reply
            <ArrowDownIcon aria-hidden="true" />
          </LinkButton>
        </WStack>
      </CardBox>
    );
  }

  return (
    <Card
      id={`thread-conclusion-${reply.id}`}
      title="Highlighted reply"
      titleIcon={<ThreadConclusionIcon color="status.success.content" />}
      url={permalink}
      body={<ContentRenderer initialValue={reply.body} />}
      shape="row"
      tone="success"
      controls={
        <LinkButton href={permalink} variant="ghost">
          View reply
        </LinkButton>
      }
    >
      {byline}
    </Card>
  );
}
