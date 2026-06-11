import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  NodeVersion,
  NodeWithChildren,
  Permission,
  PropertyType,
} from "@/api/openapi-schema";

import { LibraryPageProvider, useLibraryPageContext } from "./Context";
import { LibraryPageEditProvider, useEditState } from "./useEditState";
import { LibraryPageBlocks } from "./blocks/LibraryPageBlocks";
import { NodeStoreAPI } from "./store";

const mocks = vi.hoisted(() => ({
  contentMounts: 0,
  headingMounts: 0,
  nodeVersionCreate: vi.fn(),
  revalidate: vi.fn(),
  session: undefined as any,
}));

vi.mock("nuqs", async () => {
  const React = await vi.importActual<typeof import("react")>("react");

  return {
    parseAsBoolean: {},
    parseAsString: {},
    useQueryState: (_key: string, options?: { defaultValue?: unknown }) =>
      React.useState(options?.defaultValue ?? null),
  };
});

vi.mock("@/auth", () => ({
  useSession: () => mocks.session,
}));

vi.mock("@/api/openapi-client/nodes", () => ({
  nodeVersionCreate: mocks.nodeVersionCreate,
}));

vi.mock("@/lib/library/library", () => ({
  useLibraryMutation: () => ({
    revalidate: mocks.revalidate,
  }),
}));

vi.mock("@/lib/settings/capabilities", () => ({
  useCapability: () => false,
}));

vi.mock("@/components/content/ContentComposer/ContentComposer", async () => {
  const React = await vi.importActual<typeof import("react")>("react");

  return {
    ContentComposer: (props: { initialValue?: string }) => {
      const [mounted] = React.useState(() => {
        mocks.contentMounts += 1;
        return {
          content: props.initialValue ?? "",
          mount: mocks.contentMounts,
        };
      });

      return (
        <div data-mount={mounted.mount} data-testid="content-composer">
          {mounted.content}
        </div>
      );
    },
  };
});

vi.mock("@/components/ui/heading-input", async () => {
  const React = await vi.importActual<typeof import("react")>("react");

  return {
    HeadingInput: (props: { defaultValue?: string; value?: string }) => {
      const [mounted] = React.useState(() => {
        mocks.headingMounts += 1;
        return {
          mount: mocks.headingMounts,
          title: props.defaultValue ?? props.value ?? "",
        };
      });

      return (
        <span
          data-default-value={props.defaultValue}
          data-mount={mounted.mount}
          data-testid="heading-input"
          data-value={props.value}
        >
          {mounted.title}
        </span>
      );
    },
  };
});

