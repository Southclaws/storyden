import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Slider } from ".";

const meta = {
  title: "Components/Forms/Slider",
  component: Slider,
  args: {
    children: "Confidence",
    defaultValue: [64],
    min: 0,
    max: 100,
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof Slider>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="5" maxW="md">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Slider
          key={size}
          size={size}
          defaultValue={[64]}
          marks={[
            { value: 0, label: "0" },
            { value: 50, label: "50" },
            { value: 100, label: "100" },
          ]}
        >
          Slider {size}
        </Slider>
      ))}
    </LStack>
  ),
};
