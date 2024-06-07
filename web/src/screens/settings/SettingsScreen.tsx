"use client";

import { Heading1 } from "@/components/ui/typography-heading";
import { VStack } from "@/styled-system/jsx";

import { AuthMethodSettings } from "./components/AuthMethodSettings/AuthMethodSettings";
import { IdentitySettings } from "./components/IdentitySettings/IdentitySettings";

export function SettingsScreen() {
  return (
    <VStack alignItems="start" gap="4">
      <Heading1 size="lg">Settings</Heading1>

      <IdentitySettings />

      <AuthMethodSettings />
    </VStack>
  );
}
