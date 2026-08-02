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
import { WarningIcon } from "@/components/ui/icons/Warning";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { PermissionDetails } from "@/lib/permission/permission";
import { HStack, LStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { cardBox } from "@/styled-system/recipes";

import { useAdminAccessKeyList } from "./useAdminAccessKeyList";

type Props = {
  keys: OwnedAccessKeyList;
};

export function AccessKeySettings({ keys }: Props) {
  const { revokeKey } = useAdminAccessKeyList();

  if (keys.length === 0) {
    return (
      <LStack>
        <PageHeading>Access keys</PageHeading>
        <p>No access keys have been created yet.</p>
      </LStack>
    );
  }

  const totalKeys = keys.length;
  const totalActiveKeys = keys.filter(isKeyActive).length;
  const hasInactive = totalKeys != totalActiveKeys;

  return (
    <LStack>
      <PageHeading>Access keys</PageHeading>
      <p>All access keys created by members of this site.</p>
      <div>
        <span>
          <WarningIcon display="inline" width="4" /> <strong>Note:</strong> if
          you revoke the ability to use access keys from a role or member (by
          removing the{" "}
          {<PermissionBadge permission={Permission.USE_PERSONAL_ACCESS_KEYS} />}{" "}
          permission), this will not revoke their existing access keys.
        </span>
      </div>

      <WStack alignItems="center" color="text.subtle">
        {hasInactive ? (
          <Text variant="metadata">
            {totalKeys} access keys, {totalActiveKeys} active.
          </Text>
        ) : (
          <Text variant="metadata">{keys.length} access keys.</Text>
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
    </LStack>
  );
}

type AccessKeyItemProps = {
  accessKey: OwnedAccessKeyList[number];
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
