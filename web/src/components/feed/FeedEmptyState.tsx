"use client";

import { EmptyState } from "../site/EmptyState";
import { EmptyThreadsIcon } from "../ui/icons/Empty";
import { useI18n } from "@/i18n/provider";

export function FeedEmptyState() {
  const { t } = useI18n();

  return (
    <EmptyState w="full" icon={<EmptyThreadsIcon />}>
      <p>{t("*tumbleweed* there are no posts...")}</p>
    </EmptyState>
  );
}
