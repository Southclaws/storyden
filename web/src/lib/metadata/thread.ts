import { cache } from "react";
import "server-only";

import { threadGet } from "@/api/openapi-server/threads";

import { isNotPublic, readPublicContent } from "./public";

// Share the viewer's read between the document title and the page body.
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

export const getPublicThread = cache(async (slug: string, page: number) => {
  const viewer = await getThread(slug, page);
  if (viewer?.visibility !== "published") {
    return { status: "not-public" } as const;
  }

  return readPublicContent(async () => {
    const { data } = await threadGet(
      slug,
      { page: String(page) },
      {
        anonymous: true,
      },
    );
    if (
      data.visibility === "published" &&
      page <= Math.max(1, data.replies.total_pages)
    ) {
      return data;
    }

    return undefined;
  });
});
