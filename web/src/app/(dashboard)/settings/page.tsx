import { SettingsScreen } from "@/screens/settings/SettingsScreen";

import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  try {
    const settings = await getSettings();
    return <SettingsScreen initialSettings={settings} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
