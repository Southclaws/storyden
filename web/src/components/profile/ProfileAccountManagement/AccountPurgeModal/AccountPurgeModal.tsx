import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
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
  return (
    <ModalDrawer
      onOpen={onOpen}
      isOpen={isOpen}
      onClose={onClose}
      title={`Purge account content: ${handle}`}
    >
      <AccountPurgeScreen
        accountId={accountId}
        handle={handle}
        onSave={onClose}
      />
    </ModalDrawer>
  );
}
