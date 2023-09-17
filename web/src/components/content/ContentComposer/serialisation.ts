import markdown from "remark-parse";
import { serialize } from "remark-slate";
import slate from "remark-slate";
import { Descendant } from "slate";
import { unified } from "unified";

// Slate -> Markdown
export function serialise(value: Descendant[]): string {
  return value.map((v) => serialize(v)).join("\n");
}

// Markdown -> Slate
export function deserialise(md: string): Descendant[] {
  const vf = unified().use(markdown).use(slate).processSync(md);

  return vf.result as Descendant[];
}
