import type { Metadata } from "next";
import { cookies } from "next/headers";
import { cache } from "react";
import "server-only";

import { nodeGet } from "@/api/openapi-server/nodes";
import { getSettings } from "@/lib/settings/settings-server";

import {
  canonicalPage,
  isNotPublic,
  readPublicContent,
  socialImage,
} from "./public";

export const getLibraryNode = cache(async (slug: string) => {
  try {
    const { data } = await nodeGet(slug, undefined, { cache: "no-store" });
    return data;
  } catch (error) {
    if (!isNotPublic(error)) {
      throw error;
    }
  }

  return undefined;
});

export const getPublicNode = cache(async (slug: string) => {
  return readPublicContent(async () => {
    const session = (await cookies()).get("storyden-session");
    let data;

    if (session) {
      const response = await nodeGet(slug, undefined, {
        anonymous: true,
        cache: "no-store",
      });
      data = response.data;
    } else {
      data = await getLibraryNode(slug);
    }

    if (data?.visibility === "published") {
      return data;
    }

    return undefined;
  });
});

export const getLibraryMetadata = cache(
  async (slug: string): Promise<Metadata> => {
    const node = await getLibraryNode(slug);
    if (!node) {
      return {
        title: "Page not found",
        robots: { index: false, follow: false },
      };
    }

    const settings = await getSettings();
    const title = `${node.name} | ${settings.title}`;
    if (node.visibility !== "published") {
      return { title, robots: { index: false, follow: false } };
    }

    const result = await getPublicNode(slug);
    if (result.status === "not-public") {
      return { title, robots: { index: false, follow: false } };
    }

    if (result.status === "unavailable") {
      return { title };
    }

    const publicNode = result.data;
    const url = canonicalPage(`/l/${encodeURIComponent(publicNode.slug)}`);

    return {
      title,
      description: publicNode.description,
      alternates: { canonical: url },
      openGraph: {
        type: publicNode.content ? "article" : "website",
        title: publicNode.name,
        description: publicNode.description,
        url,
        // Catch-all routes cannot reliably use an opengraph-image handler.
        images: [
          socialImage(
            `/l/og?slug=${encodeURIComponent(publicNode.slug)}&t=${encodeURIComponent(publicNode.updatedAt)}`,
            publicNode.name,
          ),
        ],
      },
    };
  },
);
