import { useAccountSession } from "@/auth";
import { EmailSettings } from "@/components/settings/EmailSettings/EmailSettings";
import { UnreadyBanner } from "@/components/site/Unready";

export function MemberEmailSettingsScreen() {
  const { data, error } = useAccountSession();
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <EmailSettings account={data} />;
}
