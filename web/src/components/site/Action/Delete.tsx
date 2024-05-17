import { TrashIcon } from "@heroicons/react/24/outline";
import { MouseEvent, MouseEventHandler, PropsWithChildren } from "react";

import { Button, ButtonProps } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

import { ModalDrawer } from "../Modaldrawer/Modaldrawer";

import { HStack, VStack } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

type DeleteConfirmationProps = {
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
  handleConfirm: MouseEventHandler<HTMLButtonElement>;
};

type ComponentProps = ButtonProps &
  React.ButtonHTMLAttributes<HTMLButtonElement>;

export function useDeleteAction(props: {
  onClick?: MouseEventHandler<HTMLButtonElement>;
}): DeleteConfirmationProps {
  const { isOpen, onOpen, onClose } = useDisclosure();
  function handleConfirm(e: MouseEvent<HTMLButtonElement>) {
    onClose();
    props.onClick?.(e);
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
  onClick,
  ...props
}: PropsWithChildren<ComponentProps>) {
  const [bvp] = button.splitVariantProps(props);
  const deleteProps = useDeleteAction({
    onClick,
  });

  return (
    <>
      <Button
        colorPalette="red"
        size="xs"
        onClick={deleteProps.onOpen}
        {...bvp}
      >
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
          <Button w="full" colorPalette="red" onClick={handleConfirm}>
            Delete
          </Button>
        </HStack>
      </VStack>
    </ModalDrawer>
  );
}
