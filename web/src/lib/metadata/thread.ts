import { cache } from "react";
import "server-only";

import { threadGet } from "@/api/openapi-server/threads";

import { isNotPublic } from "./public";

// Share the incoming request's response between metadata and the page body.
export const getThread = cache(async (slug: string, page: number) => {
  try {
    const { data } = await threadGet(slug, { page: String(page) });
    return data;
  } catch (error) {
    if (!isNotPublic(error)) {
      throw error;
    }

    return undefined;
  }
});
