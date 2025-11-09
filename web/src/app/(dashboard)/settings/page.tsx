import { redirect } from "next/navigation";

import { SettingsScreen } from "src/screens/settings/SettingsScreen";

import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

export default async function Page() {
  const session = await getServerSession();
  if (!session) {
    redirect("/login");
  }

  try {
    const settings = await getSettings();
    return <SettingsScreen initialSettings={settings} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
