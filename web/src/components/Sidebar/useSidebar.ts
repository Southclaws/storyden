import { useRouter } from "next/router";
import { useCategoryList } from "src/api/openapi/categories";
import { z } from "zod";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export function useSidebar() {
  const { query } = useRouter();
  const { data, error } = useCategoryList();

  const { category } = QuerySchema.parse(query);

  return {
    categories: data?.categories,
    category,
    error,
  };
}
