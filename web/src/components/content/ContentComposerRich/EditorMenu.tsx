import { Portal } from "@ark-ui/react";
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
import { HStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { LinkButton } from "./LinkButton";

export type Props = {
  editor: any;
  uniqueID: string;
  format: any;
  handlers: any;
};

export function EditorMenu({ editor, uniqueID, format, handlers }: Props) {
  return (
    <HStack gap="1">
      <Menu.Root
        onSelect={(d) => format.text.set(d.value as any /* lazy */)}
        positioning={{ gutter: 0 }}
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
    </HStack>
  );
}
