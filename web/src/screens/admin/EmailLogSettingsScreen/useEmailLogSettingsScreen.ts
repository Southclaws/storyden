"use client";

import { type DateValue, parseDate } from "@internationalized/date";
import {
  parseAsArrayOf,
  parseAsInteger,
  parseAsString,
  useQueryStates,
} from "nuqs";
import { useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import { emailQueueRetry, useEmailQueueList } from "@/api/openapi-client/admin";
import {
  EmailQueueItem,
  EmailQueueListResult,
  EmailQueueStatus,
} from "@/api/openapi-schema";
import { MultiSelectPickerItem } from "@/components/ui/MultiSelectPicker";

export const EMAIL_STATUS_LABELS: Record<EmailQueueStatus, string> = {
  [EmailQueueStatus.pending]: "Pending",
  [EmailQueueStatus.processing]: "Processing",
  [EmailQueueStatus.sent]: "Sent",
  [EmailQueueStatus.failed]: "Failed",
};

export const ALL_EMAIL_STATUSES: MultiSelectPickerItem[] = Object.entries(
  EmailQueueStatus,
).map(([_, value]) => ({
  label: EMAIL_STATUS_LABELS[value],
  value,
}));

export function useEmailLogSettingsScreen() {
  const { mutate } = useSWRConfig();
  const [filters, setFilters] = useQueryStates({
    q: parseAsString,
    statuses: parseAsArrayOf(parseAsString),
    range: parseAsString,
    page: parseAsInteger.withDefault(1),
  });
  const [retryingEmailId, setRetryingEmailId] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const emailQueueList = useEmailQueueList({
    q: filters.q ?? undefined,
    statuses: (filters.statuses as EmailQueueStatus[]) ?? [],
    range: filters.range ?? undefined,
    page: filters.page.toString(),
  });
  const { data, error, swrKey, mutate: mutateEmailQueueList } = emailQueueList;

  const selectedStatuses: MultiSelectPickerItem[] =
    filters.statuses
      ?.map((status) => ({
        label: EMAIL_STATUS_LABELS[status as EmailQueueStatus],
        value: status,
      }))
      .filter((item) => item.label !== undefined) ?? [];

  const handleStatusFilterChange = async (items: MultiSelectPickerItem[]) => {
    await setFilters({
      statuses: items.length > 0 ? items.map((item) => item.value) : null,
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

    await setFilters({
      range: `${startISO}/${endISO}`,
      page: 1,
    });
  };

  const handleResetDateRange = async () => {
    await setFilters({ range: null, page: 1 });
  };

  const handleSearchChange = async (query: string) => {
    await setFilters({
      q: query.length > 0 ? query : null,
      page: 1,
    });
  };

  const parseInitialDateValue = () => {
    if (!filters.range) {
      return undefined;
    }

    const parts = filters.range.split("/").filter((part) => part);
    try {
      return parts.map((part) => parseDate(part));
    } catch {
      return undefined;
    }
  };

  const refreshEmailLog = async () => {
    if (!swrKey) {
      return;
    }

    setRefreshing(true);

    await handle(
      async () => {
        await new Promise((resolve) => setTimeout(resolve, 500));
        await mutateEmailQueueList();
      },
      {
        cleanup: async () => {
          setRefreshing(false);
        },
      },
    );
  };

  const retryEmail = async (email: EmailQueueItem) => {
    if (email.status !== "failed" || !swrKey) {
      return;
    }

    setRetryingEmailId(email.id);

    await handle(
      async () => {
        await mutateTransaction(
          mutate,
          [
            {
              key: swrKey,
              optimistic: (current: EmailQueueListResult | undefined) =>
                upsertEmailQueueItem(current, {
                  ...email,
                  updated_at: new Date().toISOString(),
                }),
              commit: (
                current: EmailQueueListResult | undefined,
                result: EmailQueueItem,
              ) => upsertEmailQueueItem(current, result),
            },
          ],
          () => emailQueueRetry(email.id),
          { revalidate: true },
        );
      },
      {
        cleanup: async () => {
          setRetryingEmailId(null);
        },
      },
    );
  };

  if (!data) {
    return {
      ready: false as const,
      error,
      filters,
      selectedStatuses,
      handleStatusFilterChange,
      handleDateRangeChange,
      handleResetDateRange,
      handleSearchChange,
      parseInitialDateValue,
      currentPage: filters.page,
      refreshEmailLog,
      refreshing,
      retryEmail,
      retryingEmailId,
    };
  }

  return {
    ready: true as const,
    data,
    filters,
    selectedStatuses,
    handleStatusFilterChange,
    handleDateRangeChange,
    handleResetDateRange,
    handleSearchChange,
    parseInitialDateValue,
    currentPage: filters.page,
    refreshEmailLog,
    refreshing,
    retryEmail,
    retryingEmailId,
  };
}

function upsertEmailQueueItem(
  current: EmailQueueListResult | undefined,
  next: EmailQueueItem,
): EmailQueueListResult | undefined {
  if (!current || !current.emails) {
    return current;
  }

  return {
    ...current,
    emails: current.emails.map((email) =>
      email.id === next.id ? next : email,
    ),
  };
}
