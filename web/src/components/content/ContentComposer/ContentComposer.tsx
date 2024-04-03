import { EditorContent } from "@tiptap/react";
import { BoldIcon, ItalicIcon, StrikethroughIcon } from "lucide-react";

import { Button } from "src/theme/components/Button";

import { css } from "@/styled-system/css";
import { HStack, LStack } from "@/styled-system/jsx";

import { Props, useContentComposer } from "./useContentComposer";

export function ContentComposer(props: Props) {
  const { editor, handlers } = useContentComposer(props);

  return (
    <LStack
      id="rich-text-editor"
      className="typography"
      w="full"
      h="full"
      gap="1"
      onDragOver={(e) => e.preventDefault()}
    >
      <HStack>
        <Button
          type="button"
          size="xs"
          kind="ghost"
          onClick={handlers.handleBold}
        >
          <BoldIcon />
        </Button>
        <Button
          type="button"
          size="xs"
          kind="ghost"
          onClick={handlers.handleItalic}
        >
          <ItalicIcon />
        </Button>
        <Button
          type="button"
          size="xs"
          kind="ghost"
          onClick={handlers.handleStrike}
        >
          <StrikethroughIcon />
        </Button>
      </HStack>

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
