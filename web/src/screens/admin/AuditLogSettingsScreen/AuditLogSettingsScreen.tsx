"use client";

import { formatDate } from "date-fns";
import Link from "next/link";
import { useState } from "react";

import {
  AuditEvent,
  AuditEventListResult,
  AuditEventType,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Badge } from "@/components/ui/badge";
import { DateRangePicker } from "@/components/ui/date-picker";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import { Flex, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { cardBox } from "@/styled-system/recipes";

import {
  ALL_EVENT_TYPES,
  EVENT_TYPE_LABELS,
  useAuditLogSettingsScreen,
} from "./useAuditLogSettingsScreen";

export function AuditLogSettingsScreen() {
  const {
    ready,
    data,
    error,
    selectedTypes,
    handleTypeFilterChange,
    handleDateRangeChange,
    handleResetDateRange,
    parseInitialDateValue,
    currentPage,
  } = useAuditLogSettingsScreen();

  const [queryInput, setQueryInput] = useState("");

  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  const filteredEventTypes = ALL_EVENT_TYPES.filter((type) =>
    type.label.toLowerCase().includes(queryInput.toLowerCase()),
  );

  return (
    <LStack gap="4">
      <PageHeading>Audit Log</PageHeading>
      <Text variant="supporting">
        View all moderation actions and administrative events on this site.
      </Text>

      <Flex
        w="full"
        gap="2"
        flexDirection={{
          base: "column",
          md: "row",
        }}
      >
        <MultiSelectPicker
          value={selectedTypes}
          onChange={handleTypeFilterChange}
          onQuery={setQueryInput}
          queryResults={filteredEventTypes}
          inputPlaceholder="Filter by event type..."
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

      <AuditLogEventList
        data={data}
        currentPage={currentPage}
        selectedTypes={selectedTypes}
      />
    </LStack>
  );
}

function AuditLogEventList({
  data,
  currentPage,
  selectedTypes,
}: {
  data: AuditEventListResult;
  currentPage: number;
  selectedTypes: MultiSelectPickerItem[];
}) {
  if (data.events.length === 0) {
    return (
      <EmptyState hideContributionLabel>No audit events found.</EmptyState>
    );
  }

  return (
    <>
      <ul className={lstack({ gap: "3" })}>
        {data.events.map((event) => (
          <AuditEventItem key={event.id} event={event} />
        ))}
      </ul>

      <PaginationControls
        path="/admin/audit-log"
        currentPage={currentPage}
        totalPages={data.total_pages}
        pageSize={data.page_size}
        params={{
          types: selectedTypes.map((v) => v.value).join(","),
        }}
      />
    </>
  );
}

type AuditEventItemProps = {
  event: AuditEvent;
};

function AuditEventItem({ event }: AuditEventItemProps) {
  const eventLabel = EVENT_TYPE_LABELS[event.type] ?? event.type;

  return (
    <li className={cardBox()}>
      <LStack gap="2">
        <WStack gap="2" alignItems="center">
          <Badge>{eventLabel}</Badge>
          <styled.time fontSize="xs" color="text.subtle">
            {formatDate(event.timestamp, "PPpp")}
          </styled.time>
        </WStack>

        <HStack>
          <Text variant="supporting">Enacted by:</Text>

          {event.enacted_by && (
            <MemberBadge
              size="sm"
              name="full-horizontal"
              profile={event.enacted_by}
            />
          )}
        </HStack>

        <EventDetails event={event} />
      </LStack>
    </li>
  );
}

function EventDetails({ event }: { event: AuditEvent }) {
  switch (event.type) {
    case AuditEventType.thread_deleted:
      return (
        <Text variant="supporting">
          Thread:{" "}
          <Link
            href={`/t/locate?id=${event.thread_id}`}
            className={css({ textDecoration: "underline" })}
          >
            <code>{event.thread_id}</code>
          </Link>
        </Text>
      );

    case AuditEventType.thread_reply_deleted:
      return (
        <Text variant="supporting">
          Reply:{" "}
          <Link
            href={`/t/locate?id=${event.reply_id}`}
            className={css({ textDecoration: "underline" })}
          >
            <code>{event.reply_id}</code>
          </Link>
        </Text>
      );

    case AuditEventType.account_suspended:
      return (
        <Text variant="supporting">
          Account:{" "}
          <Link
            href={`/m/${event.account_id}`}
            className={css({ textDecoration: "underline" })}
          >
            <code>{event.account_id}</code>
          </Link>
        </Text>
      );

    case AuditEventType.account_unsuspended:
      return (
        <Text variant="supporting">
          Account:{" "}
          <Link
            href={`/m/${event.account_id}`}
            className={css({ textDecoration: "underline" })}
          >
            <code>{event.account_id}</code>
          </Link>
        </Text>
      );

    case AuditEventType.moderation_note_created:
    case AuditEventType.moderation_note_deleted:
      return (
        <LStack gap="1">
          <Text variant="supporting">
            Account:{" "}
            <Link
              href={`/m/${event.account_id}`}
              className={css({ textDecoration: "underline" })}
            >
              <code>{event.account_id}</code>
            </Link>
          </Text>
          <Text variant="supporting">
            Note: <code>{event.note_id}</code>
          </Text>
        </LStack>
      );

    case AuditEventType.account_content_purged:
      return (
        <LStack gap="1">
          <Text variant="supporting">
            Account:{" "}
            <Link
              href={`/m/${event.account_id}`}
              className={css({ textDecoration: "underline" })}
            >
              <code>{event.account_id}</code>
            </Link>
          </Text>
          {event.included && event.included.length > 0 && (
            <HStack gap="1" flexWrap="wrap">
              <Text variant="supporting">Purged Content:</Text>
              {event.included.map((type) => (
                <Badge key={type} size="sm" variant="subtle">
                  {type}
                </Badge>
              ))}
            </HStack>
          )}
        </LStack>
      );

    default:
      return null;
  }
}
