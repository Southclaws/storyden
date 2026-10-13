import type { Metadata } from "next";

import { threadList } from "@/api/openapi-server/threads";
import { UnreadyBanner } from "@/components/site/Unready";
import { getCategoryThreadListParams } from "@/lib/feed/category";
import { canonicalPage, jsonLD, pageLinks } from "@/lib/metadata/public";
import type { Settings } from "@/lib/settings/settings";
import { getSettings } from "@/lib/settings/settings-server";
import * as z from "@/lib/validation";
import { FeedScreen } from "@/screens/feed/FeedScreen";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z.optional(
    z
      .pipe(z.string().check(z.regex(/^[1-9]\d*$/)), z.transform(Number))
      .check(z.refine(Number.isSafeInteger)),
  ),
});
type Query = z.input<typeof QuerySchema>;

export default async function Page({ searchParams }: Props) {
  let page: number | undefined;
  let settings: Settings;

  try {
    ({ page } = QuerySchema.parse(await searchParams));
    settings = await getSettings();
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }

  const structuredData = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: settings.title,
    description: settings.description,
    url: canonicalPage("/"),
  };

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: jsonLD(structuredData) }}
      />
      <FeedScreen page={page ?? 1} />
    </>
  );
}

export async function generateMetadata({
  searchParams,
}: Props): Promise<Metadata> {
  const query = QuerySchema.safeParse(await searchParams);

  if (!query.success) {
    return { title: "Page not found", robots: { index: false, follow: false } };
  }

  const settings = await getSettings();
  const threadsBlock = settings.metadata.feed.blocks.find(
    (block) => block.type === "threads",
  );
  const page = threadsBlock ? (query.data.page ?? 1) : 1;
  const url = canonicalPage("/", page);
  let pagination: Metadata["pagination"];

  if (threadsBlock) {
    try {
      const { data } = await threadList(
        {
          ...getCategoryThreadListParams(threadsBlock.source, page),
          visibility: ["published"],
        },
        { anonymous: true, cache: "no-store" },
      );
      const totalPages = Math.max(1, data.total_pages);

      if (page > totalPages) {
        return {
          title: "Page not found",
          robots: { index: false, follow: false },
        };
      }

      pagination = pageLinks("/", page, totalPages);
    } catch {
      // Public site metadata remains available when the feed requires login.
      if (page > 1) {
        return {
          title: "Page not found",
          robots: { index: false, follow: false },
        };
      }
    }
  }

  return {
    title: `${settings.title} | ${settings.description}`,
    description: settings.description,
    alternates: { canonical: url },
    pagination,
    openGraph: {
      type: "website",
      title: settings.title,
      description: settings.description,
      url,
    },
  };
}
