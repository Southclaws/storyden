import { z } from "zod";

import { SearchScreen } from "src/screens/search/SearchScreen";

import { datagraphSearch } from "@/api/openapi-server/datagraph";
import { UnreadyBanner } from "@/components/site/Unready";
import { DatagraphKindSchema } from "@/lib/datagraph/schema";

type Props = {
  searchParams: Promise<Query>;
};

export const dynamic = "force-dynamic";

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
  kind: z
    .preprocess((arg: unknown) => {
      if (typeof arg === "string") {
        return [arg];
      }

      return arg;
    }, z.array(DatagraphKindSchema))
    .optional(),
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const params = QuerySchema.parse(searchParams);

    const { data } = params.q
      ? await datagraphSearch({
          q: params.q,
          page: params.page?.toString(),
          kind: params.kind,
        })
      : {
          data: undefined,
        };

    return (
      <SearchScreen
        initialQuery={params.q ?? ""}
        initialPage={params.page ?? 1}
        initialKind={params.kind ?? []}
        initialResults={data}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
