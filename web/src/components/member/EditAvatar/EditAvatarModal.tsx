import { PropsWithChildren } from "react";

import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button, ButtonProps } from "@/components/ui/button";
import { Slot } from "@/components/ui/slot";
import { useDisclosure } from "@/utils/useDisclosure";

import { EditAvatarScreen } from "./EditAvatarScreen";
import { Props } from "./useEditAvatar";

function EditAvatarButton(props: ButtonProps) {
  return <Button {...props}>Edit</Button>;
}

export function EditAvatarModal(props: Props) {
  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title="Edit avatar"
      dismissable={false}
    >
      <EditAvatarScreen {...props} />
    </ModalDrawer>
  );
}

export function EditAvatarTrigger({
  asChild,
  ...props
}: PropsWithChildren<Props & { asChild?: boolean }>) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const Trigger = asChild ? Slot : EditAvatarButton;

  return (
    <>
      <Trigger {...props} onClick={onOpen} />

      <EditAvatarModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
