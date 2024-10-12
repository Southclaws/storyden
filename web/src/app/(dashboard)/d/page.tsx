import { categoryList } from "@/api/openapi-server/categories";
import { UnreadyBanner } from "@/components/site/Unready";
import { CategoryIndexScreen } from "@/screens/category/CategoryIndexScreen";

export default async function Page() {
  try {
    const { data } = await categoryList();

    return <CategoryIndexScreen initialCategoryList={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
