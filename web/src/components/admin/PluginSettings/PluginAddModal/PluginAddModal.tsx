import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import { AddIcon } from "@/components/ui/icons/Add";
import { useI18n } from "@/i18n/provider";
import { UseDisclosureProps, useDisclosure } from "@/utils/useDisclosure";

import { PluginAddScreen } from "./PluginAddScreen";

export function PluginAddTrigger() {
  const { t } = useI18n();
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <>
      <Button type="button" size="xs" variant="ghost" onClick={onOpen}>
        <AddIcon />
        <span>{t("Add plugin")}</span>
      </Button>

      <PluginAddModal isOpen={isOpen} onClose={onClose} />
    </>
  );
}

function PluginAddModal({ isOpen, onClose }: UseDisclosureProps) {
  const { t } = useI18n();
  return (
    <ModalDrawer
      isOpen={isOpen}
      onClose={onClose}
      title={t("Add Plugin")}
      // TODO: Do this via Context
      // dismissable={!isUploading}
    >
      <PluginAddScreen onClose={onClose} />
    </ModalDrawer>
  );
}
