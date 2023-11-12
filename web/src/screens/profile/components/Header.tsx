import { HStack, Heading, Text, VStack } from "@chakra-ui/react";

import { PublicProfile } from "src/api/openapi/schemas";
import { Avatar } from "src/components/site/Avatar/Avatar";

export function Header(props: PublicProfile) {
  return (
    <VStack alignItems="start">
      <HStack justifyContent="start">
        <Avatar handle={props.handle} />

        <VStack alignItems="start" spacing={1}>
          <Heading>{props.name}</Heading>
          <Text as="h3" size="md" color="gray.500">
            @{props.handle}
          </Text>
        </VStack>
      </HStack>
    </VStack>
  );
}
