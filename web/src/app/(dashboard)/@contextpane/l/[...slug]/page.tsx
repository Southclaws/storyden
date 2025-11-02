import { nodeGet } from "@/api/openapi-server/nodes";
import { LibraryPageScreenContextPane } from "@/screens/library/LibraryPageScreen/LibraryPageScreenContextPane";

export default async function Page(props: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await props.params;

  console.log(slug);

  try {
    const { data } = await nodeGet(slug);

    return <LibraryPageScreenContextPane slug={slug} initialNode={data} />;
  } catch (e) {
    console.error(e);
    return null;
  }
}
