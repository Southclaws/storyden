import { render, screen } from "@testing-library/react";
import { renderToReadableStream } from "react-dom/server";
import { expect, it, vi } from "vitest";

import { ThreadTagList } from "./ThreadTagList";

vi.mock("./ThreadTagListEditing", () => ({
  ThreadTagListEditing: () => <input aria-label="Edit tags" />,
}));

const tags = [{ name: "original", colour: "#123456", item_count: 1 }];

it("keeps tags readable in streamed HTML without exposing an unhydrated editor", async () => {
  const stream = await renderToReadableStream(
    <ThreadTagList editing initialTags={tags} onChange={vi.fn()} />,
  );
  await stream.allReady;
  const html = await new Response(stream).text();

  expect(html).toContain("original");
  expect(html).not.toContain("<input");
});

it("loads the tag editor after hydration when editing is requested", async () => {
  render(<ThreadTagList editing initialTags={tags} onChange={vi.fn()} />);

  expect(
    await screen.findByRole("textbox", { name: "Edit tags" }),
  ).toBeVisible();
});
