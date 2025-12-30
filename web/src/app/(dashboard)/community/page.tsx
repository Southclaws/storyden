import { z } from "zod";

import { UnreadyBanner } from "src/components/site/Unready";

import { getServerSession } from "@/auth/server-session";
import { getSettings } from "@/lib/settings/settings-server";
import { FeedScreen } from "@/screens/feed/FeedScreen";

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
    const session = await getServerSession();
    const settings = await getSettings();

    const { page } = QuerySchema.parse(await searchParams);

    return (
      <FeedScreen
        page={page ?? 1}
        initialSession={session}
        initialSettings={settings}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
