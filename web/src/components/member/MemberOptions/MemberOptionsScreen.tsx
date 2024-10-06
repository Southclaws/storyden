import { useSession } from "src/auth";

import { HStack, VStack } from "@/styled-system/jsx";

import { MemberSuspensionTrigger } from "../MemberSuspension/MemberSuspensionTrigger";

import { Props } from "./useMemberOptionsScreen";

export function MemberMenuOptionsScreen(props: Props) {
  const session = useSession();

  const showAdminOptions = session?.admin && props.handle !== session.handle;

  return (
    <VStack height="full" justify="space-between">
      {showAdminOptions && (
        <HStack w="full">
          <MemberSuspensionTrigger {...props} />
        </HStack>
      )}
    </VStack>
  );
}
