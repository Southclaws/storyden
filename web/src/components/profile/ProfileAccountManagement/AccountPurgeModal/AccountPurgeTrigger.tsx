import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { hasPermission } from "@/utils/permissions";
import { useDisclosure } from "@/utils/useDisclosure";

import { AccountPurgeModal } from "./AccountPurgeModal";
import { Props } from "./useAccountPurge";

export function AccountPurgeTrigger({ accountId, handle }: Props) {
  const session = useSession();
  const { onOpen, isOpen, onClose } = useDisclosure();

  if (!hasPermission(session, Permission.ADMINISTRATOR)) {
    return null;
  }

  return (
    <>
      <Button size="xs" variant="outline" colorPalette="red" onClick={onOpen}>
        <WarningIcon />
        Purge Content
      </Button>

      <AccountPurgeModal
        accountId={accountId}
        handle={handle}
        isOpen={isOpen}
        onClose={onClose}
        onOpen={onOpen}
      />
    </>
  );
}
