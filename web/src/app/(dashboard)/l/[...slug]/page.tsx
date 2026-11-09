import { Metadata } from "next";
import { notFound } from "next/navigation";

import { NodeListResult, NodeWithChildren } from "@/api/openapi-schema";
import { nodeListChildren } from "@/api/openapi-server/nodes";
import { ServerRenderedContent } from "@/components/content/ContentRenderer/ServerRenderedContent";
import { getTargetSlug } from "@/components/library/utils";
import {
  LibraryPageBlockTypeDirectory,
  parseNodeMetadata,
} from "@/lib/library/metadata";
import { getLibraryNode } from "@/lib/metadata/library";
import { canonicalPage, jsonLD, socialImage } from "@/lib/metadata/public";
import { getSettings } from "@/lib/settings/settings-server";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";
import { Params, ParamsSchema } from "@/screens/library/library-path";

export type Props = {
  params: Promise<Params>;
};

export default async function Page(props: Props) {
  const { slug } = ParamsSchema.parse(await props.params);

  const targetSlug = getTargetSlug(slug);

  if (!targetSlug) {
    // NOTE: This state is probably not possible to reach due to the params.
    throw new Error("Library page not found");
  }

  const data = await getLibraryNode(targetSlug);
  if (!data) {
    notFound();
  }

  // NOTE: A waterfall request which can probably be avoided by fetching the
  // subtree. However subtrees do not currently support property filtering or
  // sorting so this may need a new API endpoint or a parameter for nodeGet.
  const children = await maybeGetChildren(data);

  const isPublished = data.visibility === "published";
  const structuredData = isPublished && {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: data.name,
    description: data.description,
    url: canonicalPage(`/l/${encodeURIComponent(data.slug)}`),
    datePublished: data.createdAt,
    dateModified: data.updatedAt,
  };
  const breadcrumbs = isPublished && {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      {
        "@type": "ListItem",
        position: 1,
        name: "Library",
        item: canonicalPage("/l"),
      },
      ...[...data.ancestors, data].map((node, index) => ({
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
  if (!targetSlug) {
    return { title: "Page not found", robots: { index: false, follow: false } };
  }

  const node = await getLibraryNode(targetSlug);
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

  const url = canonicalPage(`/l/${encodeURIComponent(node.slug)}`);

  return {
    title,
    description: node.description,
    alternates: { canonical: url },
    openGraph: {
      type: node.content ? "article" : "website",
      title: node.name,
      description: node.description,
      url,
      // Catch-all routes cannot reliably use an opengraph-image handler.
      images: [
        socialImage(
          `/l/og?slug=${encodeURIComponent(node.slug)}&t=${encodeURIComponent(node.updatedAt)}`,
          node.name,
        ),
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
