import { act, render, screen } from "@testing-library/react";
import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import { expect, it } from "vitest";

import { enableWebMCPTools, useWebMCPToolsEnabled } from "./activation";

function ToolsState() {
  const enabled = useWebMCPToolsEnabled();
  return <output>{enabled ? "enabled" : "disabled"}</output>;
}

it("keeps tools disabled during SSR and browsing, then retains activation across remounts", () => {
  expect(renderToString(<ToolsState />)).toContain("disabled");
  const view = render(
    <StrictMode>
      <ToolsState />
    </StrictMode>,
  );
  expect(screen.getByRole("status")).toHaveTextContent("disabled");

  act(() => enableWebMCPTools());
  expect(screen.getByRole("status")).toHaveTextContent("enabled");
  expect(renderToString(<ToolsState />)).toContain("disabled");
  view.unmount();

  render(<ToolsState />);
  act(() => enableWebMCPTools());
  expect(screen.getByRole("status")).toHaveTextContent("enabled");
});
