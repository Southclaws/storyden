"use client";

import { formatDate } from "date-fns";
import type { ReactNode } from "react";

import {
  OAuthClientList,
  OAuthRefreshTokenList,
  Permission,
} from "@/api/openapi-schema";
import { PermissionSummary } from "@/components/role/PermissionList";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { AddIcon } from "@/components/ui/icons/Add";
import { MetaGrid, MetaItem } from "@/components/ui/MetaGrid";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox, lstack } from "@/styled-system/patterns";
import { useDisclosure } from "@/utils/useDisclosure";

import { CreateOAuthClientModal } from "./CreateOAuthClientModal";
import { useOAuthClientSettings } from "./useOAuthClientSettings";
import { useOAuthTokenSettings } from "./useOAuthTokenSettings";

type Props = {
  tokens: OAuthRefreshTokenList;
  clients: OAuthClientList;
};

export function OAuthTokenSettings({ tokens, clients }: Props) {
  const createModal = useDisclosure();

  return (
    <>
      <LStack gap="8">
        <CardBox className={lstack()} gap="6">
          <LStack w="full">
            <Heading size="md">OAuth clients</Heading>
            <p>
              Create OAuth clients for integrations you own. Each client can
              only request the permission scopes selected here.
            </p>
          </LStack>

          <LStack>
            <WStack alignItems="center" color="fg.muted">
              <styled.p>{clients.length} clients.</styled.p>
              <Button size="xs" variant="subtle" onClick={createModal.onOpen}>
                <AddIcon />
                New
              </Button>
            </WStack>

            <OAuthClientItemList clients={clients} />
          </LStack>
        </CardBox>

        <CardBox className={lstack()} gap="6">
          <LStack>
            <Heading size="md">Authorised applications</Heading>
            <p>Applications you have authorised to access this site.</p>
          </LStack>

          <OAuthTokenItemList tokens={tokens} />
        </CardBox>
      </LStack>

      <CreateOAuthClientModal
        isOpen={createModal.isOpen}
        onClose={createModal.onClose}
      />
    </>
  );
}

function OAuthClientItemList({ clients }: { clients: OAuthClientList }) {
  const { deleteClient } = useOAuthClientSettings();

  if (clients.length === 0) {
    return (
      <styled.p color="fg.muted" fontStyle="italic">
        No OAuth clients created yet.
      </styled.p>
    );
  }

  return (
    <styled.ul className={lstack({ gap: "3" })} w="full">
      {clients.map((client) => (
        <OAuthClientItem
          key={client.id}
          client={client}
          onDelete={() => deleteClient(client.id)}
        />
      ))}
    </styled.ul>
  );
}

type OAuthClientItemProps = {
  client: OAuthClientList[number];
  onDelete: () => Promise<void>;
};

function OAuthClientItem({ client, onDelete }: OAuthClientItemProps) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(onDelete);

  const permissionScopes = client.allowed_scopes.filter(
    (scope): scope is Permission =>
      Object.values(Permission).includes(scope as Permission),
  );

  return (
    <OAuthRow>
      <LStack gap="2">
        <WStack gap="2" alignItems="start">
          <LStack gap="1" minW="0">
            <Heading size="sm">{client.name}</Heading>
            <styled.p color="fg.muted" fontSize="xs" wordBreak="break-word">
              {client.client_id}
            </styled.p>
          </LStack>

          <ConfirmActions
            confirming={isConfirming}
            confirmLabel="Confirm delete"
            idleLabel="Delete"
            onConfirm={handleConfirmAction}
            onCancel={handleCancelAction}
          />
        </WStack>

        <MetaGrid
          columns={{
            base: "1fr",
            md: "minmax(12rem, 18rem) minmax(7rem, 10rem) minmax(12rem, 1fr)",
          }}
        >
          <MetaItem label="Created">
            <time>{formatDate(client.createdAt, "PPpp")}</time>
          </MetaItem>
          <MetaItem label="Type">{client.type}</MetaItem>
          <MetaItem label="Grants">
            {client.allowed_grants.map(formatGrant).join(", ")}
          </MetaItem>
        </MetaGrid>

        <PermissionSummary permissions={permissionScopes} />
      </LStack>
    </OAuthRow>
  );
}

function OAuthTokenItemList({ tokens }: Pick<Props, "tokens">) {
  const { revokeToken } = useOAuthTokenSettings();

  if (tokens.length === 0) {
    return (
      <styled.p color="fg.muted" fontStyle="italic">
        No OAuth applications authorised yet.
      </styled.p>
    );
  }

  return (
    <styled.ul className={lstack({ gap: "3" })} w="full">
      {tokens.map((token) => (
        <OAuthTokenItem
          key={token.id}
          token={token}
          onRevoke={() => revokeToken(token.id)}
        />
      ))}
    </styled.ul>
  );
}

type OAuthTokenItemProps = {
  token: OAuthRefreshTokenList[number];
  onRevoke: () => Promise<void>;
};

function OAuthTokenItem({ token, onRevoke }: OAuthTokenItemProps) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(onRevoke);

  const inactiveStatus = token.revoked_at
    ? "Revoked"
    : new Date(token.expires_at) <= new Date()
      ? "Expired"
      : undefined;

  return (
    <OAuthRow>
      <LStack gap="2">
        <WStack gap="2" alignItems="start">
          <Heading size="sm">{token.client_name}</Heading>
          {inactiveStatus ? (
            <Badge>{inactiveStatus}</Badge>
          ) : (
            <ConfirmActions
              confirming={isConfirming}
              confirmLabel="Confirm revoke"
              idleLabel="Revoke"
              onConfirm={handleConfirmAction}
              onCancel={handleCancelAction}
            />
          )}
        </WStack>

        <MetaGrid
          columns={{
            base: "1fr",
            md: "minmax(12rem, 18rem) minmax(7rem, 10rem) minmax(12rem, 1fr)",
          }}
        >
          <MetaItem label="Created">
            <time>{formatDate(token.createdAt, "PPpp")}</time>
          </MetaItem>
          <MetaItem label="Expires">
            <time>{formatDate(token.expires_at, "PPpp")}</time>
          </MetaItem>
          <MetaItem label="Client">{token.client_id}</MetaItem>
        </MetaGrid>
      </LStack>
    </OAuthRow>
  );
}

function OAuthRow({ children }: { children: ReactNode }) {
  return (
    <styled.li className={cardBox()} w="full">
      {children}
    </styled.li>
  );
}

function ConfirmActions({
  confirming,
  confirmLabel,
  idleLabel,
  onConfirm,
  onCancel,
}: {
  confirming: boolean;
  confirmLabel: string;
  idleLabel: string;
  onConfirm: () => void;
  onCancel: () => void;
}) {
  if (confirming) {
    return (
      <HStack gap="2">
        <Button
          size="xs"
          variant="subtle"
          bgColor="bg.destructive"
          onClick={onConfirm}
        >
          {confirmLabel}
        </Button>
        <Button size="xs" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
      </HStack>
    );
  }

  return (
    <Button
      size="xs"
      variant="outline"
      bgColor="bg.destructive"
      onClick={onConfirm}
    >
      {idleLabel}
    </Button>
  );
}

function formatGrant(grant: string) {
  if (grant === "client_credentials") {
    return "client credentials";
  }

  if (grant === "authorization_code") {
    return "authorization code";
  }

  if (grant === "refresh_token") {
    return "refresh token";
  }

  if (grant.includes("device_code")) {
    return "device code";
  }

  return grant;
}
