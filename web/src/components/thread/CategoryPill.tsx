import { readableColor, tint } from "polished";

import { CategoryReference } from "src/api/openapi-schema";

import { styled } from "@/styled-system/jsx";

type Props = { category: CategoryReference };

export function CategoryPill({ category }: Props) {
  const categoryColour = category.colour ?? "#e1e1e1";
  const backColour = tint(0.9, categoryColour);
  const textColour = readableColor(backColour);

  return (
    <styled.span
      px="3"
      py="1"
      borderRadius="lg"
      bgColor={backColour as any}
      color={textColour as any}
    >
      {category.name}
    </styled.span>
  );
}
