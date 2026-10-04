"use client";

import { formatDate } from "date-fns";
import { parseAsInteger, useQueryState } from "nuqs";
import { useState } from "react";

import { handle } from "@/api/client";
import {
  adminOAuthDynamicRegistrationAccessTokenRevoke,
  useAdminOAuthDynamicRegistrationAccessTokenList,
} from "@/api/openapi-client/admin";
import { OAuthDynamicRegistrationAccessToken } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { MetaGrid, MetaItem } from "@/components/ui/meta-grid";
import { StatusBadge } from "@/components/ui/status-badge";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack } from "@/styled-system/jsx";
import { useDisclosure } from "@/utils/useDisclosure";

import { CreateRegistrationAccessToken } from "./CreateRegistrationAccessToken";

const STATUS_TONES = {
  Active: "success",
  Expired: "warning",
  Revoked: "danger",
  Exhausted: "info",
} as const;

export function RegistrationAccessTokens({ enabled }: { enabled: boolean }) {
  const [page, setPage] = useQueryState(
    "tokens_page",
    parseAsInteger.withDefault(1),
  );
  const tokens = useAdminOAuthDynamicRegistrationAccessTokenList(
    { page: page.toString() },
    { swr: { refreshInterval: 30_000 } },
  );
  const create = useDisclosure();
  const pagination = tokens.data && (
    <PaginationControls
      path="/admin/authentication/agents"
      currentPage={page}
      totalPages={tokens.data.total_pages}
      pageSize={tokens.data.page_size}
      onClick={(page) => {
        void setPage(page);
      }}
    />
  );

  async function revoke(id: string) {
    await adminOAuthDynamicRegistrationAccessTokenRevoke(id);
    await tokens.mutate();
  }

  return (
    <LStack gap="3">
      <Text variant="supporting">
        Give trusted actors a token to allow their agents to register without
        individual approval. Revoking a token stops future registrations.
      </Text>
      <WStack>
        {pagination}
        <Button ml="auto" onClick={create.onOpen} disabled={!enabled}>
          Create token
        </Button>
      </WStack>
      {tokens.error ? (
        <Unready error={tokens.error} />
      ) : !tokens.data ? (
        <Unready />
      ) : tokens.data.iats.length === 0 ? (
        <EmptyState hideContributionLabel>
          No registration tokens have been created.
        </EmptyState>
      ) : (
        <LStack as="ul" gap="3" w="full">
          {tokens.data.iats.map((token) => (
            <RegistrationAccessTokenItem
              key={token.id}
              token={token}
              onRevoke={() => revoke(token.id)}
            />
          ))}
        </LStack>
      )}
      {pagination}
      {create.isOpen && (
        <CreateRegistrationAccessToken onClose={create.onClose} />
      )}
    </LStack>
  );
}

function RegistrationAccessTokenItem({
  token,
  onRevoke,
}: {
  token: OAuthDynamicRegistrationAccessToken;
  onRevoke: () => Promise<void>;
}) {
  const [revoking, setRevoking] = useState(false);
  const confirmation = useConfirmation(async () => {
    setRevoking(true);
    await handle(onRevoke, {
      cleanup: async () => {
        setRevoking(false);
      },
    });
  });
  const status = token.revoked_at
    ? "Revoked"
    : new Date(token.expires_at).getTime() <= Date.now()
      ? "Expired"
      : token.registration_count >= token.max_registrations
        ? "Exhausted"
        : "Active";

  return (
    <CardBox as="li" w="full">
      <LStack gap="3">
        <WStack>
          <Text variant="supporting" color="text.default" fontWeight="semibold">
            {token.label}
          </Text>
          <HStack>
            {status === "Active" && (
              <HStack>
                <Button
                  variant="outline"
                  intent="destructive"
                  onClick={confirmation.handleConfirmAction}
                  loading={revoking}
                >
                  {confirmation.isConfirming ? "Confirm revoke" : "Revoke"}
                </Button>
                {confirmation.isConfirming && (
                  <Button
                    variant="outline"
                    onClick={confirmation.handleCancelAction}
                    disabled={revoking}
                  >
                    Cancel
                  </Button>
                )}
              </HStack>
            )}
            <StatusBadge tone={STATUS_TONES[status]}>{status}</StatusBadge>
          </HStack>
        </WStack>
        <MetaGrid>
          <MetaItem label="Registrations used">
            {token.registration_count} / {token.max_registrations}
          </MetaItem>
          <MetaItem label="Created">
            <time dateTime={token.createdAt}>
              {formatDate(token.createdAt, "PPpp")}
            </time>
          </MetaItem>
          <MetaItem label={token.revoked_at ? "Revoked" : "Expires"}>
            <time dateTime={token.revoked_at ?? token.expires_at}>
              {formatDate(token.revoked_at ?? token.expires_at, "PPpp")}
            </time>
          </MetaItem>
        </MetaGrid>
      </LStack>
    </CardBox>
  );
}
