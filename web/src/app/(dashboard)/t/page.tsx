import { threadList } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import * as z from "@/lib/validation";
import { ThreadIndexScreen } from "@/screens/thread/ThreadIndexScreen/ThreadIndexScreen";

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
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  const searchParams = await props.searchParams;
  const [initialSession, { data }] = await Promise.all([
    getServerSession(),
    threadList({
      ...(searchParams.q ? { q: searchParams.q } : {}),
      ...(searchParams.page ? { page: searchParams.page?.toString() } : {}),
    }),
  ]);

  return (
    <ThreadIndexScreen
      initialSession={initialSession}
      threads={data}
      page={searchParams.page}
      query={searchParams.q}
    />
  );
}
