import { HStack } from "@chakra-ui/react";

import { useSession } from "src/auth";

import { ProfileReference } from "src/api/openapi/schemas";
import { Avatar } from "../site/Avatar/Avatar";
import { Anchor } from "../site/Anchor";
import { Handle } from "./Handle";

type Props =  {
  profileReference: ProfileReference;
  showHandle?: boolean;
  size?: "sm" | "lg";
};

export function ProfileReference({
  profileReference,
  showHandle = true,
  size = "sm",
}: Props) {
  const account = useSession();
  const self = account?.id === profileReference.id;
  const title = self ? `Your profile` : `${profileReference.handle}'s profile`;
  const large = size === "lg";

  return (
    <Anchor
      p={1}
      pr={showHandle ? 2 : 1}
      borderRadius="full"
      _hover={{ backgroundColor: "blackAlpha.100" }}
      href={`/p/${profileReference.handle}`}
      title={title}
    >
      <HStack>
        <Avatar handle={profileReference.handle} width={large ? 8 : 6} />
        {showHandle && <Handle profileReference={profileReference} size={size} />}
      </HStack>
    </Anchor>
  );
}
