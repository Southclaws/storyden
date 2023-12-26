import { PublicProfile } from "src/api/openapi/schemas";
import { Avatar } from "src/components/site/Avatar/Avatar";
import { Heading2, Heading3 } from "src/theme/components/Heading/Index";

import { HStack, VStack } from "@/styled-system/jsx";

export function Header(props: PublicProfile) {
  return (
    <VStack alignItems="start">
      <HStack justifyContent="start">
        <Avatar handle={props.handle} />

        <VStack alignItems="start" gap="2">
          <Heading2 size="lg">{props.name}</Heading2>
          <Heading3 size="md" color="gray.500">
            @{props.handle}
          </Heading3>
        </VStack>
      </HStack>
    </VStack>
  );
}
