import { cache } from "react";
import "server-only";

import { nodeGet } from "@/api/openapi-server/nodes";

import { isNotPublic, readPublicContent } from "./public";

// Share the incoming request's response between metadata and the page body.
export const getLibraryNode = cache(async (slug: string) => {
  try {
    const { data } = await nodeGet(slug);
    return data;
  } catch (error) {
    if (!isNotPublic(error)) {
      throw error;
    }

    return undefined;
  }
});

export const getPublicNode = cache(
  async (slug: string, signal?: AbortSignal) => {
    return readPublicContent(async () => {
      const { data } = await nodeGet(slug, undefined, {
        anonymous: true,
        signal,
      });
      if (data.visibility === "published") {
        return data;
      }

      return undefined;
    });
  },
);
