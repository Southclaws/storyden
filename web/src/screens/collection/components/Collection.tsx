import { differenceInSeconds, formatDistanceToNow } from "date-fns";

import { CollectionWithItems } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";
import { Heading1 } from "src/theme/components/Heading/Index";
import { formatDistanceDefaults } from "src/utils/date";

import { Flex, VStack, styled } from "@/styled-system/jsx";

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
      <Heading1 size="md">{props.name}</Heading1>
      <Flex alignItems="center">
        <styled.p fontSize="sm">
          <styled.span>{props.description}</styled.span>

          <styled.span px="2">•</styled.span>

          <styled.span>
            <Timestamp
              created={created}
              updated={updated}
              href={`/p/${props.owner.handle}/collections/${props.id}`}
            />
          </styled.span>
        </styled.p>
      </Flex>

      <CollectionItemList items={props.items} />
    </VStack>
  );
}
