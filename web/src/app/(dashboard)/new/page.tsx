import { redirect } from "next/navigation";
import { z } from "zod";

import { ComposeScreen } from "src/screens/compose/ComposeScreen";

import { getServerSession } from "@/auth/server-session";

const QuerySchema = z.object({
  id: z.string().optional(),
});

type Props = {
  searchParams: Promise<z.infer<typeof QuerySchema>>;
};

export default async function Page(props: Props) {
  const session = await getServerSession();
  if (!session) {
    redirect("/login");
  }

  const searchParams = await props.searchParams;
  const params = QuerySchema.parse(searchParams);
  return <ComposeScreen editing={params.id} />;
}
