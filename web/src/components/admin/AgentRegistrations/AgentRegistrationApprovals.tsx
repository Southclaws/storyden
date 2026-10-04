"use client";

import { formatDate } from "date-fns";
import { useState } from "react";

import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { MetaGrid, MetaItem } from "@/components/ui/meta-grid";
import { StatusBadge } from "@/components/ui/status-badge";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { useRegistrationApprovals } from "./useRegistrationApprovals";

type Props = { enabled: boolean; defaultRoleName?: string };

export function AgentRegistrationApprovals({
  enabled,
  defaultRoleName,
}: Props) {
  const approvals = useRegistrationApprovals(enabled);
  const [bulkSnapshot, setBulkSnapshot] = useState("");
  const approveAll = useConfirmation(() =>
    approvals.decideAll(true, bulkSnapshot),
  );
  const denyAll = useConfirmation(() =>
    approvals.decideAll(false, bulkSnapshot),
  );
  const busy = !!approvals.submitting || approvals.bulkSubmitting;
  const confirming = approveAll.isConfirming || denyAll.isConfirming;
  const pagination = approvals.queue.data && (
    <PaginationControls
      path="/admin/authentication/agents"
      currentPage={approvals.filters.approvals_page}
      totalPages={approvals.queue.data.total_pages}
      pageSize={approvals.queue.data.page_size}
      onClick={(page) => {
        void approvals.setFilters({ approvals_page: page });
      }}
    />
  );

  return (
    <LStack gap="3">
      {!enabled ? (
        <Text variant="supporting">
          Enable admin approval in authentication settings to review new bot
          accounts.
        </Text>
      ) : (
        <>
          <Text variant="supporting">
            Only approve requests you recognise. Confirm the code with the
            person running the agent. New accounts inherit their own roles, not
            your permissions.
          </Text>
          <Admonition
            value={!!approvals.outcome}
            onChange={approvals.dismissOutcome}
            kind="success"
            title={
              approvals.outcome === "approved"
                ? "Registration approved"
                : "Registration denied"
            }
          >
            {approvals.bulkUpdated !== undefined
              ? `${approvals.bulkUpdated} pending registrations ${approvals.outcome}.`
              : approvals.outcome === "approved"
                ? "The agent can finish registering when it next checks for approval."
                : "This request can no longer create an account."}
          </Admonition>
          <Admonition
            value={!!approvals.submitError}
            onChange={approvals.dismissSubmitError}
            kind="failure"
            title="Decision could not be saved"
          >
            {approvals.submitError}
          </Admonition>

          {approvals.filters.verification_code &&
            !approvals.review &&
            (approvals.lookup.error ? (
              <Admonition
                value
                kind="failure"
                title="Registration unavailable"
                onChange={() =>
                  approvals.setFilters({
                    verification_code: null,
                    tab: "reviews",
                  })
                }
              >
                {approvals.lookupErrorMessage}
              </Admonition>
            ) : (
              <Unready />
            ))}

          {confirming && (
            <Text variant="supporting">
              {approveAll.isConfirming ? "Approve" : "Deny"} all pending
              requests across all pages from this queue snapshot? Requests
              received since it was loaded are excluded.
            </Text>
          )}
          <WStack flexWrap="wrap">
            {pagination}
            <HStack flex="1" flexWrap="wrap" justifyContent="end">
              <Button
                variant="subtle"
                disabled={
                  busy ||
                  denyAll.isConfirming ||
                  !approvals.queue.data?.snapshot_at ||
                  approvals.registrations.length === 0
                }
                loading={approvals.bulkSubmitting && approveAll.isConfirming}
                onClick={() => {
                  if (!approveAll.isConfirming)
                    setBulkSnapshot(approvals.queue.data!.snapshot_at);
                  void approveAll.handleConfirmAction();
                }}
              >
                {approveAll.isConfirming
                  ? "Confirm approve all"
                  : "Approve all"}
              </Button>
              <Button
                variant="outline"
                intent="destructive"
                disabled={
                  busy ||
                  approveAll.isConfirming ||
                  !approvals.queue.data?.snapshot_at ||
                  approvals.registrations.length === 0
                }
                loading={approvals.bulkSubmitting && denyAll.isConfirming}
                onClick={() => {
                  if (!denyAll.isConfirming)
                    setBulkSnapshot(approvals.queue.data!.snapshot_at);
                  void denyAll.handleConfirmAction();
                }}
              >
                {denyAll.isConfirming ? "Confirm deny all" : "Deny all"}
              </Button>
              {confirming && (
                <Button
                  disabled={busy}
                  onClick={() => {
                    void approveAll.handleCancelAction();
                    void denyAll.handleCancelAction();
                  }}
                >
                  Cancel
                </Button>
              )}
              <Button
                onClick={() => approvals.queue.mutate()}
                loading={approvals.queue.isValidating}
                disabled={busy || confirming}
              >
                Refresh
              </Button>
            </HStack>
          </WStack>
          {approvals.queue.error ? (
            <Unready error={approvals.queue.error} />
          ) : !approvals.queue.data ? (
            <Unready />
          ) : (
            <>
              {approvals.registrations.length === 0 ? (
                <EmptyState hideContributionLabel>
                  No registrations are awaiting approval.
                </EmptyState>
              ) : (
                <styled.ul className={lstack({ gap: "3" })} w="full">
                  {approvals.registrations.map((registration) => (
                    <CardBox
                      as="li"
                      key={registration.verification_code}
                      w="full"
                    >
                      <WStack>
                        <Text>@{registration.metadata.client_name}</Text>
                        <StatusBadge tone="neutral">Pending</StatusBadge>
                      </WStack>
                      <MetaGrid>
                        <MetaItem label="Verification code">
                          {registration.verification_code}
                        </MetaItem>
                        <MetaItem label="Requested">
                          <time dateTime={registration.created_at}>
                            {formatDate(registration.created_at, "PPpp")}
                          </time>
                        </MetaItem>
                        <MetaItem label="Expires">
                          <time dateTime={registration.expires_at}>
                            {formatDate(registration.expires_at, "PPpp")}
                          </time>
                        </MetaItem>
                        <MetaItem label="Default role">
                          {defaultRoleName ?? "Member only"}
                        </MetaItem>
                      </MetaGrid>
                      <HStack w="full" justifyContent="end">
                        <Button
                          variant="subtle"
                          onClick={() =>
                            approvals.decide(
                              registration.verification_code,
                              true,
                            )
                          }
                          disabled={busy || confirming}
                          loading={
                            approvals.submitting?.code ===
                              registration.verification_code &&
                            approvals.submitting.approved
                          }
                        >
                          Approve
                        </Button>
                        <Button
                          variant="outline"
                          intent="destructive"
                          onClick={() =>
                            approvals.decide(
                              registration.verification_code,
                              false,
                            )
                          }
                          disabled={busy || confirming}
                          loading={
                            approvals.submitting?.code ===
                              registration.verification_code &&
                            !approvals.submitting.approved
                          }
                        >
                          Deny
                        </Button>
                      </HStack>
                    </CardBox>
                  ))}
                </styled.ul>
              )}
              {pagination}
            </>
          )}
        </>
      )}
    </LStack>
  );
}
