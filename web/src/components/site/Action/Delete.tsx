import { TrashIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

import { ModalDrawer } from "../Modaldrawer/Modaldrawer";

import { HStack, VStack } from "@/styled-system/jsx";

export function DeleteAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  function handleConfirm() {
    onClose();
    props.onClick();
  }

  return (
    <>
      <Button kind="destructive" size="xs" onClick={onOpen}>
        <TrashIcon width="0.5em" height="0.5em" />
        {children}
      </Button>

      <ModalDrawer title="Delete" isOpen={isOpen} onClose={onClose}>
        <VStack w="full" gap="2" alignItems="start">
          <p>Are you sure?</p>

          <HStack w="full">
            <Button w="full" onClick={onClose}>
              Cancel
            </Button>
            <Button
              w="full"
              kind="destructive"
              onClick={handleConfirm}
              {...props}
            >
              Delete
            </Button>
          </HStack>
        </VStack>
      </ModalDrawer>
    </>
  );
}
