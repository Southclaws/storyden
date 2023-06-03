import { useParams } from "next/navigation";
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
  const params = useParams();

  const { category } = QuerySchema.parse(params);

  return { category };
}
