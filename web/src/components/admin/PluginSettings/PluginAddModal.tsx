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

      <PluginAddScreen isOpen={isOpen} onClose={onClose} />
    </>
  );
}
