import { Portal } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

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
    docs: {
      description: {
        component:
          "Shared editing chrome for block-based screens. The UI layer owns gutters, handles, menus, and theme-ready styling; each domain editor owns block data, permissions, persistence, and the drag payload understood by the global drag-and-drop provider. Read mode must remain visually stable when edit mode is enabled.",
      },
    },
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

export const NestedConfiguration: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Portalled submenus stay inside the block menu's interaction boundary. Select a layout with the mouse or keyboard; clicking outside still dismisses the whole menu.",
      },
    },
  },
  render: function NestedConfiguration() {
    const [layout, setLayout] = useState("list");

    return (
      <Box maxW="2xl" ml="8">
        <BlockEditor.Root className="group">
          <BlockEditor.Gutter>
            <BlockEditor.MenuHandle>
              <Menu.ItemGroup>
                <Menu.ItemGroupLabel>Library</Menu.ItemGroupLabel>
                <Menu.Separator />
                <Menu.Root onSelect={({ value }) => setLayout(value)}>
                  <Menu.TriggerItem>
                    <LayoutIcon />
                    Layout
                  </Menu.TriggerItem>
                  <Portal>
                    <Menu.Positioner>
                      <Menu.Content>
                        <Menu.Item value="list">List</Menu.Item>
                        <Menu.Item value="grid">Grid</Menu.Item>
                      </Menu.Content>
                    </Menu.Positioner>
                  </Portal>
                </Menu.Root>
              </Menu.ItemGroup>
            </BlockEditor.MenuHandle>
          </BlockEditor.Gutter>
          <BlockEditor.Content>
            <CardBox>
              <Text>Library layout: {layout}</Text>
            </CardBox>
          </BlockEditor.Content>
        </BlockEditor.Root>
      </Box>
    );
  },
};
