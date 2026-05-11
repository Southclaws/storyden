import { formatDate } from "date-fns";

import {
  OwnedAccessKey,
  OwnedAccessKeyList,
  Permission,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { PermissionBadge } from "@/components/role/PermissionBadge";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { useI18n } from "@/i18n/provider";
import { PermissionDetails } from "@/lib/permission/permission";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox, lstack } from "@/styled-system/patterns";

import { useAdminAccessKeyList } from "./useAdminAccessKeyList";

type Props = {
  keys: OwnedAccessKeyList;
};

export function AccessKeySettings({ keys }: Props) {
  const { t } = useI18n();
  const { revokeKey } = useAdminAccessKeyList();

  if (keys.length === 0) {
    return (
      <CardBox className={lstack()}>
        <Heading size="md">{t("Access keys")}</Heading>
        <p>{t("No access keys have been created yet.")}</p>
      </CardBox>
    );
  }

  const totalKeys = keys.length;
  const totalActiveKeys = keys.filter(isKeyActive).length;
  const hasInactive = totalKeys != totalActiveKeys;

  return (
    <CardBox className={lstack()}>
      <Heading size="md">{t("Access keys")}</Heading>
      <p>{t("All access keys created by members of this site.")}</p>
      <div>
        <span>
          <WarningIcon display="inline" width="4" />{" "}
          <strong>{t("Note:")}</strong>{" "}
          {t(
            "If you revoke the ability to use access keys from a role or member (by removing the",
          )}{" "}
          {<PermissionBadge permission={Permission.USE_PERSONAL_ACCESS_KEYS} />}{" "}
          {t("permission), this will not revoke their existing access keys.")}
        </span>
      </div>

      <WStack alignItems="center" color="fg.muted">
        {hasInactive ? (
          <styled.p>
            {totalKeys} {t("access keys")}, {totalActiveKeys} {t("active")}.
          </styled.p>
        ) : (
          <styled.p>
            {keys.length} {t("access keys")}.
          </styled.p>
        )}
      </WStack>

      <ul className={lstack({ gap: "3" })}>
        {keys.map((key) => (
          <AccessKeyItem
            key={key.id}
            accessKey={key}
            onRevoke={() => revokeKey(key.id)}
          />
        ))}
      </ul>
    </CardBox>
  );
}

type AccessKeyItemProps = {
  accessKey: OwnedAccessKeyList[number];
  onRevoke: () => Promise<void>;
};

function AccessKeyItem({ accessKey, onRevoke }: AccessKeyItemProps) {
  const { t } = useI18n();
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
            <Badge>{t(inactiveStatus)}</Badge>
          )}
        </WStack>

        <WStack flexWrap="wrap">
          <styled.p fontSize="xs">
            {t("Created")}:{" "}
            <time>{formatDate(accessKey.createdAt, "PPpp")}</time>
          </styled.p>

          {accessKey.expires_at && (
            <Badge gap="1">
              <span>{t("Expiry")}:</span>
              <time>{formatDate(accessKey.expires_at, "PPpp")}</time>
            </Badge>
          )}
        </WStack>

        <MemberBadge
          size="sm"
          name="full-horizontal"
          profile={accessKey.created_by}
        />
      </LStack>
    </li>
  );
}

function isKeyActive(key: OwnedAccessKey) {
  return key.enabled && !isKeyExpired(key);
}

function isKeyExpired(key: OwnedAccessKey) {
  if (key.expires_at === undefined) {
    return false;
  }

  const expiryDate = new Date(key.expires_at);

  if (expiryDate > new Date()) {
    return false;
  }

  return true;
}
