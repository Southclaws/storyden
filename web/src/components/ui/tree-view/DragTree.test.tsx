import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, test, vi } from "vitest";

import * as Menu from "../menu";

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

type FixtureOptions = { editable?: boolean; menuActions?: boolean };

function renderTree(options?: FixtureOptions) {
  return render(<TreeFixture {...options} />);
}

function TreeFixture(options: FixtureOptions) {
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
      renderActions={({ node, tabIndex }) =>
        options.menuActions ? (
          <Menu.Root>
            <Menu.Trigger asChild>
              <button type="button" tabIndex={tabIndex}>
                Manage {node.label}
              </button>
            </Menu.Trigger>
            <Portal>
              <Menu.Positioner>
                <Menu.Content>
                  <Menu.Item value="inspect">Inspect</Menu.Item>
                </Menu.Content>
              </Menu.Positioner>
            </Portal>
          </Menu.Root>
        ) : (
          <button type="button" tabIndex={tabIndex}>
            Manage {node.label}
          </button>
        )
      }
    />
  );
}

describe("DragTree", () => {
  test("renders nested links, expansion, selection, and actions", async () => {
    const { container, rerender } = renderTree();

    expect(screen.getByRole("link", { name: "Documentation" })).toHaveAttribute(
      "href",
      "/docs/root",
    );
    const selectedLink = screen.getByRole("treeitem", {
      name: "Getting started",
    });
    expect(selectedLink).toHaveAttribute("href", "/docs/start");
    expect(selectedLink).toHaveAttribute("aria-current", "page");
    expect(screen.getByRole("treeitem", { name: "About" })).not.toHaveAttribute(
      "aria-current",
    );
    expect(
      screen.getByRole("button", { name: "Manage Getting started" }),
    ).toBeVisible();
    expect(screen.getByRole("tree", { name: "Documentation" })).toBeVisible();
    expect(
      screen.getByRole("treeitem", { name: "Documentation" }),
    ).toHaveAccessibleName("Documentation");
    expect(selectedLink).not.toHaveAttribute("aria-expanded");
    const toggle = screen.getByRole("button", { name: "Toggle Documentation" });
    expect(toggle).toHaveAttribute("aria-expanded", "true");
    expect(
      screen.getByRole("treeitem", { name: "Documentation" }),
    ).toContainElement(selectedLink);
    expect(
      screen.queryByRole("button", { name: "Toggle Getting started" }),
    ).not.toBeInTheDocument();
    expect(container.querySelector('[data-part="drop-indicator"]')).toBeNull();

    fireEvent.click(
      screen.getByRole("button", { name: "Toggle Documentation" }),
    );
    await waitFor(() =>
      expect(toggle).toHaveAttribute("aria-expanded", "false"),
    );
    await waitFor(() => expect(selectedLink).not.toBeVisible());

    rerender(<TreeFixture />);
    expect(toggle).toHaveAttribute("aria-expanded", "false");
  });

  test("keeps disclosure keyboard input separate from drag activation", async () => {
    const user = userEvent.setup();
    renderTree({ editable: true });
    const toggle = screen.getByRole("button", {
      name: "Toggle Documentation",
    });

    toggle.focus();
    await user.keyboard("{Enter}");

    expect(toggle).toHaveAttribute("aria-expanded", "false");
    expect(toggle).toHaveFocus();
    await user.keyboard(" ");
    expect(toggle).toHaveAttribute("aria-expanded", "true");
    expect(toggle).toHaveFocus();
    expect(dnd.onKeyDown).not.toHaveBeenCalled();
  });

  test("uses Ark arrow navigation, branch expansion, parent navigation, Home and End", async () => {
    const user = userEvent.setup();
    renderTree({ editable: true });
    const documentation = screen.getByRole("link", { name: "Documentation" });
    const guides = screen.getByRole("link", { name: "Guides" });
    const start = screen.getByRole("treeitem", { name: "Getting started" });
    const about = screen.getByRole("treeitem", { name: "About" });
    await user.tab();
    expect(start).toHaveFocus();
    await user.keyboard("{ArrowUp}");
    expect(guides).toHaveFocus();
    await user.keyboard("{ArrowLeft}");
    expect(guides).toHaveAttribute("aria-expanded", "false");
    await waitFor(() => expect(start).not.toBeVisible());
    await user.keyboard("{ArrowLeft}");
    expect(documentation).toHaveFocus();
    await user.keyboard("{ArrowRight}");
    expect(guides).toHaveFocus();
    await user.keyboard("{ArrowRight}");
    expect(guides).toHaveAttribute("aria-expanded", "true");
    expect(guides).toHaveFocus();
    await user.keyboard("{ArrowRight}");
    expect(start).toHaveFocus();
    await user.keyboard("{ArrowDown}");
    expect(about).toHaveFocus();
    await user.keyboard("{ArrowDown}");
    expect(about).toHaveFocus();
    await user.keyboard("{Home}");
    expect(documentation).toHaveFocus();
    await user.keyboard("{End}");
    expect(about).toHaveFocus();
    expect(dnd.onKeyDown).not.toHaveBeenCalled();
  });

  test("typeahead accepts lowercase a without navigating or changing the current page", async () => {
    const user = userEvent.setup();
    renderTree();
    await user.tab();
    await user.keyboard("abo");
    expect(screen.getByRole("treeitem", { name: "About" })).toHaveFocus();
    expect(
      screen.getByRole("treeitem", { name: "Getting started" }),
    ).toHaveAttribute("aria-current", "page");
  });

  test("typeahead skips descendants of a collapsed branch", async () => {
    const user = userEvent.setup();
    renderTree();
    const documentation = screen.getByRole("link", { name: "Documentation" });
    documentation.focus();
    await user.keyboard("{ArrowLeft}");
    await waitFor(() =>
      expect(
        screen.queryByRole("link", { name: "Guides" }),
      ).not.toBeInTheDocument(),
    );
    await user.keyboard("gui");
    expect(documentation).toHaveFocus();
  });

  test("does not announce an empty tree when there are no destinations", () => {
    render(
      <DragTree
        id="empty"
        label="Documentation"
        nodes={[]}
        getId={getId}
        getLabel={getLabel}
        getHref={getHref}
        getChildren={getChildren}
      />,
    );
    expect(screen.queryByRole("tree")).not.toBeInTheDocument();
  });

  test("Tab reaches only the focused row's actions and Escape returns to the node", async () => {
    const user = userEvent.setup();
    render(
      <>
        <TreeFixture editable />
        <button type="button">After navigation</button>
      </>,
    );
    await user.tab();
    const start = screen.getByRole("treeitem", { name: "Getting started" });
    expect(start).toHaveFocus();
    await user.tab();
    const action = screen.getByRole("button", {
      name: "Manage Getting started",
    });
    expect(action).toHaveFocus();
    await user.keyboard("{ArrowUp}");
    expect(action).toHaveFocus();
    await user.keyboard("{Escape}");
    expect(start).toHaveFocus();
    await user.keyboard("{ArrowDown}");
    expect(screen.getByRole("treeitem", { name: "About" })).toHaveFocus();
    await user.tab();
    expect(screen.getByRole("button", { name: "Manage About" })).toHaveFocus();
    await user.tab();
    expect(
      screen.getByRole("button", { name: "After navigation" }),
    ).toHaveFocus();
    await user.tab({ shift: true });
    expect(screen.getByRole("button", { name: "Manage About" })).toHaveFocus();
  });

  test("Enter activates the native branch link without collapsing its children", async () => {
    const user = userEvent.setup();
    renderTree();
    const link = screen.getByRole("link", { name: "Documentation" });
    const navigate = vi.fn((event: Event) => event.preventDefault());
    link.addEventListener("click", navigate);
    link.focus();
    await user.keyboard("{Enter}");
    expect(navigate).toHaveBeenCalledOnce();
    expect(link).toHaveAttribute("aria-expanded", "true");
    expect(
      screen.getByRole("treeitem", { name: "Getting started" }),
    ).toBeVisible();
  });

  test("Escape closes a portalled row menu before returning from its trigger to the tree", async () => {
    const user = userEvent.setup();
    renderTree({ menuActions: true });
    const node = screen.getByRole("treeitem", { name: "Getting started" });
    await user.tab();
    expect(node).toHaveFocus();
    await user.tab();
    const trigger = screen.getByRole("button", {
      name: "Manage Getting started",
    });
    expect(trigger).toHaveFocus();
    await user.keyboard("{Enter}");
    await waitFor(() => expect(screen.getByRole("menu")).toBeVisible());
    await user.keyboard("{Escape}");
    await waitFor(() =>
      expect(screen.queryByRole("menu")).not.toBeInTheDocument(),
    );
    await waitFor(() => expect(trigger).toHaveFocus());
    await user.keyboard("{Escape}");
    expect(node).toHaveFocus();
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
import { Portal } from "@ark-ui/react";
