import z from "zod";

import { Permission } from "@/api/openapi-schema";
import { robotSessionGet } from "@/api/openapi-server/robots";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { RobotSessionScreen } from "@/screens/robots/RobotSessionScreen";
import { hasPermission } from "@/utils/permissions";

type Props = {
  params: Promise<{
    id: string;
  }>;
  searchParams: Promise<Query>;
};

export const QuerySchema = z.object({
  before: z.string().optional(),
  limit: z.string().optional(),
  robot: z.string().optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  try {
    const params = await props.params;
    const session = await getServerSession();
    if (!session) {
      return (
        <UnreadyBanner error={"You must be logged in to view this page."} />
      );
    }

    if (!hasPermission(session, Permission.USE_ROBOTS)) {
      return (
        <UnreadyBanner error={"You do not have permission to use robots."} />
      );
    }

    const searchParams = QuerySchema.parse(await props.searchParams);
    const before = searchParams.before;
    const limit = searchParams.limit;
    const robot = searchParams.robot;

    if (params.id === "new") {
      return (
        <RobotSessionScreen
          initialSession={session}
          initialChatSession={null}
          initialChatBefore={before}
          initialChatLimit={limit}
          initialSelectedRobotID={robot}
        />
      );
    }

    const { data } = await robotSessionGet(
      params.id,
      {
        before,
        limit,
      },
      { cache: "no-store" },
    );

    return (
      <RobotSessionScreen
        initialSession={session}
        initialChatSession={data}
        initialChatBefore={before}
        initialChatLimit={limit}
        initialSelectedRobotID={robot}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
