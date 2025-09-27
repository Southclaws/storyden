import Link from "next/link";

import { CategoryReference } from "@/api/openapi-schema";
import { categoryColourCSS } from "@/lib/category/colours";

import { Badge, BadgeProps } from "../ui/badge";

type Props = {
  category: CategoryReference;
  asLink?: boolean;
};

export function CategoryBadge({
  category,
  asLink = true,
  ...props
}: Props & BadgeProps) {
  const cssProps = categoryColourCSS(category.colour);

  const path = `/d/${category.slug}`;

  const children = (
    <Badge
      size="sm"
      style={cssProps}
      bgColor="colorPalette.bg"
      borderColor="colorPalette.border"
      color="colorPalette.fg"
      // as any: expression produces a union that is too complex... (???)
      {...(props as any)}
    >
      {category.name}
    </Badge>
  );

  if (asLink) {
    return <Link href={path}>{children}</Link>;
  }

  return children;
}
