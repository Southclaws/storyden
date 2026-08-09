import { useCategoryList } from "@/api/openapi-client/categories";
import { CategoryLayout } from "@/components/category/CategoryIndex/CategoryCardLayout";
import { Unready } from "@/components/site/Unready";
import { buildCategoryTree } from "@/lib/category/tree";
import { FeedBlock } from "@/lib/settings/feed";

import { useFeedBlockEditor } from "../Context";

export function FeedCategoriesBlock({
  block,
}: {
  block: Extract<FeedBlock, { type: "categories" }>;
}) {
  const { initialData } = useFeedBlockEditor();
  const { data, error } = useCategoryList({
    swr: { fallbackData: initialData.initialCategoryList },
  });

  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <CategoryLayout
      layout={block.layout}
      categories={buildCategoryTree(data.categories)}
    />
  );
}
