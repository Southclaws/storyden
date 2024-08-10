"use client";

import { useGetInfo } from "src/api/openapi-client/misc";
import { useSession } from "src/auth";

// NOTE: Everything that involves data fetching here has a suitable fallback.
// Most of the time, components will render <Unready /> but this is the data for
// the navigation so it's a bit more important that we show something always.

export function useNavigation() {
  const { data: infoResult } = useGetInfo();
  const session = useSession();

  const title = infoResult?.title ?? "Storyden";

  return {
    isAdmin: session?.admin ?? false,
    title,
  };
}
