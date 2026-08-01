import { styled } from "@/styled-system/jsx";
import { cardBox } from "@/styled-system/recipes";
import type { ComponentProps } from "@/styled-system/types";

export type CardBoxProps = ComponentProps<typeof CardBox>;
export const CardBox = styled("div", cardBox);
