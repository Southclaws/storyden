import { linkList } from "@/api/openapi-server/links";
import { UnreadyBanner } from "@/components/site/Unready";
import { LinkIndexScreen } from "@/screens/library/links/LinkIndexScreen";

type Props = {
  searchParams: Promise<{
    q: string;
    page: number;
  }>;
};

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const { data } = await linkList({
      q: searchParams.q,
      page: searchParams.page?.toString(),
    });

    return (
      <LinkIndexScreen
        initialResult={data}
        query={searchParams.q}
        page={searchParams.page}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
