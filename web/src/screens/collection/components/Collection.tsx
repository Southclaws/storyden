import { Flex, Heading, Text, VStack } from "@chakra-ui/react";
import { differenceInSeconds, formatDistanceToNow } from "date-fns";

import { CollectionWithItems } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";
import { formatDistanceDefaults } from "src/utils/date";

import { CollectionItemList } from "./CollectionItemList";

export function Collection(props: CollectionWithItems) {
  const createdAt = new Date(props.createdAt);
  const updatedAt = new Date(props.updatedAt);

  const created = formatDistanceToNow(createdAt, formatDistanceDefaults);
  const updated =
    differenceInSeconds(createdAt, updatedAt) > 0
      ? formatDistanceToNow(updatedAt, formatDistanceDefaults)
      : undefined;

  return (
    <VStack alignItems="start">
      <Heading size="md">{props.name}</Heading>
      <Flex alignItems="center">
        <Text as="span">{props.description}</Text>
        <Text color="gray.500" fontSize="sm">
          <Text as="span" px={2}>
            •
          </Text>
          <Text as="span">
            <Timestamp
              created={created}
              updated={updated}
              href={`/p/${props.owner.handle}/collections/${props.id}`}
            />
          </Text>
        </Text>
      </Flex>

      <CollectionItemList items={props.items} />
    </VStack>
  );
}
