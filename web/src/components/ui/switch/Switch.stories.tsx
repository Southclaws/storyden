import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Switch } from ".";

const meta = {
  title: "UI/Switch",
  component: Switch,
  args: {
    children: "Enable notifications",
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof Switch>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Switch key={size} size={size} defaultChecked>
          Switch {size}
        </Switch>
      ))}
    </LStack>
  ),
};
