"use client";

import { Heading, VStack } from "@chakra-ui/react";

import ErrorBanner from "src/components/site/ErrorBanner";
import { Unready } from "src/components/site/Unready";

import { BrandSettings } from "./components/BrandSettings/BrandSettings";
import { useAdminScreen } from "./useAdminScreen";

export function AdminScreen() {
  const { data, error } = useAdminScreen();
  if (!data) return <Unready {...error} />;

  if (!data.admin)
    return <ErrorBanner message="Not authorised to view this page" />;

  return (
    <VStack alignItems="start" gap={4}>
      <Heading>Administration</Heading>

      <BrandSettings />
    </VStack>
  );
}
