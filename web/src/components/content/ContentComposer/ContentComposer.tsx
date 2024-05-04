import { Portal } from "@ark-ui/react";
import { EditorContent } from "@tiptap/react";
import {
  BoldIcon,
  CodeIcon,
  CodeSquareIcon,
  Heading1Icon,
  Heading2Icon,
  Heading3Icon,
  Heading4Icon,
  Heading5Icon,
  Heading6Icon,
  ImageIcon,
  ItalicIcon,
  ListIcon,
  ListOrderedIcon,
  StrikethroughIcon,
  TextIcon,
  TextQuoteIcon,
} from "lucide-react";
import { match } from "ts-pattern";

import { Button } from "src/theme/components/Button";
import {
  Menu,
  MenuContent,
  MenuItem,
  MenuPositioner,
  MenuTrigger,
} from "src/theme/components/Menu";

import "./styles.css";

import { css } from "@/styled-system/css";
import { LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { FloatingMenu } from "./plugins/MenuPlugin";
import { Props, useContentComposer } from "./useContentComposer";

export function ContentComposer(props: Props) {
  const { editor, handlers, format } = useContentComposer(props);

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
          <Menu
            size="sm"
            userSelect="none"
            onSelect={(d) => format.text.set(d.value as any /* lazy */)}
          >
            <MenuTrigger asChild>
              <Button
                type="button"
                size="xs"
                kind="ghost"
                title="Change the kind of text"
              >
                {match(format.text.active)
                  .with("p", () => <TextIcon />)
                  .with("h1", () => <Heading1Icon />)
                  .with("h2", () => <Heading2Icon />)
                  .with("h3", () => <Heading3Icon />)
                  .with("h4", () => <Heading4Icon />)
                  .with("h5", () => <Heading5Icon />)
                  .with("h6", () => <Heading6Icon />)
                  .otherwise(() => (
                    <TextIcon />
                  ))}
              </Button>
            </MenuTrigger>

            <Portal>
              {/* NOTE: Because this is a portal, we need to reference this ID
              in the FloatingMenu unfocus logic so we don't hide the menu when
              this menu is opened as it's a portal and not a child element. */}
              <MenuPositioner>
                <MenuContent
                  id="text-block-menu"
                  backdropBlur="md"
                  backdropFilter="auto"
                >
                  <MenuItem id="p">
                    <TextIcon />
                    &nbsp;Paragraph
                  </MenuItem>

                  <MenuItem id="h1" fontSize="lg" fontWeight="extrabold">
                    <Heading1Icon />
                    &nbsp;Heading 1
                  </MenuItem>
                  <MenuItem id="h2" fontSize="md" fontWeight="extrabold">
                    <Heading2Icon />
                    &nbsp;Heading 2
                  </MenuItem>
                  <MenuItem id="h3" fontSize="md" fontWeight="bold">
                    <Heading3Icon />
                    &nbsp;Heading 3
                  </MenuItem>
                  <MenuItem id="h4" fontSize="md" fontWeight="medium">
                    <Heading4Icon />
                    &nbsp;Heading 4
                  </MenuItem>
                  <MenuItem id="h5" fontSize="sm" fontWeight="normal">
                    <Heading5Icon />
                    &nbsp;Heading 5
                  </MenuItem>
                  <MenuItem id="h6" fontSize="sm" fontWeight="light">
                    <Heading6Icon />
                    &nbsp;Heading 6
                  </MenuItem>
                </MenuContent>
              </MenuPositioner>
            </Portal>
          </Menu>
          <Button
            type="button"
            size="xs"
            kind={format.bold.isActive ? "primary" : "ghost"}
            title="Toggle bold text"
            onClick={format.bold.toggle}
          >
            <BoldIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            kind={format.italic.isActive ? "primary" : "ghost"}
            title="Toggle italic text"
            onClick={format.italic.toggle}
          >
            <ItalicIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            kind={format.strike.isActive ? "primary" : "ghost"}
            title="Toggle strikeout text"
            onClick={format.strike.toggle}
          >
            <StrikethroughIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            kind={format.code.isActive ? "primary" : "ghost"}
            title="Toggle inline code snippet"
            onClick={format.code.toggle}
          >
            <CodeIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            kind={format.blockquote.isActive ? "primary" : "ghost"}
            title="Toggle quote"
            onClick={format.blockquote.toggle}
          >
            <TextQuoteIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            kind={format.pre.isActive ? "primary" : "ghost"}
            title="Toggle code block"
            onClick={format.pre.toggle}
          >
            <CodeSquareIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            kind={format.bulletList.isActive ? "primary" : "ghost"}
            title="Toggle bullet points"
            onClick={format.bulletList.toggle}
          >
            <ListIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            kind={format.orderedList.isActive ? "primary" : "ghost"}
            title="Toggle numbered list"
            onClick={format.orderedList.toggle}
          >
            <ListOrderedIcon />
          </Button>

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
