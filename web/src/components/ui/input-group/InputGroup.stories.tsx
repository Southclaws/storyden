import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { LinkIcon } from "../icons/Link";
import { SearchIcon } from "../icons/Search";
import { Input } from "../input";

import { InputGroup } from ".";

const meta = {
  title: "UI/Input Group",
  component: InputGroup,
  args: {
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg", "xl"],
    },
  },
} satisfies Meta<typeof InputGroup>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {
  render: (args) => (
    <InputGroup {...args} startElement={<SearchIcon />} maxW="md">
      <Input size={args.size} placeholder="Search..." />
    </InputGroup>
  ),
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3" maxW="md">
      {(["xs", "sm", "md", "lg", "xl"] as const).map((size) => (
        <InputGroup key={size} size={size} startElement={<SearchIcon />}>
          <Input size={size} placeholder={`Search ${size}`} />
        </InputGroup>
      ))}
      <InputGroup startElement={<LinkIcon />} endElement=".storyden.org">
        <Input placeholder="community" />
      </InputGroup>
    </LStack>
  ),
};
