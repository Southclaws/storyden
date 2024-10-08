"use client";

import { Heading } from "@/components/ui/heading";
import { VStack } from "@/styled-system/jsx";

import { AuthMethodSettings } from "./components/AuthMethodSettings/AuthMethodSettings";
import { IdentitySettings } from "./components/IdentitySettings/IdentitySettings";

export function SettingsScreen() {
  return (
    <VStack alignItems="start" gap="4">
      <Heading size="lg">Settings</Heading>

      <AuthMethodSettings />
    </VStack>
  );
}
