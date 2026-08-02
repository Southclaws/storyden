"use client";

import { blockEditor } from "styled-system/recipes";

import type { HTMLStyledProps } from "@/styled-system/types";
import { createStyleContext } from "@/utils/create-style-context";

const { withProvider, withContext } = createStyleContext(blockEditor);

export const Root = withProvider<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "root",
);
export const Gutter = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "gutter",
);
export const Handle = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "handle",
);

export const Content = withContext<HTMLDivElement, HTMLStyledProps<"div">>(
  "div",
  "content",
);
