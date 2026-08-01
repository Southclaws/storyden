import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import * as Tabs from ".";

const meta = {
  title: "UI/Tabs",
  component: Tabs.Root,
  args: {
    defaultValue: "overview",
    size: "md",
    variant: "line",
    orientation: "horizontal",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["line", "enclosed", "outline"],
    },
    orientation: {
      control: "select",
      options: ["horizontal", "vertical"],
    },
  },
} satisfies Meta<typeof Tabs.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: Tabs.RootProps) {
  return (
    <Tabs.Root {...props}>
      <Tabs.List>
        <Tabs.Trigger value="overview">Overview</Tabs.Trigger>
        <Tabs.Trigger value="activity">Activity</Tabs.Trigger>
        <Tabs.Trigger value="settings">Settings</Tabs.Trigger>
        <Tabs.Indicator />
      </Tabs.List>
      <Tabs.Content value="overview">
        Community summary and recent work.
      </Tabs.Content>
      <Tabs.Content value="activity">
        Recent posts, replies, and edits.
      </Tabs.Content>
      <Tabs.Content value="settings">
        Configuration and access controls.
      </Tabs.Content>
    </Tabs.Root>
  );
}

export const Playground: Story = {
  render: (args) => <Example {...args} />,
};

export const Variants: Story = {
  render: () => (
    <LStack gap="8">
      {(["line", "enclosed", "outline"] as const).map((variant) => (
        <Example
          key={variant}
          variant={variant}
          size="md"
          defaultValue="overview"
        />
      ))}
      <Example
        orientation="vertical"
        variant="line"
        size="sm"
        defaultValue="overview"
      />
    </LStack>
  ),
};
