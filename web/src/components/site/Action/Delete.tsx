import { TrashIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

import { ModalDrawer } from "../Modaldrawer/Modaldrawer";

import { HStack, VStack } from "@/styled-system/jsx";

type DeleteConfirmationProps = {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
  handleConfirm: () => void;
};

export function useDeleteAction(props: {
  onClick: () => void;
}): DeleteConfirmationProps {
  const { isOpen, onOpen, onClose } = useDisclosure();
  function handleConfirm() {
    onClose();
    props.onClick();
  }

  return {
    onOpen,
    onClose,
    isOpen,
    handleConfirm,
  };
}

export function DeleteAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  const deleteProps = useDeleteAction(props);

  return (
    <>
      <Button kind="destructive" size="xs" onClick={deleteProps.onOpen}>
        <TrashIcon width="0.5em" height="0.5em" />
        {children}
      </Button>

      <DeleteConfirmation {...deleteProps} />
    </>
  );
}

export function DeleteConfirmation({
  isOpen,
  onClose,
  handleConfirm,
}: DeleteConfirmationProps) {
  return (
    <ModalDrawer title="Delete" isOpen={isOpen} onClose={onClose}>
      <VStack w="full" gap="2" alignItems="start">
        <p>Are you sure?</p>

        <HStack w="full">
          <Button w="full" onClick={onClose}>
            Cancel
          </Button>
          <Button w="full" kind="destructive" onClick={handleConfirm}>
            Delete
          </Button>
        </HStack>
      </VStack>
    </ModalDrawer>
  );
}
