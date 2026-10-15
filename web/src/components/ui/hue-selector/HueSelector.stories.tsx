import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { HueSelector } from ".";

const meta = {
  title: "Components/Forms/Hue Selector",
  component: HueSelector,
  args: {
    value: "#31a89f",
    onChange: () => undefined,
    onUpdate: () => undefined,
  },
} satisfies Meta<typeof HueSelector>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Drag the hue control to preview a colour; release to commit it. Existing CSS colours are parsed with the shared Polished HSL parser.",
      },
    },
  },
  render: function HueSelectorStory() {
    const [value, setValue] = useState("#31a89f");

    return (
      <HueSelector value={value} onChange={setValue} onUpdate={setValue} />
    );
  },
};
