import {
  type ContentOutlineEntry,
  collectContentHeadings,
} from "@/components/ui/content-outline";

const blockLabels: Record<string, string> = {
  content: "Content",
  properties: "Properties",
  assets: "Assets",
  link: "Link",
  directory: "Directory",
};

export function collectLibraryOutline(
  root: HTMLElement,
): ContentOutlineEntry[] {
  const blocks = root.querySelectorAll<HTMLElement>("[data-outline-block]");
  const entries: ContentOutlineEntry[] = [];

  for (const target of blocks) {
    const type = target.dataset["outlineBlock"];

    if (!type) {
      continue;
    }

    const label = blockLabels[type];

    if (!label) {
      continue;
    }

    const visible = target.getBoundingClientRect().height > 0;

    if (!visible) {
      continue;
    }

    const hasText = Boolean(target.textContent?.trim());
    const hasMedia = target.querySelector("img, video, audio") !== null;
    const hasContent = hasText || hasMedia;

    if (!hasContent) {
      continue;
    }

    const entry: ContentOutlineEntry = {
      key: `block:${type}`,
      kind: "section",
      label,
      depth: 0,
      target,
    };

    entries.push(entry);

    if (type === "content") {
      entries.push(...collectContentHeadings(target, entry.key));
    }
  }

  return entries;
}
