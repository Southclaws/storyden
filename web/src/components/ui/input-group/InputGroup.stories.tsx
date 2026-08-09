import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { LinkIcon } from "../icons/Link";
import { SearchIcon } from "../icons/Search";
import { Input } from "../input";

import { InputGroup } from ".";

const meta = {
  title: "Components/Forms/Input Group",
  component: InputGroup,
  args: {
    size: "sm",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof InputGroup>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {
  render: (args) => {
    const size = args.size === "md" || args.size === "lg" ? args.size : "sm";

    return (
      <InputGroup {...args} startElement={<SearchIcon />} maxW="md">
        <Input size={size} placeholder="Search..." />
      </InputGroup>
    );
  },
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3" maxW="md">
      {(["sm", "md", "lg"] as const).map((size) => (
        <InputGroup key={size} size={size} startElement={<SearchIcon />}>
          <Input size={size} placeholder={`Search ${size}`} />
        </InputGroup>
      ))}
      <InputGroup startElement={<LinkIcon />} endElement=".storyden.org">
        <Input aria-label="Community address" placeholder="community" />
      </InputGroup>
    </LStack>
  ),
};
