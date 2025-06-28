import { formatDate } from "date-fns";

import { AccessKey, AccessKeyList } from "@/api/openapi-schema";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { AddIcon } from "@/components/ui/icons/Add";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox, lstack } from "@/styled-system/patterns";
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
      <CardBox className={lstack()} gap="8">
        <LStack>
          <Heading size="md">Access keys</Heading>

          <p>
            Access keys allow you to authenticate API requests. They share the
            same permissions as your account. If your account receives new
            roles, your access keys will inherit the permissions assigned to
            those roles.
          </p>
        </LStack>

        <LStack>
          <WStack alignItems="center" color="fg.muted">
            {hasInactive ? (
              <styled.p>
                {totalKeys} access keys, {totalActiveKeys} active.
              </styled.p>
            ) : (
              <styled.p>{keys.length} access keys.</styled.p>
            )}
            <Button size="xs" variant="subtle" onClick={createModal.onOpen}>
              <AddIcon />
              New
            </Button>
          </WStack>

          <AccessKeyItemList keys={keys} />
        </LStack>
      </CardBox>

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
          <Heading size="sm">{accessKey.name}</Heading>

          {inactiveStatus === undefined ? (
            <HStack>
              {isConfirming ? (
                <>
                  <Button
                    size="xs"
                    variant="subtle"
                    bgColor="bg.destructive"
                    onClick={handleConfirmAction}
                  >
                    Confirm Revoke
                  </Button>
                  <Button
                    size="xs"
                    variant="outline"
                    onClick={handleCancelAction}
                  >
                    Cancel
                  </Button>
                </>
              ) : (
                <Button
                  size="xs"
                  variant="outline"
                  bgColor="bg.destructive"
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
          <styled.p fontSize="xs">
            Created: <time>{formatDate(accessKey.createdAt, "PPpp")}</time>
          </styled.p>

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
