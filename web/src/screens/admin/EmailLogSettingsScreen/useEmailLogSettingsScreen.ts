"use client";

import { parseAsInteger, useQueryStates } from "nuqs";
import { useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import { emailQueueRetry, useEmailQueueList } from "@/api/openapi-client/admin";
import { EmailQueueItem, EmailQueueListResult } from "@/api/openapi-schema";

export function useEmailLogSettingsScreen() {
  const { mutate } = useSWRConfig();
  const [filters] = useQueryStates({
    page: parseAsInteger.withDefault(1),
  });
  const [retryingEmailId, setRetryingEmailId] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const emailQueueList = useEmailQueueList({
    page: filters.page.toString(),
  });
  const { data, error, swrKey, mutate: mutateEmailQueueList } = emailQueueList;

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
