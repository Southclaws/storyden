import { categoryList } from "@/api/openapi-server/categories";
import { CategoryList } from "@/components/category/CategoryList/CategoryList";

export async function CategoryListServer() {
  try {
    const { data: initialCategoryList } = await categoryList();
    return <CategoryList initialCategoryList={initialCategoryList} />;
  } catch (e) {
    return <CategoryList />;
  }
}
