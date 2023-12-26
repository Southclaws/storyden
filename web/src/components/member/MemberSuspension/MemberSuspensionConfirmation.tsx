import { Button } from "src/theme/components/Button";
import { WithDisclosure } from "src/utils/useDisclosure";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { Props, useMemberSuspension } from "./useMemberSuspension";

export function MemberSuspensionConfirmation(props: WithDisclosure<Props>) {
  const { handlers } = useMemberSuspension(props);

  return (
    <VStack alignItems="start">
      {props.deletedAt ? (
        <styled.p>
          Do you want to reinstate the suspended account {props.name}?
        </styled.p>
      ) : (
        <styled.p>Do you want to suspend the account {props.name}?</styled.p>
      )}

      <HStack w="full">
        <Button w="full">Cancel</Button>

        {props.deletedAt ? (
          <Button
            w="full"
            kind="destructive"
            onClick={handlers.handleReinstate}
          >
            Reinstate
          </Button>
        ) : (
          <Button
            w="full"
            kind="destructive"
            onClick={handlers.handleSuspension}
          >
            Suspend
          </Button>
        )}
      </HStack>
    </VStack>
  );
}
