import { Tag } from "@chakra-ui/react";
import { readableColor, tint } from "polished";
import { CategoryReference } from "src/api/openapi/schemas";

type Props = { category: CategoryReference };
export function CategoryPill({ category }: Props) {
  const categoryColour = category.colour ?? "#e1e1e1";
  const backColour = tint(0.9, categoryColour);
  const textColour = readableColor(backColour);

  return (
    <Tag
      size="sm"
      variant="subtle"
      px={3}
      py={1}
      borderRadius="lg"
      bgColor={backColour}
      textColor={textColour}
    >
      {category.name}
    </Tag>
  );
}
