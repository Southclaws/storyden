import { z } from "zod";

import { MemberIndexScreen } from "src/screens/library/members/MemberIndexScreen/MemberIndexScreen";

import { profileList } from "@/api/openapi-server/profiles";
import { UnreadyBanner } from "@/components/site/Unready";

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
  try {
    const params = QuerySchema.parse(await props.searchParams);

    const { data } = await profileList({
      q: params.q,
      page: params.page?.toString(),
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
