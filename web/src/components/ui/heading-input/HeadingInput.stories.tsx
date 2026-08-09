import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { LStack } from "@/styled-system/jsx";

import { HeadingInput } from ".";

type HeadingInputStoryArgs = {
  defaultValue: string;
  placeholder: string;
};

const meta = {
  title: "Components/Forms/Heading Input",
  args: {
    defaultValue: "Editable title",
    placeholder: "Untitled",
  },
} satisfies Meta<HeadingInputStoryArgs>;

export default meta;

type Story = StoryObj<HeadingInputStoryArgs>;

export const Playground: Story = {
  render: (args) => {
    const [value, setValue] = useState(String(args.defaultValue ?? ""));

    return <HeadingInput {...args} value={value} onValueChange={setValue} />;
  },
};

export const Example: Story = {
  render: () => (
    <LStack gap="3" maxW="2xl">
      <HeadingInput
        defaultValue="Editable page heading"
        onValueChange={() => undefined}
      />
    </LStack>
  ),
};
