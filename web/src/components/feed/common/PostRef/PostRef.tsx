import {
  Post,
  PostProps,
  PostReference,
  ThreadReference,
} from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";

import { FeedItemMenu } from "../FeedItemMenu/FeedItemMenu";

type Props = {
  item: PostReference;
};

export function PostRef({ item }: Props) {
  const session = useSession();

  const permalink = `/t/${item.slug}#${item.id}`;

  return (
    <Card
      id={item.id}
      title={item.title}
      text={item.description}
      url={permalink}
      shape="row"
      controls={
        session && (
          <HStack>
            <CollectionMenu thread={item} />
            <FeedItemMenu thread={item} />
          </HStack>
        )
      }
    >
      <Byline
        href={permalink}
        author={item.author}
        time={new Date(item.createdAt)}
        updated={new Date(item.updatedAt)}
      />
    </Card>
  );
}
