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

export type ImageElement = {
  type: "image";
  caption: string;
  link: string;
  children: EmptyText[];
};

type CustomText = {
  text: string;
  bold?: boolean;
  italic?: boolean;
};

type CustomElement = ParagraphElement | ImageElement;

type Formats = "bold" | "italic";

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}
