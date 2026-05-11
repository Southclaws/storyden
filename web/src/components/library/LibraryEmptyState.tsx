import { useI18n } from "@/i18n/provider";

import { EmptyState } from "../site/EmptyState";

export function LibraryEmptyState() {
  const { t } = useI18n();

  return (
    <EmptyState w="full">
      <p>{t("The community library is empty.")}</p>
    </EmptyState>
  );
}
