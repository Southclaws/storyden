import { Box, Flex, HStack } from "@chakra-ui/react";

import { AllStyledComponent } from "@remirror/styles/emotion";
import { ExtensionPriority } from "remirror";

import {
  BasicFormattingButtonGroup,
  EditorComponent,
  FloatingToolbar,
  HeadingLevelButtonGroup,
  Remirror,
  TableComponents,
  TableExtension,
  ThemeProvider,
  Toolbar,
  VerticalDivider,
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
  onChange: (md: string) => void;
};

export function Editor({ onChange }: Props) {
  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    stringHandler: "markdown",
    selection: "end",
  });

  return (
    <Box minH={32} width="full" borderRadius="2xl" mb={4}>
      <AllStyledComponent style={{ width: "100%", minHeight: "6em" }}>
        <ThemeProvider>
          <Remirror
            manager={manager}
            state={state}
            onChange={(parameter) => {
              setState(parameter.state);

              // We can't use the useHelpers hook because that can only be
              // called from a component that's a child of <Remirror>...
              const ctx = getContext();

              // This assumes the MarkdownExtension is loaded during init.
              const markdownExtension =
                ctx?.manager.getExtension(MarkdownExtension);

              const md = markdownExtension?.getMarkdown(parameter.state);

              // Note: this *looks* like a "controlled component" however it's
              // only half of it, the value from this `onChange` call cannot be
              // passed back into the component because... it's pointless!
              onChange(md as string);
            }}
          >
            <Flex flexDir="column" width="full" minHeight="6em">
              <Toolbar>
                <HStack width="full" justifyContent={["center", "start"]}>
                  <BasicFormattingButtonGroup />
                  <VerticalDivider />
                  <HeadingLevelButtonGroup />
                </HStack>
              </Toolbar>

              <EditorComponent />

              <TableComponents />
              <FloatingToolbar />
            </Flex>
          </Remirror>
        </ThemeProvider>
      </AllStyledComponent>
    </Box>
  );
}
