import { PropsWithChildren } from "react";
import { RenderLeafProps } from "slate-react";

export function Leaf({
  attributes,
  children,
  leaf,
}: PropsWithChildren<RenderLeafProps>) {
  if (leaf.bold) {
    children = <strong>{children}</strong>;
  }

  // if (leaf.code) {
  //   children = <code>{children}</code>;
  // }

  if (leaf.italic) {
    children = <em>{children}</em>;
  }

  // if (leaf.underline) {
  //   children = <u>{children}</u>;
  // }

  return <span {...attributes}>{children}</span>;
}
