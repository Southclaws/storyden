import { useRouter } from "next/router";
import { useState } from "react";
import { useCategoryList } from "src/api/openapi/categories";
import { useAuthProvider } from "src/auth/useAuthProvider";
import { z } from "zod";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export function useNavigation() {
  const { query } = useRouter();

  const { category } = QuerySchema.parse(query);

  const { data, error } = useCategoryList();

  const { account } = useAuthProvider();
  const [isExpanded, setExpanded] = useState(false);

  const onExpand = () => {
    setExpanded(!isExpanded);
  };

  return {
    categories: data?.categories ?? [],
    error,
    isAuthenticated: !!account,
    isExpanded,
    onExpand,
    category,
  };
}
