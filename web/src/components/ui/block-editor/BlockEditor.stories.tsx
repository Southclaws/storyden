import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { CardBox } from "@/components/ui/card-box";
import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { Text } from "@/components/ui/text";
import { Box } from "@/styled-system/jsx";

import * as BlockEditor from ".";

const meta = {
  title: "UI/Block Editor",
  component: BlockEditor.Root,
  parameters: {
    layout: "padded",
  },
} satisfies Meta<typeof BlockEditor.Root>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Item: Story = {
  render: () => (
    <Box maxW="2xl" ml="8">
      <BlockEditor.Root className="group">
        <BlockEditor.Gutter>
          <BlockEditor.Handle>
            <IconButton aria-label="Move or configure block" variant="subtle">
              <DragHandleIcon />
            </IconButton>
          </BlockEditor.Handle>
        </BlockEditor.Gutter>
        <BlockEditor.Content>
          <CardBox>
            <Text>Discussion categories</Text>
          </CardBox>
        </BlockEditor.Content>
      </BlockEditor.Root>
    </Box>
  ),
};
