import z from "zod";

import { UnreadyBanner } from "src/components/site/Unready";

import { Permission } from "@/api/openapi-schema";
import { robotSessionsList } from "@/api/openapi-server/robots";
import { getServerSession } from "@/auth/server-session";
import { RobotSessionListScreen } from "@/screens/robots/RobotSessionListScreen";
import { hasPermission } from "@/utils/permissions";

type Props = {
  searchParams: Promise<Query>;
};

export const QuerySchema = z.object({
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});
export type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  try {
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
    const page = searchParams.page?.toString();

    const { data } = await robotSessionsList({ page: page });

    return (
      <RobotSessionListScreen
        initialSession={session}
        initialChatSessionList={data}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
