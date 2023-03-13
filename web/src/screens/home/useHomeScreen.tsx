import { useRouter } from "next/router";
import { useThreadList } from "src/api/openapi/threads";
import { z } from "zod";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export function useHomeScreen() {
  const { query } = useRouter();

  const { category } = QuerySchema.parse(query);

  const threads = useThreadList({
    categories: category ? [category] : undefined,
  });

  return threads;
}
