import { act, fireEvent, render, screen } from "@testing-library/react";
import type { ComponentType } from "react";
import { renderToString } from "react-dom/server";
import { afterEach, expect, it, vi } from "vitest";

import { setThemeEditingEnabled } from "@/lib/theme/theme-editing";

import { ThemeEditorHost } from "./ThemeEditorHost";

const { auth, loadController } = vi.hoisted(() => ({
  auth: {
    account: undefined as { roles: { permissions: string[] }[] } | undefined,
  },
  loadController: vi.fn(),
}));

vi.mock("@/api/openapi-client/accounts", () => ({
  useAccountGet: () => ({ data: auth.account }),
}));

vi.mock("next/dynamic", async () => {
  const { lazy, Suspense } = await import("react");
  return {
    default: (load: () => Promise<ComponentType>) => {
      const Controller = lazy(async () => {
        loadController();
        return { default: await load() };
      });
      return () => (
        <Suspense fallback={null}>
          <Controller />
        </Suspense>
      );
    },
  };
});

vi.mock("./ThemeEditorController", () => ({
  ThemeEditorController: () => (
    <button onClick={() => setThemeEditingEnabled(false)}>
      Exit theme editing
    </button>
  ),
}));

afterEach(() => localStorage.clear());

it("loads only for an administrator with editing enabled, and supports closing and reopening", async () => {
  const view = render(<ThemeEditorHost />);
  expect(loadController).not.toHaveBeenCalled();

  act(() => setThemeEditingEnabled(true));
  expect(loadController).not.toHaveBeenCalled();
  auth.account = { roles: [{ permissions: ["MANAGE_SETTINGS"] }] };
  view.rerender(<ThemeEditorHost />);
  expect(loadController).not.toHaveBeenCalled();

  auth.account = { roles: [{ permissions: ["ADMINISTRATOR"] }] };
  // The persisted browser flag never enables the editor during SSR.
  expect(renderToString(<ThemeEditorHost />)).toBe("");
  view.rerender(<ThemeEditorHost />);
  const close = await screen.findByRole("button", {
    name: "Exit theme editing",
  });
  expect(loadController).toHaveBeenCalledTimes(1);

  fireEvent.click(close);
  expect(screen.queryByRole("button")).not.toBeInTheDocument();
  act(() => setThemeEditingEnabled(true));
  expect(
    await screen.findByRole("button", { name: "Exit theme editing" }),
  ).toBeVisible();
  expect(loadController).toHaveBeenCalledTimes(1);

  auth.account = undefined;
  view.rerender(<ThemeEditorHost />);
  expect(screen.queryByRole("button")).not.toBeInTheDocument();
});
