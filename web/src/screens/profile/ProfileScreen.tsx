import { Unready } from "src/components/Unready";
import { useProfileScreen } from "./useProfileScreen";
import { Profile } from "./components/Profile";
import { Flex } from "@chakra-ui/react";

export function ProfileScreen() {
  const profile = useProfileScreen();

  if (!profile.ready) return <Unready {...profile.error} />;

  return (
    <Flex py={4}>
      <Profile {...profile.data} />
    </Flex>
  );
}
