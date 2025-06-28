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
import { PermissionDetails } from "@/lib/permission/permission";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox, lstack } from "@/styled-system/patterns";

import { useAdminAccessKeyList } from "./useAdminAccessKeyList";

type Props = {
  keys: OwnedAccessKeyList;
};

export function AccessKeySettings({ keys }: Props) {
  const { revokeKey } = useAdminAccessKeyList();

  if (keys.length === 0) {
    return (
      <CardBox className={lstack()}>
        <Heading size="md">Access keys</Heading>
        <p>No access keys have been created yet.</p>
      </CardBox>
    );
  }

  const totalKeys = keys.length;
  const totalActiveKeys = keys.filter(isKeyActive).length;
  const hasInactive = totalKeys != totalActiveKeys;

  return (
    <CardBox className={lstack()}>
      <Heading size="md">Access keys</Heading>
      <p>All access keys created by members of this site.</p>
      <p>
        <WarningIcon display="inline" /> <strong>Note:</strong> if you revoke
        the ability to use access keys from a role or member (by removing the{" "}
        {<PermissionBadge permission={Permission.USE_PERSONAL_ACCESS_KEYS} />}{" "}
        permission), this will not revoke their existing access keys.
      </p>

      <WStack alignItems="center" color="fg.muted">
        {hasInactive ? (
          <styled.p>
            {totalKeys} access keys, {totalActiveKeys} active.
          </styled.p>
        ) : (
          <styled.p>{keys.length} access keys.</styled.p>
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
