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
import { getLibraryNode, getPublicNode } from "@/lib/metadata/library";
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
  const [children, result] = await Promise.all([
    maybeGetChildren(data),
    data.visibility === "published" ? getPublicNode(targetSlug) : undefined,
  ]);
  const publicNode = result?.status === "public" ? result.data : undefined;

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

  const result = await getPublicNode(targetSlug);
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
