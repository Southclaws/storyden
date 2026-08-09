import { formatDate } from "date-fns";

import { AccessKey, AccessKeyList } from "@/api/openapi-schema";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { AddIcon } from "@/components/ui/icons/Add";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { cardBox } from "@/styled-system/recipes";
import { useDisclosure } from "@/utils/useDisclosure";

import { CreateAccessKeyModal } from "./CreateAccessKeyModal";
import { useAccessKeySettings } from "./useAccessKeySettings";

type Props = {
  keys: AccessKeyList;
};

export function AccessKeysSettings({ keys }: Props) {
  const createModal = useDisclosure();

  const totalKeys = keys.length;
  const totalActiveKeys = keys.filter(isKeyActive).length;
  const hasInactive = totalKeys != totalActiveKeys;

  return (
    <>
      <LStack gap="4">
        <LStack gap="1">
          <PageHeading>Access keys</PageHeading>
          <Text variant="supporting">
            Access keys allow you to authenticate API requests. They share the
            same permissions as your account. If your account receives new
            roles, your access keys will inherit the permissions assigned to
            those roles.
          </Text>
        </LStack>

        <LStack>
          <WStack alignItems="center" color="text.subtle">
            {hasInactive ? (
              <Text variant="metadata">
                {totalKeys} access keys, {totalActiveKeys} active.
              </Text>
            ) : (
              <Text variant="metadata">{keys.length} access keys.</Text>
            )}
            <Button variant="subtle" onClick={createModal.onOpen}>
              <AddIcon />
              New
            </Button>
          </WStack>

          <AccessKeyItemList keys={keys} />
        </LStack>
      </LStack>

      <CreateAccessKeyModal
        isOpen={createModal.isOpen}
        onClose={createModal.onClose}
      />
    </>
  );
}

function AccessKeyItemList({ keys }: Props) {
  const { revokeKey } = useAccessKeySettings();

  if (keys.length === 0) {
    return (
      <p style={{ color: "var(--colors-gray-500)", fontStyle: "italic" }}>
        No access keys created yet.
      </p>
    );
  }

  return (
    <ul className={lstack({ gap: "3" })}>
      {keys.map((key) => (
        <AccessKeyItem
          key={key.id}
          accessKey={key}
          onRevoke={() => revokeKey(key.id)}
        />
      ))}
    </ul>
  );
}

type AccessKeyItemProps = {
  accessKey: AccessKeyList[number];
  onRevoke: () => Promise<void>;
};

function AccessKeyItem({ accessKey, onRevoke }: AccessKeyItemProps) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(onRevoke);

  const isExpired = isKeyExpired(accessKey);

  const inactiveStatus = isExpired
    ? "Expired"
    : accessKey.enabled
      ? undefined
      : "Revoked";

  return (
    <li className={cardBox()}>
      <LStack>
        <WStack>
          <Text variant="supporting" color="text.default" fontWeight="semibold">
            {accessKey.name}
          </Text>

          {inactiveStatus === undefined ? (
            <HStack>
              {isConfirming ? (
                <>
                  <Button
                    variant="subtle"
                    bgColor="status.danger.surface"
                    onClick={handleConfirmAction}
                  >
                    Confirm Revoke
                  </Button>
                  <Button variant="outline" onClick={handleCancelAction}>
                    Cancel
                  </Button>
                </>
              ) : (
                <Button
                  variant="outline"
                  bgColor="status.danger.surface"
                  onClick={handleConfirmAction}
                >
                  Revoke
                </Button>
              )}
            </HStack>
          ) : (
            <Badge>{inactiveStatus}</Badge>
          )}
        </WStack>

        <WStack flexWrap="wrap">
          <Text variant="metadata">
            Created: <time>{formatDate(accessKey.createdAt, "PPpp")}</time>
          </Text>

          {accessKey.expires_at && (
            <Badge gap="1">
              <span>Expiry:</span>
              <time>{formatDate(accessKey.expires_at, "PPpp")}</time>
            </Badge>
          )}
        </WStack>
      </LStack>
    </li>
  );
}

function isKeyActive(key: AccessKey) {
  return key.enabled && !isKeyExpired(key);
}

function isKeyExpired(key: AccessKey) {
  if (key.expires_at === undefined) {
    return false;
  }

  const expiryDate = new Date(key.expires_at);

  if (expiryDate > new Date()) {
    return false;
  }

  return true;
}
