import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import { AddIcon } from "@/components/ui/icons/Add";
import { UseDisclosureProps, useDisclosure } from "@/utils/useDisclosure";

import { PluginAddScreen } from "./PluginAddScreen";

export function PluginAddTrigger() {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <>
      <Button type="button" size="xs" variant="ghost" onClick={onOpen}>
        <AddIcon />
        <span>Add plugin</span>
      </Button>

      <PluginAddModal isOpen={isOpen} onClose={onClose} />
    </>
  );
}

function PluginAddModal({ isOpen, onClose }: UseDisclosureProps) {
  return (
    <ModalDrawer
      isOpen={isOpen}
      onClose={onClose}
      title="Add Plugin"
      // TODO: Do this via Context
      // dismissable={!isUploading}
    >
      <PluginAddScreen onClose={onClose} />
    </ModalDrawer>
  );
}
