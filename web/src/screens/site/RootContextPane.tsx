import { Suspense } from "react";

import { getServerSession } from "@/auth/server-session";
import { SiteContextPane } from "@/components/site/SiteContextPane/SiteContextPane";
import { Unready } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

export async function RootContextPane() {
  return (
    <Suspense fallback={<Unready />}>
      <Suspensed />
    </Suspense>
  );
}

async function Suspensed() {
  const session = await getServerSession();
  const settings = await getSettings();
  return <SiteContextPane initialSettings={settings} session={session} />;
}
