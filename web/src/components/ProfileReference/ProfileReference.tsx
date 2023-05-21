import { HStack, Text } from "@chakra-ui/react";
import { Anchor } from "../site/Anchor";

import { Avatar } from "../Avatar/Avatar";
import { useSession } from "src/auth";

type Props = {
  handle: string;
  showHandle?: boolean;
};

export function ProfileReference({ handle, showHandle = true }: Props) {
  const account = useSession();
  const self = account?.handle === handle;
  const title = self ? `Your profile` : `${handle}'s profile`;

  return (
    <Anchor
      p={1}
      pr={showHandle ? 2 : 1}
      borderRadius="full"
      _hover={{ backgroundColor: "blackAlpha.100" }}
      href={`/p/${handle}`}
      title={title}
    >
      <HStack>
        <Avatar handle={handle} />
        {showHandle && <Text>@{handle}</Text>}
      </HStack>
    </Anchor>
  );
}
