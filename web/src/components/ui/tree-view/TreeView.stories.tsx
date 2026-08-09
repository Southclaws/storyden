import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { cx } from "@/styled-system/css";
import { treeView } from "@/styled-system/recipes";

const meta = {
  title: "Components/Navigation/Tree View",
  parameters: {
    docs: {
      description: {
        component:
          "A read-only hierarchy that shares Drag Tree's dense visual language without editing or sort behavior. Use it when users need to inspect or navigate a tree but cannot reorganise it.",
      },
    },
  },
  argTypes: {
    variant: {
      control: "select",
      options: ["clamped", "scrollable"],
    },
  },
  args: {
    variant: "clamped",
  },
} satisfies Meta<{ variant: "clamped" | "scrollable" }>;

export default meta;

type Story = StoryObj<{ variant: "clamped" | "scrollable" }>;

function DemoTree(props: { variant: "clamped" | "scrollable" }) {
  const classes = treeView({ variant: props.variant });

  return (
    <div className={classes.root} style={{ maxWidth: 320 }}>
      <div className={classes.tree}>
        <div className={classes.branch} data-depth="1">
          <div className={classes.branchControl} data-depth="1">
            <span className={classes.branchIndicator}>›</span>
            <span className={classes.branchText}>Documentation Hub</span>
          </div>
          <div className={classes.branchContent} data-part="branch-content">
            <div className={classes.item} data-depth="2">
              <span className={classes.itemText}>Advanced Topics</span>
            </div>
            <div className={cx(classes.item)} data-depth="2" data-selected>
              <span className={classes.itemText}>
                Long page title that demonstrates clamped tree text
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export const Playground: Story = {
  args: {
    variant: "clamped",
  },
  render: (args) => <DemoTree {...args} />,
};

export const Variants: Story = {
  render: () => (
    <div style={{ display: "grid", gap: 24 }}>
      <DemoTree variant="clamped" />
      <DemoTree variant="scrollable" />
    </div>
  ),
};
