"use client";

import { JSX, PropsWithChildren } from "react";

import { useSession } from "@/auth";
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
  const session = useSession();

  const contributionLabel = session
    ? (authenticatedLabel ?? "Be the first to contribute!")
    : (unauthenticatedLabel ?? "Please log in to contribute.");

  return (
    <Center className={vstack(props)} p="8" gap="2" color="fg.subtle">
      {icon || <EmptyIcon />}

      <VStack gap="1" textAlign="center" fontStyle="italic">
        {children || <p>There&apos;s no content here.</p>}
        {!hideContributionLabel && <p>{contributionLabel}</p>}
      </VStack>
    </Center>
  );
}
