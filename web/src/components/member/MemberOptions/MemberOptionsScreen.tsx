import { useSession } from "src/auth";
import { ProfileLayout } from "src/screens/profile/ProfileLayout";

import { MemberSuspensionTrigger } from "../MemberSuspension/MemberSuspensionTrigger";

import { HStack, VStack } from "@/styled-system/jsx";

import { Props } from "./useMemberOptionsScreen";

export function MemberMenuOptionsScreen(props: Props) {
  const session = useSession();

  const showAdminOptions = session?.admin && props.handle !== session.handle;

  return (
    <VStack height="full" justify="space-between">
      <ProfileLayout {...props} />

      {showAdminOptions && (
        <HStack w="full">
          <MemberSuspensionTrigger {...props} />
        </HStack>
      )}
    </VStack>
  );
}
