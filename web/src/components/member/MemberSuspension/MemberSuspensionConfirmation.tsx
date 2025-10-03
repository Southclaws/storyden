import { WithDisclosure } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";
import { HStack, VStack, styled } from "@/styled-system/jsx";

import { Props, useMemberSuspension } from "./useMemberSuspension";

export function MemberSuspensionConfirmation(props: WithDisclosure<Props>) {
  const { handlers } = useMemberSuspension(props);

  return (
    <VStack alignItems="start">
      {props.profile.suspended ? (
        <styled.p>
          Do you want to reinstate the suspended account {props.profile.name}?
        </styled.p>
      ) : (
        <styled.p>
          Do you want to suspend the account {props.profile.name}?
        </styled.p>
      )}

      <HStack w="full">
        <Button type="button" flexGrow="1" onClick={props.onClose}>
          Cancel
        </Button>

        {props.profile.suspended ? (
          <Button
            // w="full"
            flexGrow="1"
            colorPalette="red"
            onClick={handlers.handleReinstate}
          >
            Reinstate
          </Button>
        ) : (
          <Button
            // w="full"
            flexGrow="1"
            colorPalette="red"
            onClick={handlers.handleSuspension}
          >
            Suspend
          </Button>
        )}
      </HStack>
    </VStack>
  );
}
