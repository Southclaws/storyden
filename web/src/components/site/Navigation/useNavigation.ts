"use client";

import { useSearchParams } from "next/navigation";
import { z } from "zod";

import { useGetInfo } from "src/api/openapi/misc";
import { useSession } from "src/auth";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

// NOTE: Everything that involves data fetching here has a suitable fallback.
// Most of the time, components will render <Unready /> but this is the data for
// the navigation so it's a bit more important that we show something always.

export function useNavigation() {
  const query = useSearchParams();
  const { data: infoResult } = useGetInfo();
  const session = useSession();

  const { category } = QuerySchema.parse(query);

  const title = infoResult?.title ?? "Storyden";

  return {
    isAdmin: session?.admin ?? false,
    title,
    category,
  };
}
