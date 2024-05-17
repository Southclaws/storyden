import { PropsWithChildren } from "react";

import { useDisclosure } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";

import { DeleteDeviceModal } from "./DeleteDeviceModal";
import { Props } from "./useDeleteDeviceScreen";

export function DeleteDeviceTrigger(props: PropsWithChildren<Props>) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <Button size="xs" colorPalette="red" onClick={onOpen}>
        {props.children ?? "Delete"}
      </Button>
      <DeleteDeviceModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
