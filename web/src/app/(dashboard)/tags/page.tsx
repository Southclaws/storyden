import { tagList } from "@/api/openapi-server/tags";
import { UnreadyBanner } from "@/components/site/Unready";
import { TagsIndexScreen } from "@/screens/tags/TagsIndexScreen";

export default async function Page() {
  try {
    const { data } = await tagList();
    return <TagsIndexScreen initialTagList={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
