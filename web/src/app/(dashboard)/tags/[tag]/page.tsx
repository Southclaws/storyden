import { tagGet } from "@/api/openapi-server/tags";
import { UnreadyBanner } from "@/components/site/Unready";
import { TagScreen } from "@/screens/tags/TagScreen";

type Props = {
  params: Promise<{
    tag: string;
  }>;
};

export default async function Page(props: Props) {
  const params = await props.params;
  try {
    const { tag } = params;

    const { data } = await tagGet(tag);
    return <TagScreen initialTag={data} slug={tag} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
