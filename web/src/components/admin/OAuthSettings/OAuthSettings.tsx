"use client";

import { formatDate } from "date-fns";
import type { ReactNode } from "react";

import {
  OAuthClientList,
  OAuthDeviceAuthorisationList,
  OAuthRefreshTokenList,
  Permission,
} from "@/api/openapi-schema";
import { PermissionSummary } from "@/components/role/PermissionList";
import { useConfirmation } from "@/components/site/useConfirmation";
import { MetaGrid, MetaItem } from "@/components/ui/MetaGrid";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox, lstack } from "@/styled-system/patterns";

import { useAdminOAuthSettings } from "./useAdminOAuthSettings";

type Props = {
  clients: OAuthClientList;
  deviceAuthorisations: OAuthDeviceAuthorisationList;
  tokens: OAuthRefreshTokenList;
};

export function OAuthSettings({
  clients,
  deviceAuthorisations,
  tokens,
}: Props) {
  const activeTokens = tokens.filter((token) => !token.revoked_at).length;

  return (
    <CardBox className={lstack()} gap="4">
      <LStack gap="2">
        <Heading size="md">OAuth</Heading>
        <p>OAuth clients, device authorisations, and refresh tokens.</p>
      </LStack>

      <OAuthClientListView clients={clients} />
      <OAuthRefreshTokenListView tokens={tokens} activeTokens={activeTokens} />
      <OAuthDeviceAuthorisationListView
        deviceAuthorisations={deviceAuthorisations}
      />
    </CardBox>
  );
}

function OAuthClientListView({ clients }: { clients: OAuthClientList }) {
  if (clients.length === 0) {
    return <Empty title="Clients" body="No OAuth clients registered yet." />;
  }

  return (
    <LStack gap="3">
      <Heading size="sm">Clients</Heading>
      <styled.ul className={lstack({ gap: "3" })} w="full">
        {clients.map((client) => (
          <OAuthClientItem key={client.id} client={client} />
        ))}
      </styled.ul>
    </LStack>
  );
}

function OAuthClientItem({ client }: { client: OAuthClientList[number] }) {
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
          <HStack gap="2">
            <Badge>{client.type}</Badge>
            <Badge>{client.scope_policy}</Badge>
          </HStack>
        </WStack>

        <MetaGrid>
          <MetaItem label="Created">
            <time>{formatDate(client.createdAt, "PPpp")}</time>
          </MetaItem>
          <MetaItem label="Grants">
            {client.allowed_grants.map(formatGrant).join(", ")}
          </MetaItem>
          <MetaItem label="Redirects">
            {client.redirect_uris?.length || "None"}
          </MetaItem>
        </MetaGrid>

        <PermissionSummary permissions={permissionScopes} />
      </LStack>
    </OAuthRow>
  );
}

function OAuthRefreshTokenListView({
  tokens,
  activeTokens,
}: {
  tokens: OAuthRefreshTokenList;
  activeTokens: number;
}) {
  const { revokeToken } = useAdminOAuthSettings();

  if (tokens.length === 0) {
    return <Empty title="Refresh tokens" body="No OAuth tokens issued yet." />;
  }

  return (
    <LStack gap="3">
      <styled.div display="flex" gap="3" alignItems="baseline" flexWrap="wrap">
        <Heading size="sm">Refresh tokens</Heading>
        <styled.p color="fg.muted" fontSize="sm">
          {tokens.length} tokens, {activeTokens} active.
        </styled.p>
      </styled.div>
      <styled.ul className={lstack({ gap: "3" })} w="full">
        {tokens.map((token) => (
          <OAuthRefreshTokenItem
            key={token.id}
            token={token}
            onRevoke={() => revokeToken(token.id)}
          />
        ))}
      </styled.ul>
    </LStack>
  );
}

function OAuthRefreshTokenItem({
  token,
  onRevoke,
}: {
  token: OAuthRefreshTokenList[number];
  onRevoke: () => Promise<void>;
}) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(onRevoke);

  return (
    <OAuthRow>
      <LStack gap="2">
        <WStack gap="2" alignItems="start">
          <Heading size="sm">{token.client_name}</Heading>
          {token.revoked_at ? (
            <Badge>Revoked</Badge>
          ) : isConfirming ? (
            <HStack flexWrap="wrap" gap="2">
              <Button
                size="xs"
                variant="subtle"
                bgColor="bg.destructive"
                onClick={handleConfirmAction}
              >
                Confirm revoke
              </Button>
              <Button size="xs" variant="outline" onClick={handleCancelAction}>
                Cancel
              </Button>
            </HStack>
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
        </WStack>

        <MetaGrid>
          <MetaItem label="Expires">
            <time>{formatDate(token.expires_at, "PPpp")}</time>
          </MetaItem>
          <MetaItem label="Client">{token.client_id}</MetaItem>
          <MetaItem label="Issued">
            <time>{formatDate(token.createdAt, "PPpp")}</time>
          </MetaItem>
        </MetaGrid>
      </LStack>
    </OAuthRow>
  );
}

function OAuthDeviceAuthorisationListView({
  deviceAuthorisations,
}: {
  deviceAuthorisations: OAuthDeviceAuthorisationList;
}) {
  if (deviceAuthorisations.length === 0) {
    return (
      <Empty
        title="Device authorisations"
        body="No device authorisation attempts yet."
      />
    );
  }

  return (
    <LStack gap="3">
      <Heading size="sm">Device authorisations</Heading>
      <styled.ul className={lstack({ gap: "3" })} w="full">
        {deviceAuthorisations.map((device) => (
          <OAuthRow key={device.id}>
            <WStack gap="2" alignItems="center">
              <Heading size="sm">{device.user_code}</Heading>
              <Badge>
                {device.approved_at
                  ? "Approved"
                  : device.denied_at
                    ? "Denied"
                    : "Pending"}
              </Badge>
            </WStack>
          </OAuthRow>
        ))}
      </styled.ul>
    </LStack>
  );
}

function OAuthRow({ children }: { children: ReactNode }) {
  return (
    <styled.li className={cardBox()} w="full">
      {children}
    </styled.li>
  );
}

function formatGrant(grant: string) {
  if (grant === "client_credentials") return "client credentials";
  if (grant === "authorization_code") return "authorization code";
  if (grant === "refresh_token") return "refresh token";
  if (grant.includes("device_code")) return "device code";

  return grant;
}

function Empty({ title, body }: { title: string; body: string }) {
  return (
    <LStack>
      <Heading size="sm">{title}</Heading>
      <styled.p color="fg.muted" fontStyle="italic">
        {body}
      </styled.p>
    </LStack>
  );
}
