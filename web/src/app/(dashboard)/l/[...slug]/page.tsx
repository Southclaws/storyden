import { Metadata } from "next";

import { NodeListResult, NodeWithChildren } from "@/api/openapi-schema";
import { nodeGet, nodeListChildren } from "@/api/openapi-server/nodes";
import { getTargetSlug } from "@/components/library/utils";
import { WEB_ADDRESS } from "@/config";
import {
  LibraryPageBlockTypeTable,
  parseNodeMetadata,
} from "@/lib/library/metadata";
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

  return <LibraryPageScreen node={data} childNodes={children} />;
}

export async function generateMetadata(props: Props) {
  try {
    const { slug } = ParamsSchema.parse(await props.params);

    const targetSlug = getTargetSlug(slug);

    if (!targetSlug) {
      // NOTE: This state is probably not possible to reach due to the params.
      throw new Error("Library page not found");
    }

    const settings = await getSettings();

    const { data } = await nodeGet(targetSlug);

    return {
      title: `${data.name} | ${settings.title}`,
      description: data.description,
      openGraph: {
        // NOTE: Massive hack because Next.js still hasn't fixed a bug with
        // catch-all routes and opengraph-image route handlers.
        images: [`${WEB_ADDRESS}/l/og?slug=${targetSlug}&t=${data.updatedAt}`],
      },
    } satisfies Metadata;
  } catch (e) {
    return {
      title: "Page not found",
      description: "The page you are looking for does not exist.",
    };
  }
}

async function maybeGetChildren(
  node: NodeWithChildren,
): Promise<NodeListResult | undefined> {
  if (!node.hide_child_tree) {
    return;
  }

  const table = parseNodeMetadata(node.meta).layout?.blocks.find(
    (b): b is LibraryPageBlockTypeTable => b.type === "table",
  );
  if (!table) {
    return;
  }

  const { data } = await nodeListChildren(node.slug);
  return data;
}
