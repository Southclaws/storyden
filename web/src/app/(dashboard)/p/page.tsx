import { z } from "zod";

import { MemberIndexScreen } from "src/screens/directory/members/MemberIndexScreen/MemberIndexScreen";

type Props = {
  searchParams: {
    q: string;
    page: string;
  };
};

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});

export default function Page(props: Props) {
  const params = QuerySchema.parse(props.searchParams);
  return <MemberIndexScreen query={params.q} page={params.page} />;
}
