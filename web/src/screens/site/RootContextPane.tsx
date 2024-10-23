import { getServerSession } from "@/auth/server-session";
import { SiteContextPane } from "@/components/site/SiteContextPane/SiteContextPane";
import { getSettings } from "@/lib/settings/settings-server";

export async function RootContextPane() {
  const session = await getServerSession();
  const settings = await getSettings();

  return <SiteContextPane initialSettings={settings} session={session} />;
}
