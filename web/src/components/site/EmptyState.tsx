"use client";

import { JSX, PropsWithChildren } from "react";

import { useSession } from "@/auth";
import { useI18n } from "@/i18n/provider";
import { Center, VStack, VstackProps } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { EmptyIcon } from "../ui/icons/Empty";

type Props = {
  icon?: JSX.Element;
  unauthenticatedLabel?: string;
  authenticatedLabel?: string;
  hideContributionLabel?: boolean;
};

export function EmptyState({
  icon,
  unauthenticatedLabel,
  authenticatedLabel,
  hideContributionLabel,
  children,
  ...props
}: PropsWithChildren<Props & VstackProps>) {
  const { t } = useI18n();
  const session = useSession();

  const contributionLabel = session
    ? (authenticatedLabel ?? t("Be the first to contribute!"))
    : (unauthenticatedLabel ?? t("Please log in to contribute."));

  return (
    <Center className={vstack(props)} p="8" gap="2" color="fg.subtle">
      {icon || <EmptyIcon />}

      <VStack gap="1" textAlign="center" fontStyle="italic">
        {children || <p>{t("There's no content here.")}</p>}
        {!hideContributionLabel && <p>{contributionLabel}</p>}
      </VStack>
    </Center>
  );
}
