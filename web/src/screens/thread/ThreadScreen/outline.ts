import {
  type ContentOutlineEntry,
  collectContentHeadings,
} from "@/components/ui/content-outline";

export function collectThreadOutline(root: HTMLElement): ContentOutlineEntry[] {
  const entries: ContentOutlineEntry[] = [];
  const original = root.querySelector<HTMLElement>("[data-outline-original]");

  if (original) {
    entries.push({
      key: "original",
      kind: "section",
      label: "Original post",
      depth: 0,
      target: original,
    });
  }

  const body = root.querySelector<HTMLElement>("[data-outline-body]");

  if (body) {
    const bodyHeight = body.getBoundingClientRect().height;
    const longContent = bodyHeight > window.innerHeight;

    if (longContent) {
      entries.push(...collectContentHeadings(body, "original"));
    }
  }

  const replies = root.querySelectorAll<HTMLElement>("[data-outline-reply]");

  for (const target of replies) {
    entries.push({
      key: target.id,
      kind: "reply",
      label: target.dataset["outlineLabel"] || "Reply",
      description: target.dataset["outlineDescription"],
      createdAt: target.dataset["outlineCreatedAt"],
      depth: 0,
      target,
    });
  }

  return entries;
}
