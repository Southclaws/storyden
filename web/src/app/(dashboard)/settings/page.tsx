import { SettingsScreen } from "src/screens/settings/SettingsScreen";

import { UnreadyBanner } from "@/components/site/Unready";

export default async function Page() {
  try {
    return <SettingsScreen />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
