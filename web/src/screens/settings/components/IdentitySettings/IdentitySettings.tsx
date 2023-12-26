import { useSession } from "src/auth";
import { Unready } from "src/components/site/Unready";
import { Heading2 } from "src/theme/components/Heading/Index";
import { Input } from "src/theme/components/Input";

import { SettingsSection } from "../SettingsSection/SettingsSection";

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
