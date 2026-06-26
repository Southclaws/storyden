import { linkGet } from "@/api/openapi-server/links";
import { UnreadyBanner } from "@/components/site/Unready";
import { LinkScreen } from "@/screens/library/links/LinkScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

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
