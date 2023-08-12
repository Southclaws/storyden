import { Heading, Input, Text } from "@chakra-ui/react";

import { useSession } from "src/auth";
import { Unready } from "src/components/Unready";

import { SettingsSection } from "../SettingsSection/SettingsSection";

export function IdentitySettings() {
  const account = useSession();

  if (!account) return <Unready />;

  return (
    <SettingsSection>
      <Heading size="sm">Identity</Heading>

      <Text>You cannot yet change your handle but this is coming soon!</Text>

      <Input disabled placeholder="@handle" value={account.handle} />
    </SettingsSection>
  );
}
