import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { Card } from ".";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ prefetch: vi.fn() }),
}));

describe("Card accessibility", () => {
  it.each(["", "<body></body>", "<body><p></p></body>", "<p></p>"])(
    "omits previews for empty rich content (%s)",
    (content) => {
      const { container } = render(
        <Card id="empty" title="Empty page" url="/l/empty" content={content} />,
      );
      expect(
        container.querySelector(".rich-card__contentContainer"),
      ).toBeNull();
      expect(screen.queryByRole("button", { name: /Show more/ })).toBeNull();
    },
  );
  it.each(
    [undefined, null, false, 0, "", " ", [], <></>, [false, <></>]].map(
      (empty) => ({ empty }),
    ),
  )("omits empty slots ($empty)", ({ empty }) => {
    const { container } = render(
      <Card
        id="sparse"
        title="Title only"
        url="/l/sparse"
        header={empty}
        menu={empty}
        body={empty}
        controls={empty}
      >
        {empty}
      </Card>,
    );

    for (const slot of [
      "headerContainer",
      "menuContainer",
      "contentContainer",
      "textArea",
      "footerContainer",
    ]) {
      expect(container.querySelector(`.rich-card__${slot}`)).toBeNull();
    }
    expect(screen.getByRole("link", { name: "Title only" })).toHaveAttribute(
      "href",
      "/l/sparse",
    );
  });

  it("renders structured body and footer controls without an empty preview", () => {
    const { container, rerender } = render(
      <Card
        id="details"
        title="Field notes"
        url="/l/notes"
        body={<p>Status: Open</p>}
        controls={<button>Bookmark</button>}
      />,
    );
    expect(
      screen.getByText("Status: Open").closest(".rich-card__contentContainer"),
    ).not.toBeNull();
    expect(
      screen
        .getByRole("button", { name: "Bookmark" })
        .closest(".rich-card__footerContainer"),
    ).not.toBeNull();
    expect(container.querySelector(".rich-card__textArea")).toBeNull();

    rerender(<Card id="details" title="Field notes" url="/l/notes" />);
    expect(container.querySelector(".rich-card__contentContainer")).toBeNull();
    expect(container.querySelector(".rich-card__footerContainer")).toBeNull();
  });

  it("offers one named destination while keeping the excerpt readable", () => {
    render(
      <Card
        id="harbour"
        title="Harbour notes"
        text="Boats are arriving."
        url="/t/harbour"
        titleStatus="Pinned discussion"
      />,
    );

    expect(
      screen.getByRole("article", {
        name: "Harbour notes Pinned discussion",
      }),
    ).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2 })).toHaveTextContent(
      "Harbour notes",
    );
    expect(screen.getAllByRole("link")).toHaveLength(1);
    expect(screen.getByRole("link", { name: "Harbour notes" })).toHaveAttribute(
      "href",
      "/t/harbour",
    );
    expect(screen.getByText("Boats are arriving.").closest("a")).toBeNull();
  });

  it("announces a meaningful foreground image once", () => {
    const { container } = render(
      <Card
        id="photo"
        title="Harbour photo"
        url="/t/photo"
        image="/photo.png"
        imageAlt="A red boat beside the wooden pier"
        headingAs="h3"
      />,
    );
    expect(screen.getByRole("heading", { level: 3 })).toBeInTheDocument();
    expect(screen.getAllByRole("img")).toHaveLength(1);
    expect(screen.getByRole("img")).toHaveAccessibleName(
      "A red boat beside the wooden pier",
    );
    expect(container.querySelector('img[aria-hidden="true"]')).toHaveAttribute(
      "alt",
      "",
    );
  });

  it("keeps inline editors named and structured details outside the navigation link", async () => {
    const user = userEvent.setup();
    render(
      <Card
        id="edit"
        title="Field notes"
        url="/l/notes"
        titleContent={
          <input aria-label="Page name" defaultValue="Field notes" />
        }
        body={<p>Status: Open for contributions</p>}
        disableAnchors
      />,
    );
    expect(
      screen.getByRole("article", { name: "Field notes" }),
    ).toBeInTheDocument();
    await user.click(screen.getByRole("textbox", { name: "Page name" }));
    await user.keyboard(" updated");
    expect(screen.getByRole("textbox")).toHaveValue("Field notes updated");
    expect(
      screen.getByText("Status: Open for contributions").closest("a"),
    ).toBeNull();
  });

  it("names image-only cards and their destination when the visible title is hidden", () => {
    render(
      <Card
        id="image-only"
        label="Field notes"
        url="/l/notes"
        image="/photo.png"
      />,
    );
    expect(
      screen.getByRole("article", { name: "Field notes" }),
    ).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Field notes" })).toHaveAttribute(
      "href",
      "/l/notes",
    );
  });

  it("keeps decorative teasers out of image navigation", () => {
    render(
      <Card id="teaser" title="Harbour" url="/t/harbour" image="/photo.png" />,
    );
    expect(screen.queryByRole("img")).not.toBeInTheDocument();
  });

  it("preserves a named content link when fill hides the title slot", () => {
    render(
      <Card
        id="fill"
        title="Alpine field notes"
        text="A quiet lake."
        url="/l/alpine"
        shape="fill"
      />,
    );
    expect(screen.getByText("A quiet lake.").closest("a")).toHaveAccessibleName(
      "Alpine field notes",
    );
    expect(screen.getByText("A quiet lake.").closest("a")).toHaveAttribute(
      "href",
      "/l/alpine",
    );
  });

  it("expands and collapses a long preview from the keyboard without losing focus", async () => {
    vi.spyOn(HTMLElement.prototype, "scrollHeight", "get").mockReturnValue(220);
    const user = userEvent.setup();
    render(
      <Card id="long" title="Field notes" text="Long preview" url="/t/notes" />,
    );

    const disclosure = screen.getByRole("button", {
      name: "Show more of Field notes",
      expanded: false,
    });
    expect(
      document.getElementById(disclosure.getAttribute("aria-controls")!),
    ).toHaveTextContent("Long preview");
    await user.tab();
    expect(screen.getByRole("link")).toHaveFocus();
    await user.tab();
    expect(disclosure).toHaveFocus();
    await user.keyboard("{Enter}");
    expect(disclosure).toHaveAttribute("aria-expanded", "true");
    expect(disclosure).toHaveAccessibleName("Show less of Field notes");
    expect(disclosure).toHaveFocus();
    await user.keyboard(" ");
    expect(disclosure).toHaveAttribute("aria-expanded", "false");
    expect(disclosure).toHaveFocus();
  });
});

