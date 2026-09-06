import { afterEach, expect, test, vi } from "vitest";

import { reloadThemeStylesheet } from "./reload-theme-stylesheet";

afterEach(() => {
  document.head.replaceChildren();
  vi.restoreAllMocks();
});

test("keeps the current theme until the refreshed stylesheet loads", async () => {
  vi.spyOn(Date, "now").mockReturnValue(1234);
  const current = document.createElement("link");
  current.rel = "stylesheet";
  current.href = "/theme.css";
  current.dataset["sdThemeAsset"] = "stylesheet";
  const staleLiveStyle = document.createElement("style");
  staleLiveStyle.dataset["sdThemeEditorLive"] = "";
  document.head.append(current, staleLiveStyle);

  const reloaded = reloadThemeStylesheet();
  const links = document.querySelectorAll<HTMLLinkElement>(
    'link[data-sd-theme-asset="stylesheet"]',
  );
  const replacement = links[1]!;

  expect(links).toHaveLength(2);
  expect(current).toBeInTheDocument();
  expect(replacement.href).toBe(
    "http://localhost:3000/theme.css?sd-theme-editor=1234",
  );

  replacement.dispatchEvent(new Event("load"));
  await reloaded;

  expect(current).not.toBeInTheDocument();
  expect(replacement).toBeInTheDocument();
  expect(staleLiveStyle).not.toBeInTheDocument();
});

test("keeps the current theme when the refreshed stylesheet fails", async () => {
  const current = document.createElement("link");
  current.rel = "stylesheet";
  current.href = "/theme.css";
  current.dataset["sdThemeAsset"] = "stylesheet";
  document.head.append(current);

  const reloaded = reloadThemeStylesheet();
  const replacement = document.querySelectorAll<HTMLLinkElement>(
    'link[data-sd-theme-asset="stylesheet"]',
  )[1]!;
  replacement.dispatchEvent(new Event("error"));

  await expect(reloaded).rejects.toThrow(
    "The updated theme stylesheet could not be loaded.",
  );
  expect(current).toBeInTheDocument();
  expect(replacement).not.toBeInTheDocument();
});
