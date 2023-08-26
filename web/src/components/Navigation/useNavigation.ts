"use client";

import { useSearchParams } from "next/navigation";
import { PropertyValue } from "styled-system/types/prop-type";
import { z } from "zod";

import { useGetInfo } from "src/api/openapi/misc";
import { useSession } from "src/auth";

// The sidebar width is shared between two components which must use the exact
// same values. The reason for this is the sidebar is position: fixed, which
// means it cannot inherit the width from the parent since its true parent is
// the viewport. To get around this, the default layout positions an empty box
// to the left of the viewport to push the content right and then the actual
// sidebar is rendered on top of this with the same width configuration.
export const SIDEBAR_WIDTH: PropertyValue<"width"> = {
  md: "1/4",
  lg: "1/3",
};

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
