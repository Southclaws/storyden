import { LinkListResult } from "src/api/openapi-schema";

import { EmptyState } from "@/components/site/EmptyState";
import { useI18n } from "@/i18n/provider";

import { LinkCardRows } from "../LinkCardList";

type Props = {
  links: LinkListResult;
  show?: number;
};

export function LinkResultList({ links, show }: Props) {
  const { t } = useI18n();

  if (links.links.length === 0) {
    return <EmptyState hideContributionLabel>{t("No links were found.")}</EmptyState>;
  }

  const shown = show ? links.links.slice(0, show) : links.links;

  return <LinkCardRows links={shown} />;
}
