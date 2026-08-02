"use client";

import { formatDate } from "date-fns";
import { useState } from "react";

import { EmailQueueItem, EmailQueueListResult } from "@/api/openapi-schema";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DateRangePicker } from "@/components/ui/date-picker";
import { Input } from "@/components/ui/input";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import { Box, Flex, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { cardBox } from "@/styled-system/recipes";

import {
  ALL_EMAIL_STATUSES,
  useEmailLogSettingsScreen,
} from "./useEmailLogSettingsScreen";

export function EmailLogSettingsScreen() {
  const {
    ready,
    data,
    error,
    filters,
    selectedStatuses,
    handleStatusFilterChange,
    handleDateRangeChange,
    handleResetDateRange,
    handleSearchChange,
    parseInitialDateValue,
    currentPage,
    refreshEmailLog,
    refreshing,
    retryEmail,
    retryingEmailId,
  } = useEmailLogSettingsScreen();
  const [statusQueryInput, setStatusQueryInput] = useState("");

  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  const filteredStatuses = ALL_EMAIL_STATUSES.filter((status) =>
    status.label.toLowerCase().includes(statusQueryInput.toLowerCase()),
  );

  return (
    <LStack gap="4">
      <WStack>
        <PageHeading>Email Log</PageHeading>

        <Button
          type="button"
          variant="subtle"
          loading={refreshing}
          onClick={refreshEmailLog}
        >
          Refresh
        </Button>
      </WStack>
      <Text variant="supporting">
        View queued emails and delivery attempts, including failure reasons.
      </Text>

      <Flex
        w="full"
        gap="2"
        flexDirection={{
          base: "column",
          md: "row",
        }}
      >
        <Input
          type="search"
          size="sm"
          value={filters.q ?? ""}
          onChange={(event) => handleSearchChange(event.target.value)}
          placeholder="Search recipient email..."
        />

        <MultiSelectPicker
          value={selectedStatuses}
          onChange={handleStatusFilterChange}
          onQuery={setStatusQueryInput}
          queryResults={filteredStatuses}
          inputPlaceholder="Filter by status..."
          size="sm"
        />

        <HStack gap="0">
          <DateRangePicker
            defaultValue={parseInitialDateValue()}
            onValueChange={handleDateRangeChange}
            active={!!parseInitialDateValue()}
            hideInputs={true}
            triggerClassName={css({
              borderRightRadius: "none",
            })}
          />
          <CancelAction
            variant="subtle"
            size="sm"
            borderLeftRadius="none"
            onClick={handleResetDateRange}
          />
        </HStack>
      </Flex>

      <EmailLogList
        data={data}
        currentPage={currentPage}
        filters={filters}
        selectedStatuses={selectedStatuses}
        onRetry={retryEmail}
        retryingEmailId={retryingEmailId}
      />
    </LStack>
  );
}

function EmailLogList({
  data,
  currentPage,
  filters,
  selectedStatuses,
  onRetry,
  retryingEmailId,
}: {
  data: EmailQueueListResult;
  currentPage: number;
  filters: {
    q?: string | null;
    range?: string | null;
  };
  selectedStatuses: MultiSelectPickerItem[];
  onRetry: (email: EmailQueueItem) => Promise<void>;
  retryingEmailId: string | null;
}) {
  const hasActiveFilters =
    !!filters.q || !!filters.range || selectedStatuses.length > 0;

  if (data.emails.length === 0) {
    return (
      <EmptyState w="full" hideContributionLabel>
        {hasActiveFilters
          ? "No emails match the current filters."
          : "No emails found."}
      </EmptyState>
    );
  }

  return (
    <>
      <ul className={lstack({ gap: "3" })}>
        {data.emails.map((email) => (
          <EmailItem
            key={email.id}
            email={email}
            onRetry={onRetry}
            retrying={retryingEmailId === email.id}
          />
        ))}
      </ul>

      <PaginationControls
        path="/admin/email"
        currentPage={currentPage}
        totalPages={data.total_pages}
        pageSize={data.page_size}
        params={buildEmailLogPaginationParams(filters, selectedStatuses)}
      />
    </>
  );
}

function EmailItem({
  email,
  onRetry,
  retrying,
}: {
  email: EmailQueueItem;
  onRetry: (email: EmailQueueItem) => Promise<void>;
  retrying: boolean;
}) {
  const lastAttempt = email.attempts[email.attempts.length - 1] ?? null;
  const error = lastAttempt?.error ?? null;
  const recipient = email.recipient_name
    ? `${email.recipient_name} <${email.recipient_address}>`
    : email.recipient_address;

  const attemptCountLabel =
    email.attempts.length === 1
      ? "1 delivery attempt"
      : `${email.attempts.length} delivery attempts`;
  const attempts = [...email.attempts].reverse();

  return (
    <li className={cardBox()}>
      <LStack gap="2">
        <WStack gap="2" alignItems="center">
          <HStack>
            <Badge
              variant="solid"
              colorPalette={email.status === "failed" ? "red" : "gray"}
            >
              {email.status}
            </Badge>
          </HStack>

          <styled.time fontSize="xs" color="text.subtle">
            {formatDate(email.queued_at, "PPpp")}
          </styled.time>
        </WStack>

        <WStack>
          <Text variant="supporting">“{email.subject}”</Text>

          <Text variant="metadata">{recipient}</Text>
        </WStack>

        <styled.details
          borderWidth="thin"
          borderColor="border.muted"
          borderRadius="md"
          bg="background.inset"
          overflow="hidden"
        >
          <styled.summary
            listStyle="none"
            cursor="pointer"
            px="3"
            py="2.5"
            fontSize="sm"
            fontWeight="medium"
            color="text.default"
            _marker={{ display: "none" }}
            css={{
              "&::-webkit-details-marker": {
                display: "none",
              },
            }}
          >
            <LStack gap="1">
              <span>{attemptCountLabel}</span>
              <Text as="span" variant="metadata" fontWeight="normal">
                {error ? error : "Open to review the full delivery history."}
              </Text>
            </LStack>
          </styled.summary>

          <LStack gap="2" px="3" pb="3">
            {attempts.map((attempt, index) => (
              <Box
                key={`${attempt.timestamp}-${index}`}
                borderTopWidth={index === 0 ? "none" : "thin"}
                borderColor="border.muted"
                pt={index === 0 ? "0" : "2"}
              >
                <Box
                  display="grid"
                  gap="2"
                  alignItems="start"
                  gridTemplateColumns={{
                    base: "1fr",
                    md: "auto minmax(0, 1fr)",
                  }}
                >
                  <WStack gap="0.5">
                    <Badge
                      size="sm"
                      variant="solid"
                      colorPalette={
                        attempt.status === "failed" ? "red" : "gray"
                      }
                    >
                      {attempt.status}
                    </Badge>

                    <styled.time fontSize="xs" color="text.subtle">
                      {formatDate(attempt.timestamp, "PPpp")}
                    </styled.time>
                  </WStack>

                  <styled.pre
                    m="0"
                    minW="0"
                    fontSize="xs"
                    lineHeight="normal"
                    color={attempt.error ? "text.default" : "text.subtle"}
                    whiteSpace="pre-wrap"
                    overflowWrap="anywhere"
                    wordBreak="break-word"
                  >
                    {attempt.error ?? "No error recorded."}
                  </styled.pre>
                </Box>
              </Box>
            ))}
          </LStack>
        </styled.details>

        <WStack alignItems="end">
          <Badge size="sm">{email.attempts.length} attempts</Badge>

          {email.status === "failed" && (
            <Button
              type="button"
              variant="subtle"
              loading={retrying}
              loadingText="Retrying..."
              onClick={() => onRetry(email)}
            >
              Retry Now
            </Button>
          )}
        </WStack>
      </LStack>
    </li>
  );
}

function buildEmailLogPaginationParams(
  filters: {
    q?: string | null;
    range?: string | null;
  },
  selectedStatuses: MultiSelectPickerItem[],
) {
  return {
    ...(filters.q ? { q: filters.q } : {}),
    ...(filters.range ? { range: filters.range } : {}),
    ...(selectedStatuses.length > 0
      ? { statuses: selectedStatuses.map((status) => status.value).join(",") }
      : {}),
  };
}
