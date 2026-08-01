import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { Button } from "../button";

import { Admonition } from ".";

const meta = {
  title: "UI/Admonition",
  component: Admonition,
  args: {
    value: true,
    kind: "neutral",
    title: "Community note",
    children: "A short message with a close affordance.",
  },
  argTypes: {
    kind: {
      control: "select",
      options: ["neutral", "success", "failure"],
    },
  },
} satisfies Meta<typeof Admonition>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {
  render: (args) => {
    const [visible, setVisible] = useState(true);

    return visible ? (
      <Admonition {...args} value={visible} onChange={setVisible} />
    ) : (
      <Button size="sm" variant="outline" onClick={() => setVisible(true)}>
        Show admonition
      </Button>
    );
  },
};

export const Variants: Story = {
  render: () => (
    <>
      <Admonition value kind="neutral" title="Neutral">
        Used for low-priority contextual information.
      </Admonition>
      <Admonition value kind="success" title="Success">
        Used when a setup step or moderation action succeeds.
      </Admonition>
      <Admonition value kind="failure" title="Failure">
        Used for destructive or blocking states.
      </Admonition>
    </>
  ),
};
