import { CollectionWithItems } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";

import { Heading1 } from "@/components/ui/typography-heading";
import { Flex, VStack, styled } from "@/styled-system/jsx";

import { CollectionItemList } from "./CollectionItemList";

export function Collection(props: CollectionWithItems) {
  return (
    <VStack alignItems="start">
      <Heading1 size="md">{props.name}</Heading1>
      <Flex alignItems="center">
        <styled.p fontSize="sm">
          <styled.span>{props.description}</styled.span>

          <styled.span px="2">•</styled.span>

          <styled.span>
            <Timestamp
              created={props.createdAt}
              updated={props.updatedAt}
              href={`/p/${props.owner.handle}/collections/${props.id}`}
            />
          </styled.span>
        </styled.p>
      </Flex>

      <CollectionItemList items={props.items} />
    </VStack>
  );
}
