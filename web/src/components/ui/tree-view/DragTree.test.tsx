import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, test, vi } from "vitest";

import * as Menu from "../menu";

import { DragTree, type DragTreeMove, isDragTreeDescendant } from ".";

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

type FixtureOptions = {
  treeNodes?: Node[];
  editable?: boolean;
  menuActions?: boolean;
  onMove?: (move: DragTreeMove<Node>) => Promise<void>;
  onSettled?: () => Promise<void>;
};

function renderTree(options?: FixtureOptions) {
  return render(<TreeFixture {...options} />);
}

function TreeFixture(options: FixtureOptions) {
  return (
    <DragTree
      id="documentation"
      label="Documentation"
      nodes={options.treeNodes ?? nodes}
      getId={getId}
      getLabel={getLabel}
      getHref={getHref}
      getChildren={getChildren}
      defaultExpandedIds={["root", "guides"]}
      move={
        options.onMove
          ? { onMove: options.onMove, onSettled: options.onSettled }
          : undefined
      }
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
      renderActions={({ node, tabIndex, startMove }) =>
        options.menuActions ? (
          <Menu.Root
            onSelect={({ value }) => {
              if (value === "move") startMove?.();
            }}
          >
            <Menu.Trigger asChild>
              <button type="button" tabIndex={tabIndex}>
                Manage {node.label}
              </button>
            </Menu.Trigger>
            <Portal>
              <Menu.Positioner>
                <Menu.Content>
                  <Menu.Item value="inspect">Inspect</Menu.Item>
                  {startMove && (
                    <Menu.Item value="move">Move or reorder…</Menu.Item>
                  )}
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
  test.each([
    { spaces: 0, placement: "inside", parent: "root" },
    { spaces: 1, placement: "below", parent: null },
    { spaces: 2, placement: "above", parent: null },
  ])(
    "keyboard placement $placement preserves the source and uses the target parent",
    async ({ spaces, placement, parent }) => {
      const onMove = vi.fn(async () => {});
      const onSettled = vi.fn(async () => {});
      const user = userEvent.setup();
      renderTree({ editable: true, onMove, onSettled });
      const source = screen.getByRole("treeitem", { name: "About" });
      source.focus();
      fireEvent.keyDown(source, {
        key: "M",
        code: "KeyM",
        altKey: true,
        shiftKey: true,
      });
      await waitFor(() =>
        expect(screen.getByRole("status")).toHaveTextContent("Carrying About"),
      );
      await user.keyboard("{Home}");
      for (let n = 0; n < spaces; n++) await user.keyboard(" ");
      expect(onMove).not.toHaveBeenCalled();
      expect(source).toHaveAttribute("href", "/docs/about");
      await user.keyboard("{Enter}");
      await waitFor(() => expect(onMove).toHaveBeenCalledTimes(1));
      expect(onMove).toHaveBeenCalledWith({
        node: nodes[1],
        oldParent: null,
        target: nodes[0],
        parent: parent ? nodes[0] : null,
        placement,
      });
      await waitFor(() =>
        expect(screen.getByRole("status")).toHaveTextContent("Moved About"),
      );
      expect(onSettled).toHaveBeenCalledTimes(1);
      await waitFor(() => expect(source).toHaveFocus());
    },
  );

  test("a carried branch cannot be placed in itself or any descendant, and Escape restores focus", async () => {
    const onMove = vi.fn(async () => {});
    const user = userEvent.setup();
    renderTree({ editable: true, onMove });
    const source = screen.getByRole("link", { name: "Guides" });
    source.focus();
    fireEvent.keyDown(source, {
      key: "M",
      code: "KeyM",
      altKey: true,
      shiftKey: true,
    });
    await user.keyboard("{Enter}{ArrowRight}{Enter} {Enter}{Escape}");
    expect(onMove).not.toHaveBeenCalled();
    await waitFor(() => expect(source).toHaveFocus());
    expect(screen.getByRole("status")).toHaveTextContent("Move cancelled");
    expect(
      screen.getByRole("treeitem", { name: "Getting started" }),
    ).toBeVisible();
  });

  test("a failed move is announced and can be retried without losing the carried node", async () => {
    const onMove = vi
      .fn()
      .mockRejectedValueOnce(new Error("Offline"))
      .mockResolvedValueOnce(undefined);
    const onSettled = vi.fn(async () => {});
    const user = userEvent.setup();
    renderTree({ editable: true, onMove, onSettled });
    const source = screen.getByRole("treeitem", { name: "About" });
    source.focus();
    fireEvent.keyDown(source, {
      key: "M",
      code: "KeyM",
      altKey: true,
      shiftKey: true,
    });
    await user.keyboard("{Home}{Enter}");
    await waitFor(() =>
      expect(screen.getByRole("alert")).toHaveTextContent("Could not move"),
    );
    expect(screen.getByRole("status")).toHaveTextContent("Carrying About");
    expect(screen.getByRole("link", { name: "Documentation" })).toHaveFocus();
    await user.keyboard("{Enter}");
    await waitFor(() =>
      expect(screen.getByRole("status")).toHaveTextContent("Moved About"),
    );
    expect(onMove).toHaveBeenCalledTimes(2);
    expect(onSettled).toHaveBeenCalledTimes(2);
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  test("pending placement ignores repeated confirmation and cancellation", async () => {
    let resolveMove!: () => void;
    const onMove = vi.fn(
      () =>
        new Promise<void>((resolve) => {
          resolveMove = resolve;
        }),
    );
    const user = userEvent.setup();
    renderTree({ editable: true, onMove });
    const source = screen.getByRole("treeitem", { name: "About" });
    source.focus();
    fireEvent.keyDown(source, {
      key: "M",
      code: "KeyM",
      altKey: true,
      shiftKey: true,
    });
    await user.keyboard("{Home}{Enter}{Enter}{Escape}");
    expect(onMove).toHaveBeenCalledTimes(1);
    expect(screen.getByRole("status")).toHaveTextContent("Placing the node");
    resolveMove();
    await waitFor(() =>
      expect(screen.getByRole("status")).toHaveTextContent("Moved About"),
    );
  });

  test("focus falls back to a visible destination when it hides the moved child", async () => {
    const onMove = vi.fn(async () => {
      view.rerender(
        <TreeFixture
          editable
          onMove={onMove}
          treeNodes={[{ id: "root", label: "Documentation", children: [] }]}
        />,
      );
    });
    const view = renderTree({ editable: true, onMove });
    const user = userEvent.setup();
    const source = screen.getByRole("treeitem", { name: "About" });
    source.focus();
    fireEvent.keyDown(source, {
      key: "M",
      code: "KeyM",
      altKey: true,
      shiftKey: true,
    });
    await user.keyboard("{Home}{Enter}");
    await waitFor(() =>
      expect(screen.getByRole("status")).toHaveTextContent(
        "Moved About inside Documentation",
      ),
    );
    await waitFor(() =>
      expect(
        screen.getByRole("treeitem", { name: "Documentation" }),
      ).toHaveFocus(),
    );
    expect(
      screen.queryByRole("treeitem", { name: "About" }),
    ).not.toBeInTheDocument();
  });

  test("Tab cancels a move without trapping keyboard navigation", async () => {
    const onMove = vi.fn(async () => {});
    const user = userEvent.setup();
    renderTree({ editable: true, onMove });
    const source = screen.getByRole("treeitem", { name: "About" });
    source.focus();
    fireEvent.keyDown(source, {
      key: "M",
      code: "KeyM",
      altKey: true,
      shiftKey: true,
    });
    await user.keyboard("{Home}{Tab}");
    expect(onMove).not.toHaveBeenCalled();
    expect(screen.getByRole("status")).toHaveTextContent("Move cancelled");
    expect(source).not.toHaveFocus();
  });

  test("the modified shortcut is unavailable to read-only users", async () => {
    const onMove = vi.fn(async () => {});
    renderTree({ onMove });
    const source = screen.getByRole("treeitem", { name: "About" });
    source.focus();
    fireEvent.keyDown(source, {
      key: "M",
      code: "KeyM",
      altKey: true,
      shiftKey: true,
    });
    expect(screen.getByRole("status")).toBeEmptyDOMElement();
    expect(source).not.toHaveAttribute("aria-keyshortcuts");
    expect(onMove).not.toHaveBeenCalled();
  });
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

  test("Shift+Tab out of the tree dismisses its keyboard shortcut hint", async () => {
    const user = userEvent.setup();
    render(
      <>
        <button type="button">Before navigation</button>
        <TreeFixture editable onMove={vi.fn(async () => {})} />
      </>,
    );
    const before = screen.getByRole("button", { name: "Before navigation" });
    await user.tab();
    await user.tab();
    await user.keyboard("{Shift}");
    expect(
      screen.getByRole("treeitem", { name: "Getting started" }),
    ).toHaveFocus();
    expect(await screen.findByRole("tooltip")).toBeVisible();

    await user.tab({ shift: true });
    expect(before).toHaveFocus();
    await waitFor(() =>
      expect(screen.queryByRole("tooltip")).not.toBeInTheDocument(),
    );

    // Leaving before the discovery delay must also cancel the pending hint.
    await user.tab();
    await user.keyboard("{Shift}");
    await user.tab({ shift: true });
    await new Promise((resolve) => setTimeout(resolve, 700));
    expect(before).toHaveFocus();
    expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
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
    renderTree({
      menuActions: true,
      editable: true,
      onMove: vi.fn(async () => {}),
    });
    const node = screen.getByRole("treeitem", { name: "Getting started" });
    await user.tab();
    expect(node).toHaveFocus();
    expect(await screen.findByRole("tooltip")).toHaveTextContent(
      "Move or reorder",
    );
    await user.tab();
    const trigger = screen.getByRole("button", {
      name: "Manage Getting started",
    });
    expect(trigger).toHaveFocus();
    await user.keyboard("{Enter}");
    await waitFor(() => expect(screen.getByRole("menu")).toBeVisible());
    await waitFor(() => expect(screen.getByRole("menu")).toHaveFocus());
    // Arrow navigation queues a focus task. Escape must invalidate it before
    // that frame runs, otherwise it steals focus back into the closing menu.
    const frames: FrameRequestCallback[] = [];
    const requestFrame = vi
      .spyOn(window, "requestAnimationFrame")
      .mockImplementation((callback) => frames.push(callback));
    fireEvent.keyDown(screen.getByRole("menu"), { key: "ArrowDown" });
    await user.keyboard("{Escape}");
    await waitFor(() => expect(trigger).toHaveFocus());
    requestFrame.mockRestore();
    expect(frames.length).toBeGreaterThan(0);
    act(() => frames.forEach((callback) => callback(performance.now())));
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
