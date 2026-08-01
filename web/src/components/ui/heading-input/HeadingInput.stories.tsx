import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { LStack } from "@/styled-system/jsx";

import { HeadingInput } from ".";

type HeadingInputStoryArgs = {
  defaultValue: string;
  placeholder: string;
  size: "md" | "lg";
};

const meta = {
  title: "UI/Heading Input",
  args: {
    defaultValue: "Editable title",
    placeholder: "Untitled",
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["md", "lg"],
    },
  },
} satisfies Meta<HeadingInputStoryArgs>;

export default meta;

type Story = StoryObj<HeadingInputStoryArgs>;

export const Playground: Story = {
  render: (args) => {
    const [value, setValue] = useState(String(args.defaultValue ?? ""));

    return (
      <HeadingInput
        {...args}
        size={args.size}
        value={value}
        onValueChange={setValue}
      />
    );
  },
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3" maxW="2xl">
      {(["md", "lg"] as const).map((size) => (
        <HeadingInput
          key={size}
          size={size}
          defaultValue={`Editable heading ${size}`}
          onValueChange={() => undefined}
        />
      ))}
    </LStack>
  ),
};
