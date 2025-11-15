import DOMPurify from "isomorphic-dompurify";
import html from "remark-html";
import remarkParse from "remark-parse";
import TurndownService from "turndown";
import { unified } from "unified";

import { deriveError } from "./error";

const SANITIZE_CONFIG = {
  ALLOWED_TAGS: [
    "p",
    "br",
    "strong",
    "em",
    "u",
    "s",
    "del",
    "code",
    "pre",
    "blockquote",
    "h1",
    "h2",
    "h3",
    "h4",
    "h5",
    "h6",
    "ul",
    "ol",
    "li",
    "a",
    "img",
    "table",
    "thead",
    "tbody",
    "tr",
    "td",
    "th",
    "hr",
  ],
  ALLOWED_ATTR: ["href", "src", "alt", "title", "class", "target", "rel"],
  ALLOW_DATA_ATTR: false,
};

const HTML_TO_MD_SANITIZE_CONFIG = {
  ...SANITIZE_CONFIG,
  ALLOWED_TAGS: [...SANITIZE_CONFIG.ALLOWED_TAGS, "b", "i", "strike"],
};

const turndownService = new TurndownService({
  headingStyle: "atx",
  codeBlockStyle: "fenced",
  bulletListMarker: "-",
});

turndownService.addRule("strikethrough", {
  filter: ["del", "s"],
  replacement: (content) => `~~${content}~~`,
});

export async function markdownToHTML(markdown: string): Promise<string> {
  try {
    const result = await unified()
      .use(remarkParse)
      .use(html, { sanitize: false })
      .process(markdown);

    const rawHTML = String(result);

    DOMPurify.addHook("afterSanitizeAttributes", (node) => {
      if (node.tagName === "A") {
        const href = node.getAttribute("href");
        if (href && /^https?:\/\//i.test(href)) {
          node.setAttribute("rel", "noopener noreferrer");
        }
      }
    });

    const sanitized = DOMPurify.sanitize(rawHTML, SANITIZE_CONFIG);

    DOMPurify.removeHooks("afterSanitizeAttributes");

    return sanitized;
  } catch (error) {
    const preview =
      markdown.length > 100 ? `${markdown.slice(0, 100)}...` : markdown;
    console.error("Failed to process markdown to HTML:", {
      error,
      markdownPreview: preview,
    });
    throw new Error(`Markdown processing failed: ${deriveError(error)}`);
  }
}

export function htmlToMarkdown(html: string): string {
  const sanitized = DOMPurify.sanitize(html, HTML_TO_MD_SANITIZE_CONFIG);

  return turndownService.turndown(sanitized);
}
