import { UnreadyBanner } from "src/components/site/Unready";

import { categoryList } from "@/api/openapi-server/categories";
import { threadList } from "@/api/openapi-server/threads";
import { CategoryScreen } from "@/screens/category/CategoryScreen";

type Props = {
  params: Promise<{
    category: string;
  }>;
};

export default async function Page(props: Props) {
  const slug = (await props.params).category;

  try {
    const { data: categoryListData } = await categoryList();

    const { data: threadListData } = await threadList({
      categories: [slug],
    });

    return (
      <CategoryScreen
        slug={slug}
        initialCategoryList={categoryListData}
        initialThreadList={threadListData}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
