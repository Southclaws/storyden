import { Button } from "@/components/ui/button";
import { HStack, VStack } from "@/styled-system/jsx";

import {
  Props,
  WithDisclosure,
  useDeleteDeviceScreen,
} from "./useDeleteDeviceScreen";

export function DeleteDeviceScreen(props: WithDisclosure<Props>) {
  const { handleConfirm } = useDeleteDeviceScreen(props);

  return (
    <VStack maxW="prose">
      <p>
        Warning: Deleting an authentication device is permanent. Make sure you
        have another authentication method or device registered to your account.
      </p>
      <HStack
        w="full"
        justifyContent="space-between"
        alignItems="center"
        justify="end"
        pb="3"
        gap="4"
      >
        <Button flexGrow="1" size="sm" variant="ghost" onClick={props.onClose}>
          Cancel
        </Button>
        <Button
          flexGrow="1"
          size="sm"
          colorPalette="red"
          onClick={handleConfirm}
        >
          Delete
        </Button>
      </HStack>
    </VStack>
  );
}
