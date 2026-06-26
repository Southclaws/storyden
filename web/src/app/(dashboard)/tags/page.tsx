import { tagList } from "@/api/openapi-server/tags";
import { UnreadyBanner } from "@/components/site/Unready";
import { TagsIndexScreen } from "@/screens/tags/TagsIndexScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  try {
    const { data } = await tagList();
    return <TagsIndexScreen initialTagList={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
