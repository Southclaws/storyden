import { pull } from "lodash";
import { useParams } from "next/navigation";

import { LibraryPath, Params, ParamsSchema } from "./library-path";

export function useLibraryPath() {
  const params = useParams<Params>();

  const parsed = ParamsSchema.safeParse(params);

  if (!parsed.success) {
    return [];
  }

  const { slug } = parsed.data;

  const cleaned = pull(slug, "new");

  const dp = cleaned as LibraryPath;

  return dp;
}
