import { useOutsideClick } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { RefObject, useRef, useState } from "react";
import { useGetInfo } from "src/api/openapi/misc";
import { useAuthProvider } from "src/auth/useAuthProvider";
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
  const overlayRef = useRef<HTMLDivElement>() as RefObject<HTMLDivElement>;
  const { account } = useAuthProvider();
  const [isExpanded, setExpanded] = useState(false);

  const { category } = QuerySchema.parse(query);

  const title = infoResult?.title ?? "Storyden";

  useOutsideClick({
    ref: overlayRef,
    handler: () => setExpanded(false),
  });

  const onExpand = () => {
    setExpanded(!isExpanded);
  };

  return {
    title,
    isAuthenticated: !!account,
    isExpanded,
    onExpand,
    category,
    overlayRef,
  };
}
