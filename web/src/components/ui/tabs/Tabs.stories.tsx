import { LocaleProvider, createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { CheckIcon } from "../icons/Check";
import { SelectIcon } from "../icons/Select";
import * as Select from "../select";

import * as Tabs from ".";

const meta = {
  title: "Components/Navigation/Tabs",
  component: Tabs.Root,
  parameters: {
    docs: {
      description: {
        component:
          "Use Tabs for a small, stable set of peer views within one product feature, such as Robots and Sessions. Use SectionNavigation for routes and hierarchies instead. Do not hide content whose structure should remain visible without tab state. For horizontal tabs with adjacent filters or actions, wrap List and Controls in Bar. Controls stay outside the tablist, stack below it on small screens, and share the line variant's full-width divider. Lists scroll independently and fade only at edges with more tabs; keyboard focus reveals offscreen tabs.",
      },
    },
  },
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

const categories = createListCollection({
  items: [
    { label: "All categories", value: "all" },
    {
      label: "Community questions and project discussions",
      value: "community",
    },
  ],
});

function WithControlsExample({
  manyTabs = false,
  ...props
}: Tabs.RootProps & { manyTabs?: boolean }) {
  const views = manyTabs
    ? ["Overview", "Activity", "Discussions", "Members", "Events", "Settings"]
    : ["Unread", "Unresolved", "Unanswered"];

  return (
    <Tabs.Root {...props} defaultValue={views[0]} orientation="horizontal">
      <Tabs.Bar>
        <Tabs.List aria-label="Community views">
          {views.map((view) => (
            <Tabs.Trigger key={view} value={view}>
              {view}
            </Tabs.Trigger>
          ))}
          <Tabs.Indicator />
        </Tabs.List>
        <Tabs.Controls>
          <Select.Root collection={categories} defaultValue={["all"]}>
            <Select.Control>
              <Select.Trigger aria-label="Filter by category">
                <Select.ValueText />
                <SelectIcon />
              </Select.Trigger>
            </Select.Control>
            <Select.Positioner>
              <Select.Content>
                {categories.items.map((item) => (
                  <Select.Item key={item.value} item={item}>
                    <Select.ItemText>{item.label}</Select.ItemText>
                    <Select.ItemIndicator>
                      <CheckIcon />
                    </Select.ItemIndicator>
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Positioner>
            <Select.HiddenSelect />
          </Select.Root>
        </Tabs.Controls>
      </Tabs.Bar>
      {views.map((view) => (
        <Tabs.Content key={view} value={view}>
          {view} discussions for the selected category.
        </Tabs.Content>
      ))}
    </Tabs.Root>
  );
}

export const WithControls: Story = {
  render: (args) => <WithControlsExample {...args} />,
};

export const Overflow: Story = {
  render: (args) => (
    <LStack maxW="lg">
      <WithControlsExample {...args} manyTabs />
    </LStack>
  ),
};

export const PartialOverflow: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "A partly visible tab must scroll fully into view when focused. Use End and Home to check both edges without activating another view.",
      },
    },
  },
  render: (args) => (
    <LStack maxW="72">
      <WithControlsExample {...args} size="sm" activationMode="manual" />
    </LStack>
  ),
};

export const RightToLeftOverflow: Story = {
  render: (args) => (
    <LocaleProvider locale="ar">
      <LStack maxW="lg" dir="rtl">
        <WithControlsExample {...args} manyTabs />
      </LStack>
    </LocaleProvider>
  ),
};
