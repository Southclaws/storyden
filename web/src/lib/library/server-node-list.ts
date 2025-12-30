"use server";

import { cache } from "react";

import { NodeListParams } from "@/api/openapi-schema";
import { nodeList } from "@/api/openapi-server/nodes";

export const nodeListCached = cache(async (params?: NodeListParams) => {
  return await nodeList(params, {
    cache: "default",
  });
});
