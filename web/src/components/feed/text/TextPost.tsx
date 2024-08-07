import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { FeedItemMenu } from "../common/FeedItemMenu/FeedItemMenu";

import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
  onDelete?: () => void;
};

export function TextPost({ thread, onDelete }: Props) {
  const session = useSession();
  const permalink = `/t/${thread.slug}`;

  return (
    <Card
      shape="row"
      id={thread.id}
      title={thread.title}
      text={thread.short}
      url={permalink}
      image={thread.assets[0]?.url}
      controls={
        session && (
          <HStack>
            <CollectionMenu thread={thread} />
            <FeedItemMenu thread={thread} onDelete={onDelete} />
          </HStack>
        )
      }
    >
      <Byline
        href={permalink}
        author={thread.author}
        time={new Date(thread.createdAt)}
        updated={new Date(thread.updatedAt)}
      />
    </Card>
  );
}
