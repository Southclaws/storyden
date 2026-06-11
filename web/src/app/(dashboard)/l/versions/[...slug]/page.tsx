import { Metadata } from "next";

import { nodeGet, nodeVersionList } from "@/api/openapi-server/nodes";
import { getTargetSlug } from "@/components/library/utils";
import { getSettings } from "@/lib/settings/settings-server";
import { LibraryPageVersionHistoryScreen } from "@/screens/library/LibraryPageVersionHistoryScreen";
import { Params, ParamsSchema } from "@/screens/library/library-path";

export type Props = {
  params: Promise<Params>;
};

export default async function Page(props: Props) {
  const { slug } = ParamsSchema.parse(await props.params);
  const targetSlug = getTargetSlug(slug);

  if (!targetSlug) {
    throw new Error("Library page not found");
  }

  const [{ data: node }, { data: versions }] = await Promise.all([
    nodeGet(targetSlug, undefined, {
      cache: "no-store",
      next: {
        tags: ["library", "node"],
        revalidate: 1,
      },
    }),
    nodeVersionList(targetSlug, undefined, {
      cache: "no-store",
      next: {
        tags: ["library", "node", "versions"],
        revalidate: 1,
      },
    }),
  ]);

  return (
    <LibraryPageVersionHistoryScreen
      node={node}
      versions={versions.versions}
      libraryPath={slug}
    />
  );
}

export async function generateMetadata(props: Props) {
  try {
    const { slug } = ParamsSchema.parse(await props.params);
    const targetSlug = getTargetSlug(slug);

    if (!targetSlug) {
      throw new Error("Library page not found");
    }

    const [settings, { data }] = await Promise.all([
      getSettings(),
      nodeGet(targetSlug),
    ]);

    return {
      title: `Version history for ${data.name} | ${settings.title}`,
      description: `Edit history for ${data.name}.`,
    } satisfies Metadata;
  } catch (e) {
    return {
      title: "Version history",
      description: "Library page version history.",
    };
  }
}
