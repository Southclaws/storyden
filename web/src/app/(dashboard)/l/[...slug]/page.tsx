import { Metadata } from "next";

import { nodeGet } from "@/api/openapi-server/nodes";
import { getTargetSlug } from "@/components/library/utils";
import { UnreadyBanner } from "@/components/site/Unready";
import { WEB_ADDRESS } from "@/config";
import { getSettings } from "@/lib/settings/settings-server";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";
import { Params, ParamsSchema } from "@/screens/library/library-path";

export type Props = {
  params: Promise<Params>;
};

export default async function Page(props: Props) {
  try {
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

    return <LibraryPageScreen node={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
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
