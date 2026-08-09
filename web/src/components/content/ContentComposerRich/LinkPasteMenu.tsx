import { Editor } from "@tiptap/react";
import { BubbleMenu } from "@tiptap/react/menus";

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
      options={{
        placement: "bottom-start",
        offset: { mainAxis: 4 },
        flip: {
          fallbackPlacements: ["top-start"],
          boundary: editor.view.dom,
          padding: 8,
        },
        shift: {
          boundary: editor.view.dom,
          crossAxis: true,
          padding: {
            top: 0,
            right: 0,
            bottom: -40,
            left: 0,
          },
          rootBoundary: "viewport",
        },
        onHide: () => {
          dismissMenu();
        },
      }}
      className={css({
        zIndex: "popover",
        borderRadius: "md",
        background: "background.overlay",
        backdropBlur: "subtle",
        backdropFilter: "auto",
        borderColor: "border.strong",
        borderWidth: "thin",
        boxShadow: "floating",
        color: "text.default",
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
