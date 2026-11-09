"use client";

import { useReplyGet } from "@/api/openapi-client/replies";
import type { Reply, Thread } from "@/api/openapi-schema";
import { Byline } from "@/components/content/Byline";
import { ContentRenderer } from "@/components/content/ContentRenderer/ContentRenderer";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { ThreadConclusionIcon } from "@/components/ui/icons/ThreadConclusion";
import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";

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
          {error ? "The answer could not be loaded." : "Loading answer…"}
        </Text>
        {Boolean(error) && (
          <Button variant="ghost" onClick={() => void mutate()}>
            Try again
          </Button>
        )}
      </CardBox>
    );
  }

  const permalink = current ? `#${reply.id}` : `/t/locate/${reply.id}`;

  return (
    <Card
      id={`thread-conclusion-${reply.id}`}
      title="Answer"
      titleIcon={<ThreadConclusionIcon color="status.success.content" />}
      url={permalink}
      body={
        !current ? <ContentRenderer initialValue={reply.body} /> : undefined
      }
      shape="row"
      tone="success"
      controls={
        <LinkButton href={permalink} variant="ghost">
          {current ? "Jump to answer" : "View in thread"}
        </LinkButton>
      }
    >
      <Byline
        author={reply.author}
        time={new Date(reply.createdAt)}
        updated={new Date(reply.updatedAt)}
        href={permalink}
      />
    </Card>
  );
}
