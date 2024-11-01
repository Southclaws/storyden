import { collectionList } from "@/api/openapi-server/collections";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { CollectionIndexScreen } from "@/screens/collection/CollectionIndexScreen";

export default async function Page() {
  try {
    const session = await getServerSession();
    const { data } = await collectionList();

    return (
      <CollectionIndexScreen session={session} initialCollections={data} />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
