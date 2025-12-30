import { z } from "zod";

import { MemberIndexScreen } from "src/screens/library/members/MemberIndexScreen/MemberIndexScreen";

import { profileList } from "@/api/openapi-server/profiles";
import { UnreadyBanner } from "@/components/site/Unready";

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
  roles: z
    .union([z.string(), z.array(z.string())])
    .transform((v) => (Array.isArray(v) ? v : v.split(",")))
    .optional(),
  invited_by: z
    .union([z.string(), z.array(z.string())])
    .transform((v) => (Array.isArray(v) ? v : v.split(",")))
    .optional(),
  joined: z.string().optional(),
  sort: z.string().optional(),
});
type Query = z.infer<typeof QuerySchema>;

type Props = {
  searchParams: Promise<Query>;
};

export default async function Page(props: Props) {
  try {
    const params = QuerySchema.parse(await props.searchParams);

    const { data } = await profileList({
      q: params.q,
      page: params.page?.toString(),
      roles: params.roles,
      invited_by: params.invited_by,
      joined: params.joined,
      sort: params.sort,
    });

    return (
      <MemberIndexScreen
        initialResult={data}
        query={params.q}
        page={params.page}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
