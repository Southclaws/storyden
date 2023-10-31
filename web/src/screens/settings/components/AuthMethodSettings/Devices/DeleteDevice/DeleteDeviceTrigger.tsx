import { useDisclosure } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { Button } from "src/theme/components/Button";

import { DeleteDeviceModal } from "./DeleteDeviceModal";
import { Props } from "./useDeleteDeviceScreen";

export function DeleteDeviceTrigger(props: PropsWithChildren<Props>) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <Button size="xs" kind="destructive" onClick={onOpen}>
        {props.children ?? "Delete"}
      </Button>
      <DeleteDeviceModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