describe("LibraryPageScreen editor source identity", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mocks.contentMounts = 0;
    mocks.headingMounts = 0;
    mocks.revalidate.mockResolvedValue(undefined);
    mocks.session = {
      id: "manager",
      roles: [
        {
          id: "manager-role",
          permissions: [Permission.MANAGE_LIBRARY],
        },
      ],
    };
  });

  it("remounts uncontrolled title and content editors when editing a draft checkpoint", async () => {
    const live = node("root", {
      content: "<p>Published content</p>",
      name: "Published title",
      current_version_id: "version-live",
    });
    mocks.nodeVersionCreate.mockResolvedValue(
      version("version-draft", {
        content: "<p>Draft content</p>",
        name: "Draft title",
      }),
    );

    render(<Harness node={live} />);

    expect(screen.getByTestId("content-composer")).toHaveTextContent(
      "Published content",
    );

    await act(async () => {
      fireEvent.click(screen.getByRole("button", { name: "edit draft" }));
    });

    await waitFor(() => {
      expect(screen.getByTestId("heading-input")).toHaveTextContent(
        "Draft title",
      );
    });
    expect(screen.getByTestId("content-composer")).toHaveTextContent(
      "Draft content",
    );
    expect(screen.getByTestId("content-composer")).toHaveAttribute(
      "data-mount",
      "2",
    );
    expect(screen.getByTestId("heading-input")).toHaveAttribute(
      "data-mount",
      "1",
    );

    await act(async () => {
      fireEvent.click(screen.getByRole("button", { name: "view page" }));
    });

    expect(screen.queryByTestId("heading-input")).not.toBeInTheDocument();
    expect(screen.getByText("Published title")).toBeInTheDocument();
    expect(screen.getByTestId("content-composer")).toHaveTextContent(
      "Published content",
    );
    expect(screen.getByTestId("content-composer")).toHaveAttribute(
      "data-mount",
      "3",
    );
  });

  it("does not remount uncontrolled editors for ordinary draft field updates", async () => {
    const live = node("root", {
      content: "<p>Published content</p>",
      name: "Published title",
    });
    mocks.nodeVersionCreate.mockResolvedValue(
      version("version-draft", {
        content: "<p>Draft content</p>",
        name: "Draft title",
      }),
    );

    let store: NodeStoreAPI | undefined;
    render(<Harness node={live} onStore={(s) => (store = s)} />);

    await act(async () => {
      fireEvent.click(screen.getByRole("button", { name: "edit draft" }));
    });

    await waitFor(() => {
      expect(screen.getByTestId("content-composer")).toHaveTextContent(
        "Draft content",
      );
    });

    act(() => {
      store?.getState().setContent("<p>Typed draft content</p>");
      store?.getState().setName("Typed draft title");
    });

    expect(screen.getByTestId("content-composer")).toHaveTextContent(
      "Draft content",
    );
    expect(screen.getByTestId("content-composer")).toHaveAttribute(
      "data-mount",
      "2",
    );
    expect(screen.getByTestId("heading-input")).toHaveTextContent(
      "Draft title",
    );
    expect(screen.getByTestId("heading-input")).toHaveAttribute(
      "data-default-value",
      "Draft title",
    );
    expect(screen.getByTestId("heading-input")).not.toHaveAttribute(
      "data-value",
    );
    expect(screen.getByTestId("heading-input")).toHaveAttribute(
      "data-mount",
      "1",
    );
  });

  it("updates the view title when the live draft state changes", () => {
    const live = node("root", {
      content: "<p>Published content</p>",
      name: "Published title",
    });

    let store: NodeStoreAPI | undefined;
    render(<Harness node={live} onStore={(s) => (store = s)} />);

    expect(screen.getByText("Published title")).toBeInTheDocument();

    act(() => {
      store?.getState().setName("Applied title");
      store?.getState().setContent("<p>Applied content</p>");
    });

    expect(screen.getByText("Applied title")).toBeInTheDocument();
    expect(screen.getByTestId("content-composer")).toHaveTextContent(
      "Published content",
    );
  });
});

function Harness({
  node,
  onStore,
}: {
  node: NodeWithChildren;
  onStore?: (store: NodeStoreAPI) => void;
}) {
  return (
    <LibraryPageProvider node={node}>
      <LibraryPageEditProvider>
        <EditControls />
        {onStore && <StoreProbe onStore={onStore} />}
        <LibraryPageBlocks />
      </LibraryPageEditProvider>
    </LibraryPageProvider>
  );
}

function EditControls() {
  const { startProposalEdit, stopEditing } = useEditState();

  return (
    <>
      <button type="button" onClick={() => void startProposalEdit()}>
        edit draft
      </button>
      <button type="button" onClick={stopEditing}>
        view page
      </button>
    </>
  );
}

function StoreProbe({ onStore }: { onStore: (store: NodeStoreAPI) => void }) {
  const { store } = useLibraryPageContext();

  onStore(store);

  return null;
}

function node(
  id: string,
  overrides: Partial<NodeWithChildren> = {},
): NodeWithChildren {
  return {
    id,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    name: `Node ${id}`,
    slug: `node-${id}`,
    description: "description",
    content: "<p>hello</p>",
    owner: profile("owner"),
    assets: [],
    tags: [],
    visibility: "published",
    hide_child_tree: false,
    meta: {
      layout: {
        blocks: [{ type: "title" }, { type: "content" }],
      },
    },
    properties: [
      {
        fid: "field-1",
        name: "Label",
        sort: "0",
        type: PropertyType.text,
        value: "old value",
      },
    ],
    child_property_schema: [],
    children: [],
    recomentations: [],
    ...overrides,
  };
}

function version(
  id: string,
  overrides: Partial<NodeVersion> = {},
): NodeVersion {
  return {
    id,
    node_id: "root",
    author: profile("author"),
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
    status: "draft",
    name: "Draft title",
    slug: "draft-title",
    description: "Draft description",
    content: "<p>Draft content</p>",
    meta: {},
    properties: [],
    ...overrides,
  };
}

function profile(id: string) {
  return {
    id,
    handle: id,
    name: `Member ${id}`,
    joined: "2026-01-01T00:00:00Z",
    roles: [],
  };
}
