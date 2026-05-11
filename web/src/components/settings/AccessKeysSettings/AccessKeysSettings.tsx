import { formatDate } from "date-fns";

import { AccessKey, AccessKeyList } from "@/api/openapi-schema";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { AddIcon } from "@/components/ui/icons/Add";
import { useI18n } from "@/i18n/provider";
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
  const { t } = useI18n();

  const totalKeys = keys.length;
  const totalActiveKeys = keys.filter(isKeyActive).length;
  const hasInactive = totalKeys != totalActiveKeys;

  return (
    <>
      <CardBox className={lstack()} gap="8">
        <LStack>
          <Heading size="md">{t("Access keys")}</Heading>

          <p>
            {t(
              "Access keys allow you to authenticate API requests. They share the same permissions as your account. If your account receives new roles, your access keys will inherit the permissions assigned to those roles.",
            )}
          </p>
        </LStack>

        <LStack>
          <WStack alignItems="center" color="fg.muted">
            {hasInactive ? (
              <styled.p>
                {t("{{total}} access keys, {{active}} active.", {
                  total: totalKeys,
                  active: totalActiveKeys,
                })}
              </styled.p>
            ) : (
              <styled.p>
                {t("{{count}} access keys.", { count: keys.length })}
              </styled.p>
            )}
            <Button size="xs" variant="subtle" onClick={createModal.onOpen}>
              <AddIcon />
              {t("New")}
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
  const { t } = useI18n();

  if (keys.length === 0) {
    return (
      <p style={{ color: "var(--colors-gray-500)", fontStyle: "italic" }}>
        {t("No access keys created yet.")}
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
  const { t } = useI18n();

  const isExpired = isKeyExpired(accessKey);

  const inactiveStatus = isExpired
    ? t("Expired")
    : accessKey.enabled
      ? undefined
      : t("Revoked");

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
                    {t("Confirm Revoke")}
                  </Button>
                  <Button
                    size="xs"
                    variant="outline"
                    onClick={handleCancelAction}
                  >
                    {t("Cancel")}
                  </Button>
                </>
              ) : (
                <Button
                  size="xs"
                  variant="outline"
                  bgColor="bg.destructive"
                  onClick={handleConfirmAction}
                >
                  {t("Revoke")}
                </Button>
              )}
            </HStack>
          ) : (
            <Badge>{inactiveStatus}</Badge>
          )}
        </WStack>

        <WStack flexWrap="wrap">
          <styled.p fontSize="xs">
            {t("Created")}: <time>{formatDate(accessKey.createdAt, "PPpp")}</time>
          </styled.p>

          {accessKey.expires_at && (
            <Badge gap="1">
              <span>{t("Expiry")}:</span>
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
