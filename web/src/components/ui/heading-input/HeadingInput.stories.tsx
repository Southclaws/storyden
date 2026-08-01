import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { LStack } from "@/styled-system/jsx";

import { HeadingInput } from ".";

type HeadingInputStoryArgs = {
  defaultValue: string;
  placeholder: string;
  size: "xs" | "sm" | "md" | "lg" | "xl" | "2xl";
};

const meta = {
  title: "UI/Heading Input",
  args: {
    defaultValue: "Editable title",
    placeholder: "Untitled",
    size: "lg",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg", "xl", "2xl"],
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
        size={args.size as any}
        value={value}
        onValueChange={setValue}
      />
    );
  },
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3" maxW="2xl">
      {(["xs", "sm", "md", "lg", "xl", "2xl"] as const).map((size) => (
        <HeadingInput
          key={size}
          size={size as any}
          defaultValue={`Editable heading ${size}`}
          onValueChange={() => undefined}
        />
      ))}
    </LStack>
  ),
};
