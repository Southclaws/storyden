import { cacheLife } from "next/cache";

import { NodeListParams } from "@/api/openapi-schema";
import { nodeList } from "@/api/openapi-server/nodes";

export async function nodeListCached(params?: NodeListParams) {
  "use cache";
  cacheLife("minutes");
  return await nodeList(params, { cache: "no-store" });
}
