import { act, renderHook, waitFor } from "@testing-library/react";
import { type PropsWithChildren } from "react";
import { SWRConfig, useSWRConfig } from "swr";
import { expect, it, vi } from "vitest";

import { getCategoryListKey } from "@/api/openapi-client/categories";

import { useCategorySelect } from "./useCategorySelect";

function CategoriesCache({ children }: PropsWithChildren) {
  return (
    <SWRConfig value={{ provider: () => new Map(), errorRetryCount: 0 }}>
      {children}
    </SWRConfig>
  );
}

it("uses server categories immediately without fetching again and still supports revalidation", async () => {
  const fetch = vi
    .spyOn(globalThis, "fetch")
    .mockResolvedValue(Response.json({ categories: [] }));
  const { result } = renderHook(
    () => ({
      select: useCategorySelect({ categories: [] }),
      cache: useSWRConfig(),
    }),
    { wrapper: CategoriesCache },
  );

  expect(result.current.select.ready).toBe(true);
  expect(result.current.select.collection?.items).toEqual([
    { label: "No category", value: "__none__" },
  ]);
  await act(async () => {});
  expect(fetch).not.toHaveBeenCalled();

  await act(async () => {
    await result.current.cache.mutate(getCategoryListKey());
  });
  expect(fetch).toHaveBeenCalledTimes(1);
});

it("fetches categories when no server data is available", async () => {
  const fetch = vi
    .spyOn(globalThis, "fetch")
    .mockResolvedValue(Response.json({ categories: [] }));
  const { result } = renderHook(() => useCategorySelect(), {
    wrapper: CategoriesCache,
  });

  await waitFor(() => expect(result.current.ready).toBe(true));
  expect(fetch).toHaveBeenCalledTimes(1);
});
