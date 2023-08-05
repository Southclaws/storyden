import { Flex, Heading, LinkBox, LinkOverlay, Text } from "@chakra-ui/react";
import NextLink from "next/link";

import { CollectionItem } from "src/api/openapi/schemas";
import { Byline } from "src/screens/thread/components/Byline";

export function CollectionItem(props: { item: CollectionItem }) {
  const permalink = `/t/${props.item.slug}`;

  return (
    <Flex as="section" flexDir="column" py={2} width="full" gap={2}>
      <LinkBox as="article">
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay as={NextLink} href={permalink}>
              {props.item.title}
            </LinkOverlay>
          </Heading>
        </Flex>

        <Text noOfLines={3}>{props.item.short}</Text>
      </LinkBox>

      <Flex justifyContent="space-between">
        <Byline
          href={permalink}
          author={props.item.author.handle}
          time={new Date(props.item.createdAt)}
          updated={new Date(props.item.updatedAt)}
          //   more={<ThreadMenu {...props.item} />}
        />

        {/* Tags list */}
      </Flex>
    </Flex>
  );
}
