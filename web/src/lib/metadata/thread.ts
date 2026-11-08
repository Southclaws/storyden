import type { Metadata } from "next";
import { cookies } from "next/headers";
import { cache } from "react";
import "server-only";

import { threadGet } from "@/api/openapi-server/threads";
import { getSettings } from "@/lib/settings/settings-server";

import {
  canonicalPage,
  isNotPublic,
  pageLinks,
  readPublicContent,
  socialImage,
} from "./public";

// Share the viewer's read between the document title and the page body.
export const getThread = cache(async (slug: string, page: number) => {
  try {
    const { data } = await threadGet(
      slug,
      { page: String(page) },
      { cache: "no-store" },
    );
    return data;
  } catch (error) {
    if (!isNotPublic(error)) {
      throw error;
    }
  }

  return undefined;
});

// Public metadata and JSON-LD must use the anonymous view of the thread.
export const getPublicThread = cache(async (slug: string, page: number) => {
  return readPublicContent(async () => {
    const session = (await cookies()).get("storyden-session");
    let data;

    if (session) {
      const response = await threadGet(
        slug,
        { page: String(page) },
        { anonymous: true, cache: "no-store" },
      );
      data = response.data;
    } else {
      // The viewer's read is already anonymous for crawlers and guests.
      data = await getThread(slug, page);
    }

    if (
      data?.visibility === "published" &&
      page <= Math.max(1, data.replies.total_pages)
    ) {
      return data;
    }

    return undefined;
  });
});

export const getThreadMetadata = cache(
  async (slug: string, page: number): Promise<Metadata> => {
    const thread = await getThread(slug, page);
    if (!thread) {
      return {
        title: "Thread Not Found",
        robots: { index: false, follow: false },
      };
    }

    const settings = await getSettings();
    const title = `${thread.title} | ${settings.title}`;
    if (thread.visibility !== "published") {
      return { title, robots: { index: false, follow: false } };
    }

    const result = await getPublicThread(slug, page);
    if (result.status === "not-public") {
      return { title, robots: { index: false, follow: false } };
    }

    if (result.status === "unavailable") {
      return { title };
    }

    const publicThread = result.data;
    const path = `/t/${encodeURIComponent(publicThread.slug)}`;
    const url = canonicalPage(path, page);
    const totalPages = Math.max(1, publicThread.replies.total_pages);

    return {
      title,
      description: publicThread.description,
      alternates: { canonical: url },
      pagination: pageLinks(path, page, totalPages),
      openGraph: {
        images: [socialImage(undefined, settings.title)],
        type: "article",
        title: publicThread.title,
        description: publicThread.description,
        url,
        publishedTime: publicThread.createdAt,
        modifiedTime: publicThread.updatedAt,
      },
    };
  },
);
