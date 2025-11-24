import { Plugin, PluginKey } from "@tiptap/pm/state";
import { Extension } from "@tiptap/react";
import { find } from "linkifyjs";

export type LinkPasteMenuState = {
  isVisible: boolean;
  url: string | null;
  position: number;
  range: { from: number; to: number } | null;
};

export const linkPasteMenuKey = new PluginKey<LinkPasteMenuState>(
  "linkPasteMenu",
);

export const LinkPasteMenuPlugin = Extension.create({
  name: "linkPasteMenu",

  addProseMirrorPlugins() {
    return [
      new Plugin<LinkPasteMenuState>({
        key: linkPasteMenuKey,

        state: {
          init() {
            return {
              isVisible: false,
              url: null,
              position: 0,
              range: null,
            };
          },

          apply(tr, value) {
            const meta = tr.getMeta(linkPasteMenuKey);

            if (meta) {
              return meta;
            }

            if (!tr.docChanged) {
              return value;
            }

            if (value.isVisible) {
              return {
                isVisible: false,
                url: null,
                position: 0,
                range: null,
              };
            }

            return value;
          },
        },

        props: {
          handlePaste(view, event) {
            const { state } = view;
            const { selection } = state;
            const { $from } = selection;

            const parentNode = $from.parent;
            const isEmptyParagraph =
              parentNode.type.name === "paragraph" &&
              parentNode.content.size === 0;

            if (!isEmptyParagraph) {
              return false;
            }

            const clipboardData = event.clipboardData;
            if (!clipboardData) {
              return false;
            }

            const pastedText = clipboardData.getData("text/plain").trim();

            if (!pastedText) {
              return false;
            }

            const links = find(pastedText);
            if (links.length === 0 || links[0] === undefined) {
              return false;
            }

            const link = links[0];

            const isSingleURL = links.length === 1 && link.value === pastedText;

            if (!isSingleURL) {
              return false;
            }

            const url = link.href;

            event.preventDefault();

            const insertPos = selection.from;
            const tr = state.tr.insertText(pastedText, insertPos);

            const newPos = insertPos + pastedText.length;

            tr.setMeta(linkPasteMenuKey, {
              isVisible: true,
              url,
              position: newPos,
              range: { from: insertPos, to: newPos },
            });

            view.dispatch(tr);

            return true;
          },
        },
      }),
    ];
  },
});
