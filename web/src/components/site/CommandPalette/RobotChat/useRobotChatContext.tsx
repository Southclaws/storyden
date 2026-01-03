"use client";

import { useParams, usePathname } from "next/navigation";
import { useCallback, useRef } from "react";

import { nodeGet } from "@/api/openapi-client/nodes";
import { profileGet } from "@/api/openapi-client/profiles";
import { threadGet } from "@/api/openapi-client/threads";
import { RobotChatContext } from "@/api/robots";

export function useRobotPageContext(): () => Promise<RobotChatContext> {
  const pathname = usePathname();
  const params = useParams();

  const pathnameRef = useRef(pathname);
  const paramsRef = useRef(params);

  pathnameRef.current = pathname;
  paramsRef.current = params;

  const getContext = useCallback(async (): Promise<RobotChatContext> => {
    const currentPathname = pathnameRef.current;
    const currentParams = paramsRef.current;

    if (currentPathname?.startsWith("/t/") && currentParams["slug"]) {
      try {
        const slugParam = currentParams["slug"];
        const slug = Array.isArray(slugParam) ? slugParam[0] : slugParam;
        if (!slug) return {};

        const thread = await threadGet(slug);
        return {
          datagraph_item: {
            id: thread.id,
            slug: thread.slug,
            kind: "thread",
          },
        };
      } catch (e) {
        console.error("Failed to fetch thread for context:", e);
      }
    }

    if (currentPathname?.startsWith("/l/") && currentParams["slug"]) {
      try {
        const slugParam = currentParams["slug"];
        const slug = Array.isArray(slugParam) ? slugParam.join("/") : slugParam;
        if (!slug) return {};

        const node = await nodeGet(slug);
        return {
          datagraph_item: {
            id: node.id,
            slug: node.slug,
            kind: "node",
          },
        };
      } catch (e) {
        console.error("Failed to fetch node for context:", e);
      }
    }

    if (currentPathname?.startsWith("/m/") && currentParams["handle"]) {
      try {
        const handleParam = currentParams["handle"];
        const handle = Array.isArray(handleParam)
          ? handleParam[0]
          : handleParam;
        if (!handle) return {};

        const profile = await profileGet(handle);
        return {
          datagraph_item: {
            id: profile.id,
            slug: profile.handle,
            kind: "profile",
          },
        };
      } catch (e) {
        console.error("Failed to fetch profile for context:", e);
      }
    }

    if (currentPathname === "/") {
      return { page_type: "Index page" };
    }

    if (currentPathname?.startsWith("/settings")) {
      return { page_type: "Settings page" };
    }

    if (currentPathname?.startsWith("/admin")) {
      return { page_type: "Admin page" };
    }

    if (currentPathname?.startsWith("/search")) {
      return { page_type: "Search page" };
    }

    if (currentPathname?.startsWith("/notifications")) {
      return { page_type: "Notifications page" };
    }

    if (currentPathname?.startsWith("/d/")) {
      return { page_type: "Directory page" };
    }

    if (currentPathname?.startsWith("/c/")) {
      return { page_type: "Category page" };
    }

    if (currentPathname?.startsWith("/tags")) {
      return { page_type: "Tags page" };
    }

    if (currentPathname?.startsWith("/links")) {
      return { page_type: "Links page" };
    }

    if (currentPathname?.startsWith("/new")) {
      return { page_type: "New post page" };
    }

    if (currentPathname?.startsWith("/drafts")) {
      return { page_type: "Drafts page" };
    }

    if (currentPathname?.startsWith("/queue")) {
      return { page_type: "Queue page" };
    }

    if (currentPathname?.startsWith("/reports")) {
      return { page_type: "Reports page" };
    }

    if (currentPathname?.startsWith("/roles")) {
      return { page_type: "Roles page" };
    }

    return {};
  }, []);

  return getContext;
}
