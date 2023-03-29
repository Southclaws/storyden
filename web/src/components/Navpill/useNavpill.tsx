import { useOutsideClick } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { RefObject, useRef, useState } from "react";
import { useCategoryList } from "src/api/openapi/categories";
import { useAuthProvider } from "src/auth/useAuthProvider";
import { z } from "zod";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export function useNavpill() {
  const { query } = useRouter();

  // TODO: Check if this is the correct way to handle refs in strict TypeScript.
  const ref = useRef<HTMLDivElement>() as RefObject<HTMLDivElement>;

  const { category } = QuerySchema.parse(query);

  const { data, error } = useCategoryList();

  const { account } = useAuthProvider();
  const [isExpanded, setExpanded] = useState(false);

  useOutsideClick({
    ref: ref,
    handler: () => setExpanded(false),
  });

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
    ref,
  };
}
