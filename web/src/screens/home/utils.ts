import { useRouter } from "next/router";
import { WEB_ADDRESS } from "src/config";
import { z } from "zod";

export function getPermalinkForThread(slug: string) {
  return `${WEB_ADDRESS}/t/${slug}`;
}

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export function useQueryParameters() {
  const { query } = useRouter();

  const { category } = QuerySchema.parse(query);

  return { category };
}
