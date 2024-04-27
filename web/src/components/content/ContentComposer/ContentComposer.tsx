import { EditorContent } from "@tiptap/react";
import {
  BoldIcon,
  ImageIcon,
  ItalicIcon,
  StrikethroughIcon,
} from "lucide-react";

import { Button } from "src/theme/components/Button";

import "./styles.css";

import { css } from "@/styled-system/css";
import { LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { FloatingMenu } from "./plugins/MenuPlugin";
import { Props, useContentComposer } from "./useContentComposer";

export function ContentComposer(props: Props) {
  const { editor, handlers } = useContentComposer(props);

  return (
    <LStack
      id="rich-text-editor"
      containerType="inline-size"
      className="typography"
      w="full"
      h="full"
      gap="1"
      onDragOver={(e) => e.preventDefault()}
    >
      {editor && (
        <FloatingMenu editor={editor}>
          <Button
            type="button"
            size="xs"
            kind="ghost"
            title="Toggle bold text"
            onClick={handlers.handleBold}
          >
            <BoldIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            kind="ghost"
            title="Toggle italic text"
            onClick={handlers.handleItalic}
          >
            <ItalicIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            kind="ghost"
            title="Toggle strikeout text"
            onClick={handlers.handleStrike}
          >
            <StrikethroughIcon />
          </Button>
          &nbsp;
          <label
            className={button({
              size: "xs",
              kind: "ghost",
            })}
            htmlFor="filepicker"
            title="Insert an image"
          >
            <ImageIcon />
          </label>
          <styled.input
            id="filepicker"
            type="file"
            multiple
            display="none"
            onChange={handlers.handleFileUpload}
          />
        </FloatingMenu>
      )}

      <EditorContent
        id="editor-content"
        className={css({
          // NOTE: We want to make the clickable area expand to the full height.
          height: "full",
          width: "full",
        })}
        editor={editor}
      />
    </LStack>
  );
}
