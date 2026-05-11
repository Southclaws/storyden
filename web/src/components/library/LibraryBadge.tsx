import Link from "next/link";

import { Node } from "@/api/openapi-schema";
import { useI18n } from "@/i18n/provider";
import { styled } from "@/styled-system/jsx";

export function LibraryBadge() {
  const { t } = useI18n();

  return (
    <styled.span
      position="relative"
      backgroundColor="bg.accent"
      color="fg.accent"
      px="1"
      borderRadius="md"
    >
      <Link href="/l">{t("library")}</Link>
    </styled.span>
  );
}

export function LibraryPageBadge(props: Node) {
  const url = `/l/${props.slug}`;
  return (
    <styled.span
      position="relative"
      backgroundColor="bg.accent"
      color="fg.accent"
      px="1"
      borderRadius="md"
      lineClamp="1"
    >
      <Link href={url}>{props.name}</Link>
    </styled.span>
  );
}

// TODO: Make this a recipe component.
export function NewBadge() {
  const { t } = useI18n();

  return (
    <styled.span
      fontSize="xs"
      fontWeight="bold"
      backgroundColor="bg.accent"
      color="fg.accent"
      px="1"
      py="0.5"
      borderRadius="sm"
    >
      {t("New")}
    </styled.span>
  );
}
