import { CardGrid, CardRows } from "@/components/ui/rich-card";
import { CategoryTree } from "@/lib/category/tree";

import { CategoryCard } from "./CategoryCard";

type Props = {
  categories: CategoryTree[];
};

type LayoutProps = {
  layout: "grid" | "list";
} & Props;

export function CategoryLayout(props: LayoutProps) {
  switch (props.layout) {
    case "grid":
      return <CategoryCardGrid categories={props.categories} />;

    case "list":
      return <CategoryCardList categories={props.categories} />;
    default:
      return null;
  }
}

export function CategoryCardGrid({ categories }: Props) {
  return (
    <CardGrid>
      {categories.map((c) => (
        <CategoryCard key={c.id} category={c} showChildren={true} />
      ))}
    </CardGrid>
  );
}

export function CategoryCardList({ categories }: Props) {
  return (
    <CardRows>
      {categories.map((c) => (
        <CategoryCard key={c.id} category={c} showChildren={true} />
      ))}
    </CardRows>
  );
}
