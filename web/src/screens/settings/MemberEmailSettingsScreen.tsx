import { useAccountGet } from "@/api/openapi-client/accounts";
import { EmailSettings } from "@/components/settings/EmailSettings/EmailSettings";
import { UnreadyBanner } from "@/components/site/Unready";

export function MemberEmailSettingsScreen() {
  const { data, error } = useAccountGet();
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <EmailSettings account={data} />;
}
