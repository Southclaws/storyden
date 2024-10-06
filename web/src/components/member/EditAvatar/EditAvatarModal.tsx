import { PropsWithChildren } from "react";

import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { Slot } from "@/components/ui/Slot";
import { Button, ButtonProps } from "@/components/ui/button";
import { useDisclosure } from "@/utils/useDisclosure";

import { EditAvatarScreen } from "./EditAvatarScreen";
import { Props } from "./useEditAvatar";

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

  const Trigger = asChild
    ? Slot
    : (bp: ButtonProps) => <Button {...bp}>Edit</Button>;

  return (
    <>
      <Trigger {...props} onClick={onOpen} />

      <EditAvatarModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
