import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { useI18n } from "@/i18n/provider";
import { UseDisclosureProps } from "@/utils/useDisclosure";

import { AccountPurgeScreen } from "./AccountPurgeScreen";
import { Props } from "./useAccountPurge";

export function AccountPurgeModal({
  accountId,
  handle,
  onClose,
  onOpen,
  isOpen,
}: UseDisclosureProps & Props) {
  const { t } = useI18n();

  return (
    <ModalDrawer
      onOpen={onOpen}
      isOpen={isOpen}
      onClose={onClose}
      title={t("Purge account content: {{handle}}", { handle })}
    >
      <AccountPurgeScreen
        accountId={accountId}
        handle={handle}
        onSave={onClose}
      />
    </ModalDrawer>
  );
}
