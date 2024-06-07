"use client";

import ErrorBanner from "src/components/site/ErrorBanner";
import { Unready } from "src/components/site/Unready";

import { Heading1 } from "@/components/ui/typography-heading";
import { VStack } from "@/styled-system/jsx";

import { BrandSettings } from "./components/BrandSettings/BrandSettings";
import { useAdminScreen } from "./useAdminScreen";

export function AdminScreen() {
  const { data, error } = useAdminScreen();
  if (!data) return <Unready {...error} />;

  if (!data.admin)
    return <ErrorBanner message="Not authorised to view this page" />;

  return (
    <VStack alignItems="start" gap="4">
      <Heading1 size="lg">Administration</Heading1>

      <BrandSettings />
    </VStack>
  );
}
