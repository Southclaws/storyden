import { Heading, VStack } from "@chakra-ui/react";

import { AuthMethodSettings } from "./components/AuthMethodSettings/AuthMethodSettings";
import { IdentitySettings } from "./components/IdentitySettings/IdentitySettings";

export function SettingsScreen() {
  return (
    <VStack alignItems="start" gap={4}>
      <Heading>Settings</Heading>

      <IdentitySettings />

      <AuthMethodSettings />
    </VStack>
  );
}
