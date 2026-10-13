import { keyBy } from "lodash/fp";

import { categoryList } from "@/api/openapi-server/categories";
import { datagraphSearch } from "@/api/openapi-server/datagraph";
import { UnreadyBanner } from "@/components/site/Unready";
import { DatagraphKindSchema } from "@/lib/datagraph/schema";
import * as z from "@/lib/validation";
import { SearchScreen } from "@/screens/search/SearchScreen";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  q: z.optional(z.string()),
  page: z.optional(
    z.pipe(
      z.string(),
      z.transform((v) => parseInt(v, 10)),
    ),
  ),
  kind: z.optional(
    z.pipe(
      z.transform((arg: unknown) => {
        if (typeof arg === "string") {
          return [arg];
        }

        return arg;
      }),
      z.array(DatagraphKindSchema),
    ),
  ),
  authors: z.optional(
    z.pipe(
      z.transform((arg: unknown) => {
        if (typeof arg === "string") {
          return [arg];
        }

        return arg;
      }),
      z.array(z.string()),
    ),
  ),
  categories: z.optional(
    z.pipe(
      z.transform((arg: unknown) => {
        if (typeof arg === "string") {
          return [arg];
        }

        return arg;
      }),
      z.array(z.string()),
    ),
  ),
  tags: z.optional(
    z.pipe(
      z.transform((arg: unknown) => {
        if (typeof arg === "string") {
          return [arg];
        }

        return arg;
      }),
      z.array(z.string()),
    ),
  ),
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const params = QuerySchema.parse(searchParams);

    let categoryIDs: string[] = [];
    if (params.categories?.length) {
      const { data } = await categoryList();
      const bySlug = keyBy("slug", data.categories);

      for (const slug of params.categories) {
        const category = bySlug[slug];
        if (category) {
          categoryIDs.push(category.id);
        }
      }
    }

    const { data } = params.q
      ? await datagraphSearch({
          q: params.q,
          page: params.page?.toString(),
          kind: params.kind,
          authors: params.authors,
          categories: categoryIDs,
          tags: params.tags,
        })
      : {
          data: undefined,
        };

    return (
      <SearchScreen
        initialQuery={params.q ?? ""}
        initialPage={params.page ?? 1}
        initialKind={params.kind ?? []}
        initialAuthors={params.authors ?? []}
        initialCategories={params.categories ?? []}
        initialTags={params.tags ?? []}
        initialResults={data}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
