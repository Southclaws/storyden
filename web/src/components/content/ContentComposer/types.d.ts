import type { BaseEditor } from "slate";
import { ReactEditor } from "slate-react";

//
// Why the custom types? Details here:
//
// https://docs.slatejs.org/concepts/12-typescript
//

type ParagraphElement = {
  type: "paragraph";
  children: CustomText[];
};

type LinkElement = {
  type: "link";
  link: string;
  children: CustomText[];
};

type HeadingOneElement = {
  type: "heading_one";
  children: CustomText[];
};
type HeadingTwoElement = {
  type: "heading_two";
  children: CustomText[];
};
type HeadingThreeElement = {
  type: "heading_three";
  children: CustomText[];
};
type HeadingFourElement = {
  type: "heading_four";
  children: CustomText[];
};
type HeadingFiveElement = {
  type: "heading_five";
  children: CustomText[];
};
type HeadingSixElement = {
  type: "heading_six";
  children: CustomText[];
};

type OrderedListElement = {
  type: "ol_list";
  children: CustomText[];
};

type UnorderedListElement = {
  type: "ul_list";
  children: CustomText[];
};

type ListItemElement = {
  type: "list_item";
  children: CustomText[];
};

type ImageElement = {
  type: "image";
  caption: string;
  link: string;
  children: EmptyText[];
};

type BlockQuoteElement = {
  type: "block_quote";
  children: CustomText[];
};

type CodeBlockElement = {
  type: "code_block";
  children: CustomText[];
};

type ThematicBreakElement = {
  type: "thematic_break";
  children: CustomText[];
};

export type CustomText = {
  text: string;
  bold?: boolean;
  italic?: boolean;
  underline?: boolean;
  code?: boolean;
};

export type CustomElement =
  | ParagraphElement
  | LinkElement
  | HeadingOneElement
  | HeadingTwoElement
  | HeadingThreeElement
  | HeadingFourElement
  | HeadingFiveElement
  | HeadingSixElement
  | OrderedListElement
  | UnorderedListElement
  | ListItemElement
  | ImageElement
  | BlockQuoteElement
  | CodeBlockElement
  | ThematicBreakElement;

type Formats = "bold" | "italic" | "underline";

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}
