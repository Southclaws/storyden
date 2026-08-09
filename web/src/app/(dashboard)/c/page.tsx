import { collectionList } from "@/api/openapi-server/collections";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { CollectionIndexScreen } from "@/screens/collection/CollectionIndexScreen";

type Props = {
  searchParams: Promise<{ page?: string }>;
};

export default async function Page(props: Props) {
  const { page } = await props.searchParams;

  try {
    const session = await getServerSession();
    const { data } = await collectionList({ page });

    return (
      <CollectionIndexScreen session={session} initialCollections={data} />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
