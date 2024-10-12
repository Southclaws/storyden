import Link from "next/link";

import { Category } from "@/api/openapi-schema";
import { categoryColourCSS } from "@/lib/category/colours";

import { Badge, BadgeProps } from "../ui/badge";

type Props = {
  category: Category;
};

export function CategoryBadge({ category, ...props }: Props & BadgeProps) {
  const cssProps = categoryColourCSS(category.colour);

  const path = `/d/${category.slug}`;

  return (
    <Link href={path}>
      <Badge
        size="sm"
        style={cssProps}
        bgColor="colorPalette"
        borderColor="colorPalette.muted"
        color="colorPalette.text"
        // as any: expression produces a union that is too complex... (???)
        {...(props as any)}
      >
        {category.name}
      </Badge>
    </Link>
  );
}
