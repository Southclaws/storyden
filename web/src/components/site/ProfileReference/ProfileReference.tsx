import { Box, HStack } from "@/styled-system/jsx";

import { useSession } from "src/auth";
import { ProfileReference } from "src/api/openapi/schemas";
import { Avatar } from "src/components/site/Avatar/Avatar";
import { Anchor } from "src/components/site/Anchor";

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
      className="profile-reference"
      p={1}
      pr={showHandle ? 2 : 1}
      borderRadius="full"
      _hover={{ backgroundColor: "blackAlpha.100" }}
      href={`/p/${profileReference.handle}`}
      title={title}
      minW={0}
    >
      <HStack>
        <Avatar flexShrink={0} handle={profileReference.handle} width={large ? 8 : 6} />
        {showHandle && <Box minW={0} flexShrink={1}>
          <Handle profileReference={profileReference} size={size} />
        </Box>}
      </HStack>
    </Anchor>
  );
}
