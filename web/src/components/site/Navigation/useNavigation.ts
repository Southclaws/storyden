"use client";

import { useParams } from "next/navigation";

import { useSession } from "src/auth";

import { useSettings } from "@/lib/settings/settings-client";

// NOTE: Everything that involves data fetching here has a suitable fallback.
// Most of the time, components will render <Unready /> but this is the data for
// the navigation so it's a bit more important that we show something always.

export function useNavigation() {
  const { settings } = useSettings();
  const session = useSession();
  const { slug } = useParams();

  const nodeSlug = slug?.[0];

  const title = settings?.title ?? "Storyden";

  return {
    isAdmin: session?.admin ?? false,
    title,
    nodeSlug,
  };
}
