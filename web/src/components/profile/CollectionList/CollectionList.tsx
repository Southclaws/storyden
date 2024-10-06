import { Fragment } from "react";

import { Collection } from "src/api/openapi-schema";

import { Heading } from "@/components/ui/heading";
import { Box, Flex, LinkOverlay, styled } from "@/styled-system/jsx";
import { Divider, VStack } from "@/styled-system/jsx";

type Props = {
  collections: Collection[];
};
export function CollectionList(props: Props) {
  return (
    <VStack alignItems="start">
      {/* TODO: Actually design this lol
      // {isSelf && <CollectionCreateTrigger />} */}

      <styled.ol gap="4" display="flex" flexDir="column" width="full" m="0">
        {props.collections.map((c) => (
          <Fragment key={c.id}>
            <Divider />
            <CollectionListItem {...c} />
          </Fragment>
        ))}
      </styled.ol>
    </VStack>
  );
}

export function CollectionListItem(props: Collection) {
  return (
    <Box key={props.id} position="relative">
      <styled.li key={props.id} listStyleType="none">
        <Flex id={props.id} flexDir="column" gap="1">
          <LinkOverlay
            href={`/m/${props.owner.handle}/collections/${props.id}`}
          >
            <Heading size="md">{props.name}</Heading>
          </LinkOverlay>
          <p>{props.description}</p>
        </Flex>
      </styled.li>
    </Box>
  );
}
