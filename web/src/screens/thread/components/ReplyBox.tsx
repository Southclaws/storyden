import { Box, Button, ChakraProvider, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { extended } from "src/theme";

import { AllStyledComponent } from "@remirror/styles/emotion";
import { ExtensionPriority } from "remirror";

import {
  EditorComponent,
  FloatingToolbar,
  Remirror,
  TableComponents,
  TableExtension,
  ThemeProvider,
  useHelpers,
  useRemirror,
} from "@remirror/react";

import {
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  CodeBlockExtension,
  CodeExtension,
  HardBreakExtension,
  HeadingExtension,
  ItalicExtension,
  LinkExtension,
  ListItemExtension,
  MarkdownExtension,
  OrderedListExtension,
  StrikeExtension,
  TrailingNodeExtension,
  wysiwygPreset,
} from "remirror/extensions";

const extensions = () => [
  new LinkExtension({ autoLink: true }),
  new BoldExtension(),
  new StrikeExtension(),
  new ItalicExtension(),
  new HeadingExtension(),
  new LinkExtension(),
  new BlockquoteExtension(),
  new BulletListExtension({ enableSpine: true }),
  new OrderedListExtension(),
  new ListItemExtension({
    priority: ExtensionPriority.High,
    enableCollapsible: true,
  }),
  new CodeExtension(),
  new CodeBlockExtension({ supportedLanguages: [] }),
  new TrailingNodeExtension(),
  new TableExtension(),
  new MarkdownExtension({ copyAsMarkdown: false }),
  new HardBreakExtension(),
  ...wysiwygPreset(),
];

type Props = {
  onSave: (md: string) => void;
};

export function ReplyBox({ onSave }: Props) {
  const onSaveAll = (md: string) => {
    onSave(md);
  };

  return (
    <Box minH={32} width="full" borderRadius="2xl" p={2}>
      <Editor>
        <Save onSave={onSaveAll} />
      </Editor>
    </Box>
  );
}

export function Editor({ children }: PropsWithChildren) {
  const { manager } = useRemirror({
    extensions,
    stringHandler: "markdown",
    content: "**Markdown** content is the _best_",
    selection: "end",
  });

  return (
    <AllStyledComponent style={{ width: "100%", minHeight: "6em" }}>
      <ThemeProvider>
        <Remirror manager={manager}>
          <Flex flexDir="column" width="full" minHeight="6em">
            <EditorComponent />

            <TableComponents />
            <FloatingToolbar />

            {/* NOTE: Remirror is doing some janky stuff and blocking Chakra */}
            <ChakraProvider theme={extended}>{children}</ChakraProvider>
          </Flex>
        </Remirror>
      </ThemeProvider>
    </AllStyledComponent>
  );
}

type SaveProps = { onSave: (md: string) => void };
function Save({ onSave }: SaveProps) {
  const helpers = useHelpers();

  const onClick = () => {
    const md = helpers.getMarkdown();
    if (md.length === 0) return;
    onSave(md);
  };

  return (
    <Flex mt={4} justifyContent="end">
      <Button onClick={onClick}>Post</Button>
    </Flex>
  );
}
