import { Node } from "@tiptap/core";

const TAG = "div";

export const LinkPreviewNode = Node.create<{}>({
  name: "linkPreview",

  group: "block",

  atom: true,

  selectable: true,

  addAttributes() {
    return {
      href: {
        default: null,
        parseHTML: (element) => element.getAttribute("data-href"),
        renderHTML: (attributes: Record<string, any>) => {
          if (!attributes["href"]) {
            return {};
          }
          return { "data-href": attributes["href"] };
        },
      },
      "data-display": {
        default: "card",
        parseHTML: (element) => element.getAttribute("data-display"),
        renderHTML: (attributes: Record<string, any>) => {
          return { "data-display": attributes["data-display"] };
        },
      },
      "data-diff": {
        default: null,
        parseHTML: (element) => element.getAttribute("data-diff"),
        renderHTML: (attributes: Record<string, any>) => {
          if (!attributes["data-diff"]) {
            return {};
          }
          return { "data-diff": attributes["data-diff"] };
        },
      },
    };
  },

  parseHTML() {
    return [
      {
        tag: `${TAG}[data-display="card"]`,
        priority: 100,
      },
    ];
  },

  renderHTML({ node }) {
    const href = node.attrs["href"] || "";
    const dataDiff = node.attrs["data-diff"];

    const attrs: Record<string, string> = {
      "data-href": href,
      "data-display": "card",
      class: "link-card",
    };

    if (dataDiff) {
      attrs["data-diff"] = dataDiff;
    }

    return [TAG, attrs];
  },
});
