import { PublicProfile } from "src/api/openapi/schemas";
import { Avatar } from "src/components/site/Avatar/Avatar";

import { Heading } from "@/components/ui/heading";
import { HStack, VStack } from "@/styled-system/jsx";

export function Header(props: PublicProfile) {
  return (
    <VStack className="profile__header" alignItems="start" minW="0" w="full">
      <HStack justifyContent="start" minW="0" w="full">
        <Avatar handle={props.handle} />

        <VStack
          alignItems="start"
          gap="2"
          overflow="hidden"
          minW="0"
          width="full"
          containerType="inline-size"
        >
          <Heading size="lg">{props.name}</Heading>
          <Heading
            w="full"
            size="md"
            color="gray.500"
            className="fluid-font-size"
            textWrap="nowrap"
            textOverflow="ellipsis"
            overflow="hidden"
          >
            @{props.handle}
          </Heading>
        </VStack>
      </HStack>
    </VStack>
  );
}
