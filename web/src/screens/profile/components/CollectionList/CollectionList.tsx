import {
  Box,
  Divider,
  Flex,
  Heading,
  LinkBox,
  LinkOverlay,
  ListItem,
  OrderedList,
  VStack,
} from "@chakra-ui/react";
import Link from "next/link";

import { Collection } from "src/api/openapi/schemas";
import { ContentViewer } from "src/components/ContentViewer/ContentViewer";

import { CollectionCreate } from "./CollectionCreate/CollectionCreate";

type Props = {
  collections: Collection[];
};
export function CollectionList(props: Props) {
  return (
    <VStack alignItems="end">
      <CollectionCreate />

      <OrderedList gap={4} display="flex" flexDir="column" width="full" m={0}>
        {props.collections.map((c) => (
          <LinkBox key={c.id} as="article">
            <Divider />
            <ListItem key={c.id} listStyleType="none">
              <Flex id={c.id} flexDir="column" gap={1}>
                <LinkOverlay
                  as={Link}
                  href={`/p/${c.owner.handle}/collections/${c.id}`}
                >
                  <Heading size="md">{c.name}</Heading>
                </LinkOverlay>
                <ContentViewer value={c.description} />
              </Flex>
            </ListItem>
          </LinkBox>
        ))}
      </OrderedList>
    </VStack>
  );
}
