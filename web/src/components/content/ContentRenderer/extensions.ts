import { Link } from "@tiptap/extension-link";
import StarterKit from "@tiptap/starter-kit";

import { css } from "@/styled-system/css";

import { ImageNode } from "./ImageNode";
import { LinkPreviewNode } from "./LinkPreviewNode";

export const ContentLink = Link.extend({
  inclusive: false,
  parseHTML() {
    return [
      {
        tag: "a[href]:not([data-display])",
        getAttrs: (dom) => {
          const href = dom.getAttribute("href");
          return href ? { href } : false;
        },
      },
    ];
  },
});

export const contentExtensions = [
  StarterKit.configure({ link: false }),
  LinkPreviewNode,
  ContentLink,
  ImageNode.configure({
    allowBase64: false,
    HTMLAttributes: { class: css({ borderRadius: "md" }) },
  }),
];
