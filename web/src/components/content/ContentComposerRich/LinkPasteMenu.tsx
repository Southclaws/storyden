import { BubbleMenu, Editor } from "@tiptap/react";

import { Button } from "@/components/ui/button";
import { CardIcon } from "@/components/ui/icons/Card";
import { LinkIcon } from "@/components/ui/icons/Link";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";

import { linkPasteMenuKey } from "./plugins/LinkPasteMenuPlugin";

type Props = {
  editor: Editor;
};

export function LinkPasteMenu({ editor }: Props) {
  const dismissMenu = () => {
    const tr = editor.state.tr.setMeta(linkPasteMenuKey, {
      isVisible: false,
      url: null,
      position: 0,
      range: null,
    });
    editor.view.dispatch(tr);
  };

  const handleLinkChoice = () => {
    const menuState = linkPasteMenuKey.getState(editor.state);
    if (!menuState?.range || !menuState?.url) return;

    editor
      .chain()
      .focus()
      .setTextSelection(menuState.range)
      .setLink({ href: menuState.url })
      .run();

    dismissMenu();
  };

  const handleCardChoice = () => {
    const menuState = linkPasteMenuKey.getState(editor.state);
    if (!menuState?.range || !menuState?.url) return;

    editor
      .chain()
      .focus()
      .deleteRange(menuState.range)
      .setLinkPreview({ href: menuState.url })
      .run();

    dismissMenu();
  };

  return (
    <BubbleMenu
      editor={editor}
      shouldShow={() => {
        const menuState = linkPasteMenuKey.getState(editor.state);
        return menuState?.isVisible ?? false;
      }}
      tippyOptions={{
        placement: "bottom-start",
        maxWidth: "100%",
        onHide: () => {
          dismissMenu();
        },
        popperOptions: {
          modifiers: [
            {
              name: "offset",
              options: {
                offset: [0, 4],
              },
            },
            {
              name: "flip",
              options: {
                fallbackPlacements: ["top-start"],
                boundary: editor.view.dom,
                padding: 8,
              },
            },
            {
              name: "preventOverflow",
              options: {
                boundary: editor.view.dom,
                altAxis: true,
                padding: {
                  top: 0,
                  right: 0,
                  bottom: -40,
                  left: 0,
                },
                rootBoundary: "viewport",
                tether: false,
              },
            },
          ],
        },
      }}
      className={css({
        zIndex: "popover",
        borderRadius: "md",
        backgroundColor: "background.overlay/80",
        backdropBlur: "frosted",
        backdropFilter: "auto",
        boxShadow: "floating",
        padding: "1",
      })}
    >
      <Text variant="supporting">Show link as</Text>
      <HStack gap="1">
        <Button type="button" variant="subtle" onClick={handleLinkChoice}>
          <LinkIcon /> Text
        </Button>
        <Button type="button" variant="subtle" onClick={handleCardChoice}>
          <CardIcon /> Preview
        </Button>
      </HStack>
    </BubbleMenu>
  );
}
