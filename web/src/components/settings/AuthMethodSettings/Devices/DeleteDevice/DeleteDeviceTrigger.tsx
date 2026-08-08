import { PropsWithChildren } from "react";

import { Button } from "@/components/ui/button";
import { useDisclosure } from "@/utils/useDisclosure";

import { DeleteDeviceModal } from "./DeleteDeviceModal";
import { Props } from "./useDeleteDeviceScreen";

export function DeleteDeviceTrigger(props: PropsWithChildren<Props>) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <Button colorPalette="red" onClick={onOpen}>
        {props.children ?? "Delete"}
      </Button>
      <DeleteDeviceModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
