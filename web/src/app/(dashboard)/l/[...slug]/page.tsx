import { Metadata } from "next";
import { cache } from "react";

import { NodeListResult, NodeWithChildren } from "@/api/openapi-schema";
import { nodeGet, nodeListChildren } from "@/api/openapi-server/nodes";
import { ServerRenderedContent } from "@/components/content/ContentRenderer/ServerRenderedContent";
import { getTargetSlug } from "@/components/library/utils";
import { WEB_ADDRESS } from "@/config";
import {
  LibraryPageBlockTypeDirectory,
  parseNodeMetadata,
} from "@/lib/library/metadata";
import { canonicalPage, jsonLD } from "@/lib/metadata/public";
import { getSettings } from "@/lib/settings/settings-server";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";
import { Params, ParamsSchema } from "@/screens/library/library-path";

export type Props = {
  params: Promise<Params>;
};

// Metadata and JSON-LD share one anonymous read for this render.
const getPublicNode = cache(async (slug: string) => {
  try {
    const { data } = await nodeGet(slug, undefined, {
      anonymous: true,
      cache: "no-store",
    });

    if (data.visibility === "published") {
      return data;
    }
  } catch {
    // Public metadata is unavailable if the anonymous read fails.
  }

  return undefined;
});

export default async function Page(props: Props) {
  const { slug } = ParamsSchema.parse(await props.params);

  const targetSlug = getTargetSlug(slug);

  if (!targetSlug) {
    // NOTE: This state is probably not possible to reach due to the params.
    throw new Error("Library page not found");
  }

  const { data } = await nodeGet(targetSlug, undefined, {
    cache: "no-store",
    next: {
      tags: ["library", "node"],
      revalidate: 1,
    },
  });

  // NOTE: A waterfall request which can probably be avoided by fetching the
  // subtree. However subtrees do not currently support property filtering or
  // sorting so this may need a new API endpoint or a parameter for nodeGet.
  const children = await maybeGetChildren(data);
  const publicNode =
    data.visibility === "published"
      ? await getPublicNode(targetSlug)
      : undefined;
  const structuredData = publicNode && {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: publicNode.name,
    description: publicNode.description,
    url: canonicalPage(`/l/${encodeURIComponent(publicNode.slug)}`),
    datePublished: publicNode.createdAt,
    dateModified: publicNode.updatedAt,
  };
  const breadcrumbs = publicNode && {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      {
        "@type": "ListItem",
        position: 1,
        name: "Library",
        item: canonicalPage("/l"),
      },
      ...[...publicNode.ancestors, publicNode].map((node, index) => ({
        "@type": "ListItem",
        position: index + 2,
        name: node.name,
        item: canonicalPage(`/l/${encodeURIComponent(node.slug)}`),
      })),
    ],
  };

  return (
    <>
      {structuredData && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: jsonLD([structuredData, breadcrumbs]),
          }}
        />
      )}
      <ServerRenderedContent values={[data.content]}>
        <LibraryPageScreen node={data} childNodes={children} />
      </ServerRenderedContent>
    </>
  );
}

export async function generateMetadata(props: Props): Promise<Metadata> {
  const params = ParamsSchema.safeParse(await props.params);
  const targetSlug = params.success
    ? getTargetSlug(params.data.slug)
    : undefined;
  const data = targetSlug ? await getPublicNode(targetSlug) : undefined;

  if (!data) {
    return {
      title: "Page not found",
      description: "The page you are looking for does not exist.",
      robots: { index: false, follow: false },
    };
  }

  const settings = await getSettings();
  const url = canonicalPage(`/l/${encodeURIComponent(data.slug)}`);

  return {
    title: `${data.name} | ${settings.title}`,
    description: data.description,
    alternates: { canonical: url },
    openGraph: {
      type: data.content ? "article" : "website",
      title: data.name,
      description: data.description,
      url,
      // NOTE: Massive hack because Next.js still hasn't fixed a bug with
      // catch-all routes and opengraph-image route handlers.
      images: [
        `${WEB_ADDRESS}/l/og?slug=${encodeURIComponent(data.slug)}&t=${encodeURIComponent(data.updatedAt)}`,
      ],
    },
  };
}

async function maybeGetChildren(
  node: NodeWithChildren,
): Promise<NodeListResult | undefined> {
  // NOTE: This is confusingly inverted, the reason is, this flag dictates that
  // the children of this node are hidden in the *sidebar* (technically speaking
  // hidden in the tree traversal API) but not in the *page* view. So when this
  // is false, the nodes are displayed in the sidebar and not in the page, and
  // thus are not fetched alongside the node when rendering the node's page.
  if (!node.hide_child_tree) {
    return;
  }

  const directory = parseNodeMetadata(node.meta).layout?.blocks.find(
    (b): b is LibraryPageBlockTypeDirectory => b.type === "directory",
  );
  if (!directory) {
    return;
  }

  const { data } = await nodeListChildren(node.slug);
  return data;
}
