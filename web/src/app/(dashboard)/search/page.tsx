import { z } from "zod";

import { EmptyState } from "src/components/site/EmptyState";
import { SearchScreen } from "src/screens/search/SearchScreen";

import { datagraphSearch } from "@/api/openapi-server/datagraph";
import { UnreadyBanner } from "@/components/site/Unready";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const params = QuerySchema.parse(searchParams);

    if (!params.q) {
      return (
        <EmptyState>
          <p>Search anything.</p>
        </EmptyState>
      );
    }

    const { data } = await datagraphSearch({ q: params.q });

    return (
      <SearchScreen
        query={params.q}
        page={params.page ?? 1}
        initialResults={data}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
