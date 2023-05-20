import { HStack, Text } from "@chakra-ui/react";
import { Anchor } from "../site/Anchor";

import { Avatar } from "../Avatar/Avatar";

type Props = {
  handle: string;
  showHandle?: boolean;
};

export function ProfileReference({ handle, showHandle = true }: Props) {
  return (
    <Anchor
      href={`/p/${handle}`}
      _hover={{ backgroundColor: "blackAlpha.100" }}
      p={1}
      pr={showHandle ? 2 : 1}
      borderRadius="full"
    >
      <HStack>
        <Avatar handle={handle} />
        {showHandle && <Text>@{handle}</Text>}
      </HStack>
    </Anchor>
  );
}
