import z from "zod";

import { threadList } from "@/api/openapi-server/threads";
import { UnreadyBanner } from "@/components/site/Unready";
import { categoryListCached } from "@/lib/category/server-category-list";
import { CategoryIndexScreen } from "@/screens/category/CategoryIndexScreen";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});
type Query = z.infer<typeof QuerySchema>;

export default async function Page({ searchParams }: Props) {
  try {
    const { page } = QuerySchema.parse(await searchParams);

    const { data: categories } = await categoryListCached();
    const { data: threads } = await threadList({
      page: page?.toString() ?? "1",
      // NOTE: The string "null" is a special case that yields all threads that
      // do not have a category. Why a string null and not just null? Because
      // the OpenAPI generators don't play well with query parameters that use
      // an algebraic type like string[] | null. We even tried using an array
      // with a string like string[] | string where the second string is set
      // to filter against the regex "^null$" but that also didn't work. So,
      // simply dumb solution: backend checks for a list item of "null". Lol.
      categories: ["null"],
    });

    return (
      <CategoryIndexScreen
        layout={"list"}
        threadListMode="uncategorised"
        showQuickShare={true}
        initialCategoryList={categories}
        initialThreadList={threads}
        initialThreadListPage={page}
        paginationBasePath="/d"
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
