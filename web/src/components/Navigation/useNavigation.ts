import { useRouter } from "next/router";
import { useGetInfo } from "src/api/openapi/misc";
import { z } from "zod";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

// NOTE: Everything that involves data fetching here has a suitable fallback.
// Most of the time, components will render <Unready /> but this is the data for
// the navigation so it's a bit more important that we show something always.

export function useNavigation() {
  const { query } = useRouter();
  const { data: infoResult } = useGetInfo();

  const { category } = QuerySchema.parse(query);

  const title = infoResult?.title ?? "Storyden";

  return {
    title,
    category,
  };
}
