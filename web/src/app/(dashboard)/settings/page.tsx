import { SettingsScreen } from "src/screens/settings/SettingsScreen";

import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

export default async function Page() {
  try {
    const settings = await getSettings();
    return <SettingsScreen initialSettings={settings} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
