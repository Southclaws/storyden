import { redirect } from "next/navigation";

import { getServerSession } from "src/auth/server-session";

import { nodeGet } from "@/api/openapi-server/nodes";
import { getTargetSlug } from "@/components/library/utils";
import { UnreadyBanner } from "@/components/site/Unready";
import { LibraryPageContainerScreen } from "@/screens/library/LibraryPageContainerScreen/LibraryPageContainerScreen";
import { LibraryPageCreateScreen } from "@/screens/library/LibraryPageCreateScreen/LibraryPageCreateScreen";
import { Params, ParamsSchema, Query } from "@/screens/library/library-path";

type Props = {
  params: Params;
  searchParams: Query;
};

export default async function Page(props: Props) {
  try {
    const { slug } = ParamsSchema.parse(props.params);
    const session = await getServerSession();

    const [targetSlug, isNew] = getTargetSlug(slug);

    if (targetSlug) {
      if (isNew) {
        if (!session) {
          redirect(`/login`); // TODO: ?return= back to this path.
        }

        return <LibraryPageCreateScreen session={session} />;
      }

      const { data } = await nodeGet(targetSlug);

      if (data) {
        return <LibraryPageContainerScreen slug={targetSlug} node={data} />;
      }
    }

    // Creating a new item or node from the root: "/l/new"
    if (isNew) {
      if (!session) {
        redirect(`/login`); // TODO: ?return= back to this path.
      }

      return <LibraryPageCreateScreen session={session} />;
    }

    // NOTE: This state is probably not possible to reach due to the params.
    throw new Error("Library page not found");
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
