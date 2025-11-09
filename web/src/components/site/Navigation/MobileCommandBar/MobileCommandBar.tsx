import { Suspense } from "react";

import { getServerSession } from "@/auth/server-session";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { MobileCommandBarClient } from "./MobileCommandBarClient";

export async function MobileCommandBar() {
  const session = await getServerSession();

  return (
    <MobileCommandBarClient hasAccount={!!session}>
      <Suspense>
        <ContentNavigationList />
      </Suspense>
    </MobileCommandBarClient>
  );
}
