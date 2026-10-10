import type { Metadata } from "next";
import { z } from "zod";

import { UnreadyBanner } from "@/components/site/Unready";
import {
  canonicalPage,
  getAnonymousPageCount,
  jsonLD,
  pageLinks,
  pageNumber,
  websiteSchema,
} from "@/lib/metadata/public";
import { getSettings } from "@/lib/settings/settings-server";
import { FeedScreen } from "@/screens/feed/FeedScreen";

type Props = {
  searchParams: Promise<{ page?: string }>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .regex(/^[1-9]\d*$/)
    .transform(Number)
    .refine(Number.isSafeInteger)
    .optional(),
});
export default async function Page({ searchParams }: Props) {
  try {
    const { page } = QuerySchema.parse(await searchParams);
    const settings = await getSettings();

    return (
      <>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: jsonLD(websiteSchema(settings.title, settings.description)),
          }}
        />
        <FeedScreen page={page ?? 1} />
      </>
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}

export async function generateMetadata({
  searchParams,
}: Props): Promise<Metadata> {
  const settings = await getSettings();
  const { page } = await searchParams;
  if (page !== undefined && pageNumber(page) === undefined) {
    return { title: "Page not found", robots: { index: false, follow: false } };
  }
  const currentPage = pageNumber(page) ?? 1;
  const threadsBlock = settings.metadata.feed.blocks.find(
    (block) => block.type === "threads",
  );
  const canonicalPageNumber = threadsBlock ? page : undefined;
  const url = canonicalPage("/", canonicalPageNumber);
  const count = threadsBlock
    ? await getAnonymousPageCount(
        `/threads?page=${currentPage}${threadsBlock.source === "uncategorised" ? "&categories=null" : ""}`,
      )
    : undefined;
  if (count && currentPage > count) {
    return { title: "Page not found", robots: { index: false, follow: false } };
  }

  return {
    title: `${settings.title} | ${settings.description}`,
    description: settings.description,
    alternates: {
      canonical: url,
      types: {
        "text/markdown": canonicalPage("/index.md", canonicalPageNumber),
      },
    },
    ...(count ? { pagination: pageLinks("/", currentPage, count) } : {}),
    openGraph: {
      type: "website",
      title: settings.title,
      description: settings.description,
      url,
    },
  };
}
