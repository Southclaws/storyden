import { categoryGet } from "@/api/openapi-server/categories";
import { threadList } from "@/api/openapi-server/threads";
import { CategoryScreenContextPane } from "@/screens/category/CategoryScreenContextPane";

export default async function Page(props: {
  params: Promise<{ category: string }>;
}) {
  const { category } = await props.params;

  try {
    const { data: categoryData } = await categoryGet(category);

    const { data: threadListData } = await threadList({
      categories: [category],
    });

    return (
      <CategoryScreenContextPane
        slug={category}
        initialCategory={categoryData}
        initialThreadList={threadListData}
      />
    );
  } catch (e) {
    console.error(e);
    return null;
  }
}
