"use client";

import { Unready, UnreadyBanner } from "src/components/site/Unready";

import { Heading } from "@/components/ui/heading";
import { VStack } from "@/styled-system/jsx";

import { BrandSettings } from "./components/BrandSettings/BrandSettings";
import { useAdminScreen } from "./useAdminScreen";

export function AdminScreen() {
  const { data, error } = useAdminScreen();
  if (!data) return <Unready error={error} />;

  if (!data.admin)
    return <UnreadyBanner error="Not authorised to view this page" />;

  return (
    <VStack alignItems="start" gap="4">
      <Heading size="lg">Administration</Heading>

      <BrandSettings />
    </VStack>
  );
}
