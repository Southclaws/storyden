import { nodeGet } from "@/api/openapi-server/nodes";
import { getTargetSlug } from "@/components/library/utils";
import { UnreadyBanner } from "@/components/site/Unready";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";
import { Params, ParamsSchema } from "@/screens/library/library-path";

type Props = {
  params: Params;
};

export default async function Page(props: Props) {
  try {
    const { slug } = ParamsSchema.parse(props.params);

    const targetSlug = getTargetSlug(slug);

    if (!targetSlug) {
      // NOTE: This state is probably not possible to reach due to the params.
      throw new Error("Library page not found");
    }

    const { data } = await nodeGet(targetSlug);

    return <LibraryPageScreen node={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
