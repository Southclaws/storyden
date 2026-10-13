import { categoryGet } from "@/api/openapi-server/categories";
import { threadList } from "@/api/openapi-server/threads";
import { UnreadyBanner } from "@/components/site/Unready";
import { canonicalPage, jsonLD } from "@/lib/metadata/public";
import * as z from "@/lib/validation";
import { CategoryScreen } from "@/screens/category/CategoryScreen";

type Props = {
  params: Promise<{
    category: string;
  }>;
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

export default async function Page({ params, searchParams }: Props) {
  const { category: slug } = await params;
  let page;
  let categoryData;
  let threadListData;

  try {
    ({ page } = QuerySchema.parse(await searchParams));

    ({ data: categoryData } = await categoryGet(slug));

    ({ data: threadListData } = await threadList({
      categories: [slug],
      page: page?.toString(),
    }));
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }

  let structuredData;

  try {
    const [category, threads] = await Promise.all([
      categoryGet(slug, { anonymous: true, cache: "no-store" }),
      threadList(
        {
          categories: [slug],
          page: String(page ?? 1),
          visibility: ["published"],
        },
        { anonymous: true, cache: "no-store" },
      ),
    ]);

    if ((page ?? 1) <= Math.max(1, threads.data.total_pages)) {
      structuredData = {
        "@context": "https://schema.org",
        "@type": "CollectionPage",
        name: category.data.name,
        description: category.data.description,
        url: canonicalPage(
          `/d/${encodeURIComponent(category.data.slug)}`,
          page,
        ),
        mainEntity: {
          "@type": "ItemList",
          itemListElement: threads.data.threads
            .filter(
              (thread) =>
                thread.visibility === "published" &&
                threadListData.threads.some(
                  (visibleThread) => visibleThread.id === thread.id,
                ),
            )
            .map((thread, index) => ({
              "@type": "ListItem",
              position: index + 1,
              name: thread.title,
              url: canonicalPage(`/t/${encodeURIComponent(thread.slug)}`),
            })),
        },
      };
    }
  } catch {
    // The signed-in category can still render without public structured data.
  }

  return (
    <>
      {structuredData && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: jsonLD(structuredData) }}
        />
      )}
      <CategoryScreen
        initialPage={page ?? 1}
        slug={slug}
        initialCategory={categoryData}
        initialThreadList={threadListData}
      />
    </>
  );
}
