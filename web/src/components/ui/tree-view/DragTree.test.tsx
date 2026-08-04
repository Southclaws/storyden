import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { DragTree, isDragTreeDescendant } from ".";

type Node = {
  id: string;
  label: string;
  children: Node[];
};

const nodes: Node[] = [
  {
    id: "root",
    label: "Documentation",
    children: [
      {
        id: "guides",
        label: "Guides",
        children: [{ id: "start", label: "Getting started", children: [] }],
      },
    ],
  },
  { id: "about", label: "About", children: [] },
];

const getId = (node: Node) => node.id;
const getLabel = (node: Node) => node.label;
const getHref = (node: Node) => `/docs/${node.id}`;
const getChildren = (node: Node) => node.children;

type DndRegistration = {
  id: string;
  disabled?: boolean;
  data?: Record<string, unknown>;
};

const dnd = vi.hoisted(() => ({
  active: null as null | { data: { current: Record<string, unknown> } },
  over: null as null | { data: { current: Record<string, unknown> } },
  draggables: [] as DndRegistration[],
  droppables: [] as DndRegistration[],
  onKeyDown: vi.fn((event: KeyboardEvent) => event.preventDefault()),
}));

vi.mock("@dnd-kit/core", () => ({
  useDraggable: (options: DndRegistration) => {
    dnd.draggables.push(options);
    return {
      active: dnd.active,
      attributes: {},
      isDragging: false,
      listeners: { onKeyDown: dnd.onKeyDown },
      over: dnd.over,
      setNodeRef: vi.fn(),
      transform: null,
    };
  },
  useDroppable: (options: DndRegistration) => {
    dnd.droppables.push(options);
    return { setNodeRef: vi.fn() };
  },
}));

beforeEach(() => {
  dnd.active = null;
  dnd.over = null;
  dnd.draggables.length = 0;
  dnd.droppables.length = 0;
  dnd.onKeyDown.mockClear();
});

function renderTree(options?: { editable?: boolean }) {
  return render(<TreeFixture {...options} />);
}

function TreeFixture(options: { editable?: boolean }) {
  return (
    <DragTree
      id="documentation"
      label="Documentation"
      nodes={nodes}
      getId={getId}
      getLabel={getLabel}
      getHref={getHref}
      getChildren={getChildren}
      defaultExpandedIds={["root", "guides"]}
      isSelected={(node) => node.id === "start"}
      drag={
        options?.editable
          ? {
              enabled: true,
              getNodeData: ({ node, parent }) => ({
                type: "docs-node",
                nodeId: node.id,
                parentId: parent?.id ?? null,
              }),
              getDividerData: ({ node, parent, direction }) => ({
                type: "docs-divider",
                nodeId: node.id,
                parentId: parent?.id ?? null,
                direction,
              }),
            }
          : undefined
      }
      renderActions={({ node }) => (
        <button type="button">Manage {node.label}</button>
      )}
    />
  );
}

describe("DragTree", () => {
  test("renders nested links, expansion, selection, and actions", () => {
    const { container, rerender } = renderTree();

    expect(screen.getByRole("link", { name: "Documentation" })).toHaveAttribute(
      "href",
      "/docs/root",
    );
    expect(
      screen.getByRole("link", { name: "Getting started" }),
    ).toHaveAttribute("href", "/docs/start");
    expect(
      screen.getByRole("button", { name: "Manage Getting started" }),
    ).toBeVisible();
    expect(
      screen
        .getByRole("link", { name: "Documentation" })
        .closest('[role="treeitem"]'),
    ).toHaveAttribute("aria-expanded", "true");
    expect(
      screen
        .getByRole("link", { name: "Getting started" })
        .closest('[role="treeitem"]'),
    ).toHaveAttribute("aria-selected", "true");
    expect(container.querySelector('[data-part="drop-indicator"]')).toBeNull();

    fireEvent.click(
      screen.getByRole("button", { name: "Toggle Documentation" }),
    );
    expect(
      screen
        .getByRole("link", { name: "Documentation" })
        .closest('[role="treeitem"]'),
    ).toHaveAttribute("aria-expanded", "false");

    rerender(<TreeFixture />);
    expect(
      screen
        .getByRole("link", { name: "Documentation" })
        .closest('[role="treeitem"]'),
    ).toHaveAttribute("aria-expanded", "false");
  });

  test("keeps disclosure keyboard input separate from drag activation", async () => {
    const user = userEvent.setup();
    renderTree({ editable: true });
    const toggle = screen.getByRole("button", {
      name: "Toggle Documentation",
    });
    const treeItem = screen
      .getByRole("link", { name: "Documentation" })
      .closest('[role="treeitem"]');

    toggle.focus();
    await user.keyboard("{Enter}");

    expect(treeItem).toHaveAttribute("aria-expanded", "false");
    await user.keyboard(" ");
    expect(treeItem).toHaveAttribute("aria-expanded", "true");
    expect(dnd.onKeyDown).not.toHaveBeenCalled();
  });

  test("preserves caller drag payloads and renders one divider per insertion point", () => {
    const { container } = renderTree({ editable: true });

    expect(dnd.draggables).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          id: "start",
          data: expect.objectContaining({
            type: "docs-node",
            nodeId: "start",
            parentId: "guides",
            dragTree: expect.objectContaining({
              treeId: "documentation",
              kind: "node",
              nodeId: "start",
            }),
          }),
        }),
      ]),
    );
    expect(
      container.querySelectorAll(
        '[data-part="drop-indicator"][data-direction="above"]',
      ),
    ).toHaveLength(4);
    expect(
      container.querySelectorAll(
        '[data-part="drop-indicator"][data-direction="below"]',
      ),
    ).toHaveLength(3);
  });

  test("disables the active node and every deep descendant as drop targets", () => {
    dnd.active = {
      data: {
        current: {
          dragTree: {
            treeId: "documentation",
            kind: "node",
            nodeId: "root",
          },
        },
      },
    };

    renderTree({ editable: true });

    for (const id of ["root", "guides", "start"]) {
      expect(dnd.droppables.find((options) => options.id === id)).toEqual(
        expect.objectContaining({ disabled: true }),
      );

      for (const direction of ["above", "below"]) {
        const divider = dnd.droppables.find(
          (options) => options.id === `${id}_${direction}`,
        );
        if (divider) {
          expect(divider).toEqual(expect.objectContaining({ disabled: true }));
        }
      }
    }
    expect(dnd.droppables.find((options) => options.id === "about")).toEqual(
      expect.objectContaining({ disabled: false }),
    );
  });

  test("finds descendants below non-root ancestors", () => {
    expect(
      isDragTreeDescendant(nodes, "guides", "start", getId, getChildren),
    ).toBe(true);
    expect(
      isDragTreeDescendant(nodes, "start", "guides", getId, getChildren),
    ).toBe(false);
  });
});
