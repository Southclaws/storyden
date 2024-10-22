import { z } from "zod";

import { MemberIndexScreen } from "src/screens/library/members/MemberIndexScreen/MemberIndexScreen";

type Props = {
  searchParams: Promise<{
    q: string;
    page: string;
  }>;
};

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});

export default async function Page(props: Props) {
  const params = QuerySchema.parse(await props.searchParams);
  return <MemberIndexScreen query={params.q} page={params.page} />;
}
