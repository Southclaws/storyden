import { PostProps, ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { FeedItemMenu } from "../FeedItemMenu/FeedItemMenu";

import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";

type Props =
  | {
      kind: "thread";
      item: ThreadReference;
    }
  | {
      kind: "post";
      item: PostProps;
    };

export function PostRef({ kind, item }: Props) {
  const session = useSession();

  const data =
    kind === "thread"
      ? {
          id: item.id,
          title: item.title,
          permalink: `/t/${item.slug}`,
          short: item.short,
          author: item.author,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt,
        }
      : {
          id: item.id,
          title: item.root_slug, // TODO: Include parent thread title on API
          permalink: `/t/${item.root_slug}#${item.id}`,
          short: item.body,
          author: item.author,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt,
        };

  return (
    <Card
      id={data.id}
      title={data.title}
      text={data.short}
      url={data.permalink}
      shape="row"
      controls={
        session &&
        kind === "thread" && (
          <HStack>
            <CollectionMenu thread={item} />
            <FeedItemMenu thread={item} />
          </HStack>
        )
      }
    >
      <Byline
        href={data.permalink}
        author={data.author}
        time={new Date(data.createdAt)}
        updated={new Date(data.updatedAt)}
      />
    </Card>
  );
}
