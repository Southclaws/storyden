import {
  DndContext,
  MouseSensor,
  TouchSensor,
  pointerWithin,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { Fragment, type ReactNode } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { MoreIcon } from "@/components/ui/icons/More";
import { styled } from "@/styled-system/jsx";

import { DragTree } from ".";

type DemoNode = {
  id: string;
  label: string;
  children: DemoNode[];
};

const demoNodes: DemoNode[] = [
  {
    id: "documentation",
    label: "Documentation",
    children: [
      { id: "getting-started", label: "Getting started", children: [] },
      {
        id: "guides",
        label: "Guides",
        children: [
          { id: "moderation", label: "Moderation", children: [] },
          { id: "theming", label: "Themes and appearance", children: [] },
        ],
      },
    ],
  },
  {
    id: "community",
    label: "Community",
    children: [
      { id: "conduct", label: "Code of conduct", children: [] },
      { id: "contributing", label: "Contributing", children: [] },
    ],
  },
];

const getId = (node: DemoNode) => node.id;
const getLabel = (node: DemoNode) => node.label;
const getHref = (node: DemoNode) => `#${node.id}`;
const getChildren = (node: DemoNode) => node.children;

function DragTreeStoryContext({ children }: { children: ReactNode }) {
  const sensors = useSensors(
    useSensor(MouseSensor, {
      activationConstraint: { distance: 2 },
    }),
    useSensor(TouchSensor, {
      activationConstraint: { delay: 150, tolerance: 5 },
    }),
  );

  return (
    <DndContext collisionDetection={pointerWithin} sensors={sensors}>
      {children}
    </DndContext>
  );
}

const meta = {
  title: "UI/Drag Tree",
  parameters: {
    layout: "centered",
  },
  decorators: [
    (Story) => (
      <styled.div width="72" minHeight="80">
        <Story />
      </styled.div>
    ),
  ],
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const Nested: Story = {
  render: () => (
    <DragTreeStoryContext>
      <DragTree
        id="story-nested"
        label="Documentation"
        nodes={demoNodes}
        getId={getId}
        getLabel={getLabel}
        getHref={getHref}
        getChildren={getChildren}
        defaultExpandedIds={["documentation", "guides", "community"]}
        isSelected={(node) => node.id === "theming"}
      />
    </DragTreeStoryContext>
  ),
};

export const Editable: Story = {
  render: () => (
    <DragTreeStoryContext>
      <DragTree
        id="story-editable"
        label="Documentation"
        nodes={demoNodes}
        getId={getId}
        getLabel={getLabel}
        getHref={getHref}
        getChildren={getChildren}
        defaultExpandedIds={["documentation", "guides", "community"]}
        drag={{
          enabled: true,
          getNodeData: ({ node, parent }) => ({
            type: "story-node",
            nodeId: node.id,
            parentId: parent?.id ?? null,
          }),
          getDividerData: ({ node, parent, direction }) => ({
            type: "story-divider",
            nodeId: node.id,
            parentId: parent?.id ?? null,
            direction,
          }),
        }}
        renderActions={({ node }) => (
          <Fragment>
            <IconButton
              aria-label={`Add child to ${node.label}`}
              variant="plain"
            >
              <AddIcon />
            </IconButton>
            <IconButton
              aria-label={`Actions for ${node.label}`}
              variant="plain"
            >
              <MoreIcon />
            </IconButton>
          </Fragment>
        )}
      />
    </DragTreeStoryContext>
  ),
};

export const LongLabels: Story = {
  render: () => (
    <DragTreeStoryContext>
      <DragTree
        id="story-long-labels"
        label="Documentation"
        nodes={[
          {
            id: "long-root",
            label: "A section with a deliberately long navigation label",
            children: [
              {
                id: "long-child",
                label:
                  "A deeply nested page whose title must remain constrained",
                children: [],
              },
            ],
          },
        ]}
        getId={getId}
        getLabel={getLabel}
        getHref={getHref}
        getChildren={getChildren}
        defaultExpandedIds={["long-root"]}
      />
    </DragTreeStoryContext>
  ),
};
