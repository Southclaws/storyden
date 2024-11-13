import { linkGet } from "@/api/openapi-server/links";
import { UnreadyBanner } from "@/components/site/Unready";
import { LinkScreen } from "@/screens/library/links/LinkScreen";

type Props = {
  params: Promise<{
    slug: string;
  }>;
};

export default async function Page(props: Props) {
  try {
    const params = await props.params;

    const { data } = await linkGet(params.slug);

    return <LinkScreen initialLink={data} slug={params.slug} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
