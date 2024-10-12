import { categoryList } from "@/api/openapi-server/categories";
import { threadList } from "@/api/openapi-server/threads";
import { CategoryScreenContextPane } from "@/screens/category/CategoryScreenContextPane";

export default async function Page(props: { params: { category: string } }) {
  const { category } = props.params;

  try {
    const { data: categoryListData } = await categoryList();

    const { data: threadListData } = await threadList({
      categories: [category],
    });

    return (
      <CategoryScreenContextPane
        slug={category}
        initialCategoryList={categoryListData}
        initialThreadList={threadListData}
      />
    );
  } catch (e) {
    console.error(e);
    return null;
  }
}
