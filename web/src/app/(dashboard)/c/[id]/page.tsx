import { UnreadyBanner } from "@/components/site/Unready";
import { CollectionScreen } from "@/screens/collection/CollectionScreen";

import { collectionGet } from "@/api/openapi-server/collections";
import { getServerSession } from "@/auth/server-session";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

type Props = {
  params: Promise<{
    id: string;
  }>;
};

export default async function Page(props: Props) {
  try {
    const params = await props.params;
    const session = await getServerSession();
    const { data } = await collectionGet(params.id);

    return <CollectionScreen session={session} initialCollection={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
