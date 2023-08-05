import {
  Flex,
  Heading,
  LinkBox,
  LinkOverlay,
  ListItem,
  Text,
} from "@chakra-ui/react";
import Link from "next/link";

import { Collection } from "src/api/openapi/schemas";

export function CollectionListItem(props: Collection) {
  return (
    <LinkBox key={props.id} as="article">
      <ListItem key={props.id} listStyleType="none">
        <Flex id={props.id} flexDir="column" gap={1}>
          <LinkOverlay
            as={Link}
            href={`/p/${props.owner.handle}/collections/${props.id}`}
          >
            <Heading size="md">{props.name}</Heading>
          </LinkOverlay>
          <Text>{props.description}</Text>
        </Flex>
      </ListItem>
    </LinkBox>
  );
}
