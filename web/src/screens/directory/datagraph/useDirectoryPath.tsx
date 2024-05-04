import { pull } from "lodash";
import { useParams } from "next/navigation";

import { DirectoryPath, Params, ParamsSchema } from "./directory-path";

export function useDirectoryPath() {
  const params = useParams<Params>();

  const parsed = ParamsSchema.safeParse(params);

  if (!parsed.success) {
    return [];
  }

  const { slug } = parsed.data;

  const cleaned = pull(slug, "new");

  const dp = cleaned as DirectoryPath;

  return dp;
}
