import { cacheLife } from "next/cache";

import { categoryList } from "@/api/openapi-server/categories";

export async function categoryListCached() {
  "use cache";
  cacheLife("minutes");
  return await categoryList({ cache: "no-store" });
}
