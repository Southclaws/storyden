import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { CardBox } from "@/components/ui/card-box";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LayoutIcon } from "@/components/ui/icons/Layout";
import * as Menu from "@/components/ui/menu";
import { Text } from "@/components/ui/text";
import { Box } from "@/styled-system/jsx";

import * as BlockEditor from ".";

const meta = {
  title: "Components/Layout/Block Editor",
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
          <BlockEditor.MenuHandle>
            <Menu.ItemGroup>
              <Menu.ItemGroupLabel>Discussion categories</Menu.ItemGroupLabel>
              <Menu.Separator />
              <Menu.Item value="layout">
                <LayoutIcon />
                Layout
              </Menu.Item>
              <Menu.Item value="delete">
                <DeleteIcon />
                Delete
              </Menu.Item>
            </Menu.ItemGroup>
          </BlockEditor.MenuHandle>
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
