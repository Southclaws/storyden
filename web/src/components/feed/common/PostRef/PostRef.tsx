import NextLink from "next/link";

import { PostProps, ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Heading, LinkBox, LinkOverlay } from "src/theme/components";

import { FeedItem } from "../FeedItem/FeedItem";

import { Flex, HStack, styled } from "@/styled-system/jsx";

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
          title: item.title,
          permalink: `/t/${item.slug}`,
          short: item.short,
          author: item.author,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt,
        }
      : {
          title: item.root_slug, // TODO: Include parent thread title on API
          permalink: `/t/${item.root_slug}#${item.id}`,
          short: item.body,
          author: item.author,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt,
        };

  return (
    <LinkBox>
      <FeedItem>
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay as={NextLink} href={data.permalink}>
              {data.title}
            </LinkOverlay>
          </Heading>
        </Flex>

        <styled.p lineClamp={3}>{data.short}</styled.p>

        <Flex justifyContent="space-between">
          <Byline
            href={data.permalink}
            author={data.author}
            time={new Date(data.createdAt)}
            updated={new Date(data.updatedAt)}
          />

          <HStack>
            {session && kind === "thread" && <CollectionMenu thread={item} />}
          </HStack>
        </Flex>
      </FeedItem>
    </LinkBox>
  );
}
