import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { LStack } from "@/styled-system/jsx";

import { UnreadyBanner } from "./Unready";

const meta = {
  title: "Components/Feedback/Unready Banner",
  component: UnreadyBanner,
  decorators: [
    (Story) => (
      <LStack width="full" maxW="2xl">
        <Story />
      </LStack>
    ),
  ],
  parameters: {
    docs: {
      description: {
        component:
          "A prominent unavailable-content state with a heading, error message, and optional recovery content. With no error, the banner reserves a 384px loading region. Use it when the main content cannot be displayed; use Unready for feedback inside an existing region, Empty State when loading succeeded but the collection is empty, and Admonition for an operation outcome alongside usable content. The banner has one layout size; compact sizing belongs to Unready.",
      },
    },
  },
  argTypes: {
    error: {
      control: "text",
    },
    size: {
      table: {
        disable: true,
      },
    },
    children: {
      control: false,
    },
  },
} satisfies Meta<typeof UnreadyBanner>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Loading: Story = {};

export const Error: Story = {
  args: {
    error: "We could not load this page. Please try again.",
  },
};

export const LongMessage: Story = {
  args: {
    error:
      "We could not load the latest discussions because the connection was interrupted. Check your connection and try again. Your existing content is still available.",
  },
};

export const NarrowError: Story = {
  ...Error,
  render: (args) => (
    <LStack width="64" maxW="full">
      <UnreadyBanner {...args} />
    </LStack>
  ),
};

export const ExceptionMessage: Story = {
  argTypes: {
    error: {
      control: false,
    },
  },
  render: (args) => (
    <UnreadyBanner
      {...args}
      error={new globalThis.Error("The connection timed out. Try again.")}
    />
  ),
};

export const WithRecoveryAction: Story = {
  ...Error,
  parameters: {
    docs: {
      description: {
        story:
          "Children provide the recovery action for the failed operation. Try again switches this preview to the loading state; reset the story to return to the error. This preview does not make a network request.",
      },
    },
  },
  render: (args) => {
    const [retrying, setRetrying] = useState(false);
    const error = retrying ? undefined : args.error;

    return (
      <UnreadyBanner {...args} error={error}>
        <Button onClick={() => setRetrying(true)}>Try again</Button>
      </UnreadyBanner>
    );
  },
};