describe("Minimal card media", () => {
  it.each(["minimal", "row", "responsive", "box", "fill"] as const)(
    "supports caller-provided fallback content in the %s media slot",
    (shape) => {
      const { container } = render(
        <Card
          id="processing"
          title="Field notes"
          url="/l/notes"
          shape={shape}
          imageFallback={<span>Cover is processing</span>}
        />,
      );
      expect(
        screen
          .getByText("Cover is processing")
          .closest(".rich-card__mediaContainer"),
      ).not.toBeNull();
      expect(container.querySelector("img")).toBeNull();
    },
  );

  it("renders a responsive linked cover without a backdrop", () => {
    const { container } = render(
      <Card
        id="minimal"
        title="Field notes"
        url="/l/notes"
        shape="minimal"
        image="/photo.png"
        imageAlt="Mountain lake"
        imageAspectRatio={4 / 3}
        imageProps={{
          srcSet: "/small.png 320w, /large.png 960w",
          sizes: "408px",
          width: 1200,
          height: 800,
          loading: "lazy",
          decoding: "async",
        }}
      />,
    );
    const image = screen.getByRole("img", { name: "Mountain lake" });
    expect(image).toHaveAttribute("srcset", "/small.png 320w, /large.png 960w");
    expect(image).toHaveAttribute("sizes", "408px");
    expect(image).toHaveAttribute("width", "1200");
    expect(image).toHaveAttribute("height", "800");
    expect(image).toHaveAttribute("loading", "lazy");
    expect(image.closest("a")).toHaveAccessibleName("Field notes");
    expect(image.parentElement).toHaveStyle({
      aspectRatio: "1.3333333333333333",
    });
    expect(
      container.querySelector(".rich-card__mediaBackdropContainer"),
    ).toBeNull();
    expect(
      container.querySelector(".rich-card__root--shape_minimal"),
    ).not.toBeNull();
  });

  it("keeps a fallback-only card named and linked", () => {
    const { container, rerender } = render(
      <Card
        id="fallback"
        label="Untitled notebook"
        url="/l/notebook"
        shape="minimal"
        imageFallback={<span>Untitled notebook</span>}
      />,
    );
    expect(
      screen.getByRole("article", { name: "Untitled notebook" }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("link", { name: "Untitled notebook" }),
    ).toHaveAttribute("href", "/l/notebook");
    expect(container.querySelector("img")).toBeNull();
    expect(
      container.querySelector(".rich-card__mediaMissing"),
    ).toHaveTextContent("Untitled notebook");

    rerender(
      <Card
        id="fallback"
        label="Untitled notebook"
        url="/l/notebook"
        shape="minimal"
        imageFallback={false}
      />,
    );
    expect(container.querySelector(".rich-card__mediaContainer")).toBeNull();
  });

  it("replaces the fallback when a cover is available", () => {
    const { container } = render(
      <Card
        id="cover"
        title="Notes"
        url="/l/notes"
        shape="minimal"
        image="/photo.png"
        imageFallback="Processing cover"
      />,
    );
    expect(container.querySelector("img")).toHaveAttribute("src", "/photo.png");
    expect(screen.queryByText("Processing cover")).toBeNull();
  });

  it("keeps media out of navigation while editing and leaves the title editable", async () => {
    const user = userEvent.setup();
    const { container } = render(
      <Card
        id="editing-minimal"
        title="Field notes"
        url="/l/notes"
        shape="minimal"
        image="/photo.png"
        titleContent={<input aria-label="Name" defaultValue="Field notes" />}
        disableAnchors
      />,
    );
    expect(container.querySelector("img")?.closest("a")).toBeNull();
    await user.type(screen.getByRole("textbox", { name: "Name" }), " updated");
    expect(screen.getByRole("textbox", { name: "Name" })).toHaveValue(
      "Field notes updated",
    );
  });
});
