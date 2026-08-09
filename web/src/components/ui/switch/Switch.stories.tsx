import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Switch } from ".";

const meta = {
  title: "Components/Forms/Switch",
  component: Switch,
  parameters: {
    docs: {
      description: {
        component:
          "An immediate on/off setting whose effect is clear from its label. Use Checkbox when values are staged for form submission or several independent options are selected together. The track, thumb, label, focus, and disabled states must remain distinct in both themes.",
      },
    },
  },
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
