import { Portal } from "@ark-ui/react";
import { EditorContent } from "@tiptap/react";
import { match } from "ts-pattern";

import { Button } from "@/components/ui/button";
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
} from "@/components/ui/icons/Typography";
import * as Menu from "@/components/ui/menu";
import { css, cx } from "@/styled-system/css";
import { LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import "./styles.css";

import { LinkButton } from "./LinkButton";
import { FloatingMenu } from "./plugins/MenuPlugin";
import { ContentComposerProps, useContentComposer } from "./useContentComposer";

export function ContentComposer(props: ContentComposerProps) {
  const { editor, initialValueHTML, uniqueID, handlers, format } =
    useContentComposer(props);

  return (
    <LStack
      id={`rich-text-editor-${uniqueID}`}
      containerType="inline-size"
      className={cx("typography", props.className)}
      // NOTE: Relative positioning is for the floating menu to work.
      position="relative"
      w="full"
      gap="1"
      onDragOver={(e) => e.preventDefault()}
    >
      {editor ? (
        <FloatingMenu editor={editor}>
          <Menu.Root
            onSelect={(d) => format.text.set(d.value as any /* lazy */)}
          >
            <Menu.Trigger asChild>
              <Button
                type="button"
                size="xs"
                variant="ghost"
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
            </Menu.Trigger>

            <Portal>
              {/* NOTE: Because this is a portal, we need to reference this ID
              in the FloatingMenu unfocus logic so we don't hide the menu when
              this menu is opened as it's a portal and not a child element. */}
              <Menu.Positioner>
                <Menu.Content
                  id="text-block-menu"
                  userSelect="none"
                  backdropBlur="md"
                  backdropFilter="auto"
                >
                  <Menu.Item value="p">
                    <TextIcon />
                    &nbsp;Paragraph
                  </Menu.Item>

                  <Menu.Item value="h1" fontSize="lg" fontWeight="extrabold">
                    <Heading1Icon />
                    &nbsp;Heading 1
                  </Menu.Item>
                  <Menu.Item value="h2" fontSize="md" fontWeight="extrabold">
                    <Heading2Icon />
                    &nbsp;Heading 2
                  </Menu.Item>
                  <Menu.Item value="h3" fontSize="md" fontWeight="bold">
                    <Heading3Icon />
                    &nbsp;Heading 3
                  </Menu.Item>
                  <Menu.Item value="h4" fontSize="md" fontWeight="medium">
                    <Heading4Icon />
                    &nbsp;Heading 4
                  </Menu.Item>
                  <Menu.Item value="h5" fontSize="sm" fontWeight="normal">
                    <Heading5Icon />
                    &nbsp;Heading 5
                  </Menu.Item>
                  <Menu.Item value="h6" fontSize="sm" fontWeight="light">
                    <Heading6Icon />
                    &nbsp;Heading 6
                  </Menu.Item>
                </Menu.Content>
              </Menu.Positioner>
            </Portal>
          </Menu.Root>
          <Button
            type="button"
            size="xs"
            variant={format.bold.isActive ? "subtle" : "ghost"}
            title="Toggle bold text"
            onClick={format.bold.toggle}
          >
            <BoldIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            variant={format.italic.isActive ? "subtle" : "ghost"}
            title="Toggle italic text"
            onClick={format.italic.toggle}
          >
            <ItalicIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            variant={format.strike.isActive ? "subtle" : "ghost"}
            title="Toggle strikeout text"
            onClick={format.strike.toggle}
          >
            <StrikethroughIcon />
          </Button>
          <Button
            type="button"
            size="xs"
            variant={format.code.isActive ? "subtle" : "ghost"}
            title="Toggle inline code snippet"
            onClick={format.code.toggle}
          >
            <CodeIcon />
          </Button>

          <LinkButton editor={editor} />

          <Button
            type="button"
            size="xs"
            variant={format.blockquote.isActive ? "subtle" : "ghost"}
            title="Toggle quote"
            onClick={format.blockquote.toggle}
          >
            <TextQuoteIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            variant={format.pre.isActive ? "subtle" : "ghost"}
            title="Toggle code block"
            onClick={format.pre.toggle}
          >
            <CodeSquareIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            variant={format.bulletList.isActive ? "subtle" : "ghost"}
            title="Toggle bullet points"
            onClick={format.bulletList.toggle}
          >
            <ListIcon />
          </Button>

          <Button
            type="button"
            size="xs"
            variant={format.orderedList.isActive ? "subtle" : "ghost"}
            title="Toggle numbered list"
            onClick={format.orderedList.toggle}
          >
            <ListOrderedIcon />
          </Button>

          <label
            className={button({
              size: "xs",
              variant: "ghost",
            })}
            htmlFor={`filepicker-${uniqueID}`}
            title="Insert an image"
          >
            <ImageIcon />
          </label>
          <styled.input
            id={`filepicker-${uniqueID}`}
            type="file"
            multiple
            display="none"
            onChange={handlers.handleFileUpload}
          />
        </FloatingMenu>
      ) : (
        <div dangerouslySetInnerHTML={{ __html: initialValueHTML }} />
      )}

      <EditorContent
        id={`editor-content-${uniqueID}`}
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
