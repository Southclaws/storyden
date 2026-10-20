import { createTreeCollection } from "@ark-ui/react/tree-view";
import {
  DndContext,
  MouseSensor,
  TouchSensor,
  pointerWithin,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { Fragment, type ReactNode, useState } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { MoreIcon } from "@/components/ui/icons/More";
import * as MenuPrimitive from "@/components/ui/menu";
import { Menu } from "@/components/ui/menu/Menu";
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
  title: "Components/Navigation/Drag Tree",
  args: { variant: "clamped" as const },
  argTypes: {
    variant: { control: "select", options: ["clamped", "scrollable"] },
  },
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "The draggable navigation hierarchy used for Library pages and discussion categories. Each tree item owns its destination link and row actions, for both leaves and branches; its accessible name comes from the destination title. Ark owns focus, hierarchy, and typeahead: Up/Down move between visible nodes; Right opens a branch or enters its first child; Left closes a branch or returns to its parent; Home/End jump to the first/last visible node. Typing a title prefix moves focus without navigating. Enter opens the focused destination. Clicking a title preserves branch expansion; the chevron toggles it separately. Tab reaches only the focused row's actions; Escape from those actions returns to its destination. Action renderers must forward the supplied tabIndex to each control, including menu triggers. The current destination uses aria-current=page. dnd-kit owns pointer dragging; its keyboard activator is excluded so tree navigation cannot start a drag. Callers provide domain data, persistence, links, actions, and a route-agnostic selected matcher. Links do not prefetch the hierarchy's destinations. Use Tree View's visual recipe for a static hierarchy without navigation or editing controls. Use the variant control to compare clamped and scrollable layouts.",
      },
    },
  },
  decorators: [
    (Story) => (
      <styled.div width="72" minHeight="80">
        <Story />
      </styled.div>
    ),
  ],
} satisfies Meta<{ variant: "clamped" | "scrollable" }>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Nested: Story = {
  render: (args) => (
    <DragTreeStoryContext>
      <DragTree
        variant={args.variant}
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
  render: (args) => (
    <DragTreeStoryContext>
      <DragTree
        variant={args.variant}
        id="story-editable"
        label="Documentation"
        nodes={demoNodes}
        getId={getId}
        getLabel={getLabel}
        getHref={getHref}
        getChildren={getChildren}
        defaultExpandedIds={["documentation", "guides", "community"]}
        isSelected={(node) => node.id === "theming"}
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
        renderActions={({ node, tabIndex }) => (
          <Fragment>
            <IconButton
              tabIndex={tabIndex}
              aria-label={`Add child to ${node.label}`}
              variant="plain"
            >
              <AddIcon />
            </IconButton>
            <IconButton
              tabIndex={tabIndex}
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
  parameters: {
    docs: {
      description: {
        story:
          "Themes and appearance is the current page. Arrow to another node to compare its full-row focus outline with the current page's persistent background. Tab into the focused row's actions to see each button's own focus indicator.",
      },
    },
  },
};

export const LongLabels: Story = {
  render: (args) => (
    <DragTreeStoryContext>
      <DragTree
        variant={args.variant}
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

function KeyboardMovingExample({
  variant,
  fail = false,
}: {
  variant: "clamped" | "scrollable";
  fail?: boolean;
}) {
  const [nodes, setNodes] = useState(demoNodes);
  return (
    <DragTreeStoryContext>
      <DragTree
        id="story-keyboard-moving"
        label="Documentation"
        nodes={nodes}
        variant={variant}
        getId={getId}
        getLabel={getLabel}
        getChildren={getChildren}
        getHref={getHref}
        isSelected={(node) => node.id === "theming"}
        defaultExpandedIds={["documentation", "guides", "community"]}
        drag={{
          enabled: true,
          getNodeData: ({ node }) => ({ type: "story-node", node }),
          getDividerData: ({ node, direction }) => ({
            type: "story-divider",
            node,
            direction,
          }),
        }}
        move={{
          onMove: async ({ node, target, placement }) => {
            if (fail) throw new Error("Example persistence failure");
            const collection = createTreeCollection({
              rootNode: { children: nodes } as DemoNode,
              nodeToValue: getId,
              nodeToString: getLabel,
              nodeToChildren: getChildren,
            });
            const without = collection.remove([
              collection.getIndexPath(node.id)!,
            ]);
            const path = without.getIndexPath(target.id)!;
            const result =
              placement === "inside"
                ? without.replace(path, {
                    ...target,
                    children: [...target.children, node],
                  })
                : placement === "above"
                  ? without.insertBefore(path, [node])!
                  : without.insertAfter(path, [node])!;
            setNodes(result.rootNode.children);
          },
        }}
        renderActions={({ node, tabIndex, startMove }) => (
          <Menu
            onSelect={({ value }) => {
              if (value === "move") startMove?.();
            }}
            trigger={
              <IconButton
                variant="plain"
                tabIndex={tabIndex}
                aria-label={`Actions for ${node.label}`}
              >
                <MoreIcon />
              </IconButton>
            }
          >
            <MenuPrimitive.Item value="move">
              Move or reorder…
            </MenuPrimitive.Item>
          </Menu>
        )}
      />
    </DragTreeStoryContext>
  );
}

export const KeyboardMoving: Story = {
  render: (args) => <KeyboardMovingExample variant={args.variant} />,
  parameters: {
    docs: {
      description: {
        story:
          "Focus a node with the keyboard to discover Alt+Shift+M (Option+Shift+M on Mac), or choose Move or reorder in its menu. The node and its descendants are carried together. Arrows, Home/End, and typeahead browse destinations; Space cycles before/inside/after. Enter saves the placement, Escape restores the original position, and Tab cancels without trapping focus. Before/after uses the same focus-coloured divider as pointer sorting, without outlining the destination row. Before/after a root node returns the carried node to the top level. The source stays in place until confirmation. The shortcut hint closes when focus leaves the tree, including Shift+Tab back to the preceding control. The preview respects reduced motion. This example persists moves in local state.",
      },
    },
  },
};

export const MoveFailure: Story = {
  render: (args) => <KeyboardMovingExample variant={args.variant} fail />,
  parameters: {
    docs: {
      description: {
        story:
          "An unsuccessful move retains the carried node and destination for retry. The inline error remains available to assistive technology. Escape still cancels and returns focus to the source.",
      },
    },
  },
};
