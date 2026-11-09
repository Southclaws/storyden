import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box, LStack } from "@/styled-system/jsx";

import { Unready } from "./Unready";

const meta = {
  title: "Components/Feedback/Unready",
  component: Unready,
  decorators: [
    (Story) => (
      <LStack width="full" maxW="lg">
        <Story />
      </LStack>
    ),
  ],
  parameters: {
    docs: {
      description: {
        component:
          "Loading or failure feedback for an existing content region. With no error, Unready centres a labelled spinner inside its parent; give that parent a height. With an error, it shows the derived message and warning icon. Use sm for compact controls and dropdowns, and md for ordinary content regions. Use Unready Banner when the unavailable region needs a heading and recovery actions, Empty State for an empty collection, and field error text for invalid input.",
      },
    },
  },
  args: {
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md"],
    },
    error: {
      control: "text",
    },
  },
} satisfies Meta<typeof Unready>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Loading: Story = {
  render: (args) => (
    <Box width="full" height="24">
      <Unready {...args} />
    </Box>
  ),
};

export const CompactLoading: Story = {
  ...Loading,
  args: {
    size: "sm",
  },
};

export const Error: Story = {
  args: {
    error: "Capabilities could not be loaded. Try again.",
  },
};

export const CompactError: Story = {
  args: {
    size: "sm",
    error: "Capabilities could not be loaded. Try again.",
  },
};

export const NarrowError: Story = {
  ...Error,
  render: (args) => (
    <LStack width="64" maxW="full">
      <Unready {...args} />
    </LStack>
  ),
};

export const NarrowCompactError: Story = {
  ...NarrowError,
  args: {
    size: "sm",
    error: "Capabilities could not be loaded. Try again.",
  },
};

export const LongMessage: Story = {
  args: {
    error:
      "We could not load the latest discussions because the connection was interrupted. Check your connection and try again. Your existing content is still available.",
  },
};

export const UnbreakableMessage: Story = {
  ...NarrowCompactError,
  args: {
    size: "sm",
    error:
      "Could not load discussion-archive-with-a-very-long-unbroken-identifier-0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz.",
  },
};

export const ExceptionMessage: Story = {
  argTypes: {
    error: {
      control: false,
    },
  },
  render: (args) => (
    <Unready
      {...args}
      error={new globalThis.Error("The connection timed out. Try again.")}
    />
  ),
};
