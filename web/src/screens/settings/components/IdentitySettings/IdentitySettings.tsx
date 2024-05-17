import { useSession } from "src/auth";
import { Unready } from "src/components/site/Unready";

import { SettingsSection } from "../SettingsSection/SettingsSection";

import { Input } from "@/components/ui/input";
import { Heading2 } from "@/components/ui/typography-heading";

export function IdentitySettings() {
  const account = useSession();

  if (!account) return <Unready />;

  return (
    <SettingsSection>
      <Heading2 size="md">Identity</Heading2>

      <p>You cannot yet change your handle but this is coming soon!</p>

      <Input disabled placeholder="@handle" value={account.handle} />
    </SettingsSection>
  );
}
