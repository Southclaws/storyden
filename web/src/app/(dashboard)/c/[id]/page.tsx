import { UnreadyBanner } from "src/components/site/Unready";
import { CollectionScreen } from "src/screens/collection/CollectionScreen";

import { collectionGet } from "@/api/openapi-server/collections";
import { getServerSession } from "@/auth/server-session";

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
