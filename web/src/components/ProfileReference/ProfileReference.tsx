import { HStack, Text } from "@chakra-ui/react";

import { useSession } from "src/auth";

import { Avatar } from "../Avatar/Avatar";
import { Anchor } from "../site/Anchor";

type Props = {
  handle: string;
  showHandle?: boolean;
  size?: "sm" | "lg";
};

export function ProfileReference({
  handle,
  showHandle = true,
  size = "sm",
}: Props) {
  const account = useSession();
  const self = account?.handle === handle;
  const title = self ? `Your profile` : `${handle}'s profile`;
  const large = size === "lg";

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
        <Avatar handle={handle} width={large ? 8 : 6} />
        {showHandle && <Text fontSize={large ? "md" : "sm"}>@{handle}</Text>}
      </HStack>
    </Anchor>
  );
}
