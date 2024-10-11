import { useParams, usePathname, useRouter } from "next/navigation";

import { LibraryPath, Params, ParamsSchema } from "./library-path";

export function useLibraryPath() {
  const params = useParams<Params>();

  const parsed = ParamsSchema.safeParse(params);

  if (!parsed.success) {
    return [];
  }

  const { slug } = parsed.data;

  const path = slug as LibraryPath;

  return path;
}
