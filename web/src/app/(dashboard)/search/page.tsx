import { z } from "zod";

import { EmptyState } from "src/components/site/EmptyState";
import { SearchScreen } from "src/screens/search/SearchScreen";

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
  const params = QuerySchema.parse(await props.searchParams);

  if (!params.q) {
    return <EmptyState />;
  }

  return <SearchScreen query={params.q} />;
}
