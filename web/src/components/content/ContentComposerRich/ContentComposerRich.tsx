import { BubbleMenu, EditorContent } from "@tiptap/react";

import { EditIcon } from "@/components/ui/icons/Edit";
import { css, cx } from "@/styled-system/css";
import { LStack } from "@/styled-system/jsx";

import { ComposerTools } from "../ComposerTools";
import { ContentComposerProps } from "../composer-props";

import "./styles.css";

import { EditorMenu } from "./EditorMenu";
import { useContentComposer } from "./useContentComposerRich";

export function ContentComposerRich(props: ContentComposerProps) {
  const { editor, initialValueHTML, uniqueID, handlers, format } =
    useContentComposer(props);

  return (
    <LStack
      id={`rich-text-editor-${uniqueID}`}
      containerType="inline-size"
      className={cx("typography", props.className)}
      position="relative"
      w="full"
      gap="1"
      minHeight="8"
      onDragOver={(e) => e.preventDefault()}
    >
      {editor ? (
        <>
          <ComposerTools icon={<EditIcon />}>
            <EditorMenu
              editor={editor}
              uniqueID={uniqueID}
              format={format}
              handlers={handlers}
            />
          </ComposerTools>
          <BubbleMenu
            editor={editor}
            tippyOptions={{
              placement: "bottom-start",
              maxWidth: "100%",
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
                      boundary: editor?.view.dom,
                      padding: 8,
                    },
                  },
                  {
                    name: "preventOverflow",
                    options: {
                      boundary: editor?.view.dom,
                      altAxis: true,
                      padding: {
                        top: 0,
                        right: 0,
                        // Some negative padding on the bottom allows the menu
                        // to overflow the bottom of the editor area for cases
                        // where the editor is only a single line. Without this,
                        // the menu can only be placed over the text itself.
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
              display: "flex",
              flexWrap: "wrap",
              minW: "0",
              maxW: "full",
              gap: "1",
              padding: "1",
              backgroundColor: "bg.subtle",
              backdropBlur: "frosted",
              backdropFilter: "auto",
              boxShadow: "md",
            })}
          >
            <EditorMenu
              editor={editor}
              uniqueID={uniqueID}
              format={format}
              handlers={handlers}
            />
          </BubbleMenu>
          <EditorContent
            id={`editor-content-${uniqueID}`}
            className={css({
              height: "full",
              width: "full",
            })}
            editor={editor}
          />
        </>
      ) : (
        <div dangerouslySetInnerHTML={{ __html: initialValueHTML }} />
      )}
    </LStack>
  );
}
