import { useOutsideClick } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { RefObject, useRef, useState } from "react";
import { useCategoryList } from "src/api/openapi/categories";
import { useGetInfo } from "src/api/openapi/misc";
import { useAuthProvider } from "src/auth/useAuthProvider";
import { z } from "zod";

export const QuerySchema = z.object({
  category: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export function useNavigation() {
  const { query } = useRouter();
  const { data: infoResult, error: infoError } = useGetInfo();
  const { data: categoriesResult, error: categoriesError } = useCategoryList();
  const overlayRef = useRef<HTMLDivElement>() as RefObject<HTMLDivElement>;
  const { account } = useAuthProvider();
  const [isExpanded, setExpanded] = useState(false);

  const { category } = QuerySchema.parse(query);

  const error = infoError ?? categoriesError ?? undefined;
  const categories = categoriesResult?.categories ?? [];
  const title = infoResult?.title ?? "be";

  useOutsideClick({
    ref: overlayRef,
    handler: () => setExpanded(false),
  });

  const onExpand = () => {
    setExpanded(!isExpanded);
  };

  return {
    categories,
    title,
    error,
    isAuthenticated: !!account,
    isExpanded,
    onExpand,
    category,
    overlayRef,
  };
}
