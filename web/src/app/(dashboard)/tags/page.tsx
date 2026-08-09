import { tagList } from "@/api/openapi-server/tags";
import { UnreadyBanner } from "@/components/site/Unready";
import { TagsIndexScreen } from "@/screens/tags/TagsIndexScreen";

type Props = {
  searchParams: Promise<{ page?: string }>;
};

export default async function Page(props: Props) {
  const { page } = await props.searchParams;

  try {
    const { data } = await tagList({ page });
    return <TagsIndexScreen initialTagList={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
