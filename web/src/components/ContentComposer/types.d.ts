import type { BaseEditor } from "slate";
import { ReactEditor } from "slate-react";

//
// Why the custom types? Details here:
//
// https://docs.slatejs.org/concepts/12-typescript
//

type CustomElement = {
  type: "paragraph";
  children: CustomText[];
};

type CustomText = {
  text: string;
  bold?: boolean;
  italic?: boolean;
};

type Formats = "bold" | "italic";

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}
