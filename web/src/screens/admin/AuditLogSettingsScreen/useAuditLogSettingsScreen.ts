"use client";

import { type DateValue, parseDate } from "@internationalized/date";
import {
  parseAsArrayOf,
  parseAsInteger,
  parseAsString,
  useQueryStates,
} from "nuqs";

import { useAuditEventList } from "@/api/openapi-client/admin";
import { AuditEventType } from "@/api/openapi-schema";
import { MultiSelectPickerItem } from "@/components/ui/MultiSelectPicker";

export const EVENT_TYPE_LABELS: Record<AuditEventType, string> = {
  [AuditEventType.thread_deleted]: "Thread Deleted",
  [AuditEventType.thread_reply_deleted]: "Reply Deleted",
  [AuditEventType.account_suspended]: "Account Suspended",
  [AuditEventType.account_unsuspended]: "Account Unsuspended",
  [AuditEventType.account_content_purged]: "Content Purged",
};

export const ALL_EVENT_TYPES: MultiSelectPickerItem[] = Object.entries(
  AuditEventType,
).map(([_, value]) => ({
  label: EVENT_TYPE_LABELS[value],
  value,
}));

export function useAuditLogSettingsScreen() {
  const [filters, setFilters] = useQueryStates({
    types: parseAsArrayOf(parseAsString),
    range: parseAsString,
    page: parseAsInteger.withDefault(1),
  });

  const { data, error } = useAuditEventList({
    types: (filters.types as AuditEventType[]) ?? [],
    range: filters.range ?? undefined,
    page: filters.page.toString(),
  });

  const selectedTypes: MultiSelectPickerItem[] =
    filters.types
      ?.map((type) => ({
        label: EVENT_TYPE_LABELS[type as AuditEventType],
        value: type,
      }))
      .filter((item) => item.label !== undefined) ?? [];

  const handleTypeFilterChange = async (items: MultiSelectPickerItem[]) => {
    await setFilters({
      types: items.length > 0 ? items.map((i) => i.value) : null,
      page: 1,
    });
  };

  const handleDateRangeChange = async (details: { value: DateValue[] }) => {
    const [start, end] = details.value;

    if (!start && !end) {
      await setFilters({ range: null, page: 1 });
      return;
    }

    const startISO = start ? start.toString() : "";
    const endISO = end ? end.toString() : "";
    const rangeString = `${startISO}/${endISO}`;

    await setFilters({ range: rangeString, page: 1 });
  };

  const handleResetDateRange = async () => {
    await setFilters({ range: null, page: 1 });
  };

  const parseInitialDateValue = () => {
    if (!filters.range) return undefined;

    const parts = filters.range.split("/").filter((p) => p);
    try {
      return parts.map((p) => parseDate(p));
    } catch {
      return undefined;
    }
  };

  if (!data) {
    return {
      ready: false as const,
      error,
      selectedTypes,
      handleTypeFilterChange,
      handleDateRangeChange,
      handleResetDateRange,
      parseInitialDateValue,
      currentPage: filters.page,
    };
  }

  return {
    ready: true as const,
    data,
    selectedTypes,
    handleTypeFilterChange,
    handleDateRangeChange,
    handleResetDateRange,
    parseInitialDateValue,
    currentPage: filters.page,
  };
}
