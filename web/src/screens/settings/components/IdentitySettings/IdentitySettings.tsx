import { useSession } from "src/auth";
import { Unready } from "src/components/site/Unready";

import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";

import { SettingsSection } from "../SettingsSection/SettingsSection";

export function IdentitySettings() {
  const account = useSession();

  if (!account) return <Unready />;

  return (
    <SettingsSection>
      <Heading size="md">Identity</Heading>

      <p>You cannot yet change your handle but this is coming soon!</p>

      <Input disabled placeholder="@handle" value={account.handle} />
    </SettingsSection>
  );
}
