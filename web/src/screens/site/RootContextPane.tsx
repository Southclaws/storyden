import { getServerSession } from "@/auth/server-session";
import { SiteContextPane } from "@/components/site/SiteContextPane/SiteContextPane";
import { getInfo } from "@/utils/info";

export async function RootContextPane() {
  const session = await getServerSession();
  const info = await getInfo();

  return <SiteContextPane info={info} session={session} />;
}
