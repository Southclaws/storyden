import { Collection } from "src/api/openapi-schema";

import { Heading } from "@/components/ui/heading";
import { Box, Flex, LinkOverlay, styled } from "@/styled-system/jsx";

export function CollectionListItem(props: Collection) {
  return (
    <Box key={props.id} position="relative">
      <styled.li key={props.id} listStyleType="none">
        <Flex id={props.id} flexDir="column" gap="1">
          <LinkOverlay
            href={`/p/${props.owner.handle}/collections/${props.id}`}
          >
            <Heading size="md">{props.name}</Heading>
          </LinkOverlay>
          <p>{props.description}</p>
        </Flex>
      </styled.li>
    </Box>
  );
}
