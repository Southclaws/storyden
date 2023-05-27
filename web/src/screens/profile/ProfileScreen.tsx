import { Unready } from "src/components/Unready";
import { Props, useProfileScreen } from "./useProfileScreen";
import { Profile } from "./components/Profile";
import { Flex } from "@chakra-ui/react";

export function ProfileScreen(props: Props) {
  const profile = useProfileScreen(props);

  if (!profile.ready) return <Unready {...profile.error} />;

  return (
    <Flex py={4}>
      <Profile {...profile.data} />
    </Flex>
  );
}
