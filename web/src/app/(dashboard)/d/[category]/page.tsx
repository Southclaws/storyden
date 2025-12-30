import { z } from "zod";

import { UnreadyBanner } from "src/components/site/Unready";

import { categoryGet } from "@/api/openapi-server/categories";
import { threadList } from "@/api/openapi-server/threads";
import { CategoryScreen } from "@/screens/category/CategoryScreen";

type Props = {
  params: Promise<{
    category: string;
  }>;
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});
type Query = z.infer<typeof QuerySchema>;

export default async function Page({ params, searchParams }: Props) {
  try {
    const { category: slug } = await params;
    const { page } = QuerySchema.parse(await searchParams);

    const { data: categoryData } = await categoryGet(slug);

    const { data: threadListData } = await threadList({
      categories: [slug],
      page: page?.toString(),
    });

    return (
      <CategoryScreen
        initialPage={page ?? 1}
        slug={slug}
        initialCategory={categoryData}
        initialThreadList={threadListData}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
