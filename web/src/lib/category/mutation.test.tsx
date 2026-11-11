import { act, renderHook, waitFor } from "@testing-library/react";
import type { PropsWithChildren } from "react";
import useSWR, { SWRConfig } from "swr";
import { expect, it, vi } from "vitest";

import {
  categoryUpdate,
  getCategoryListKey,
} from "@/api/openapi-client/categories";
import type { CategoryListOKResponse } from "@/api/openapi-schema";

import { useCategoryMutations } from "./mutation";

vi.mock("@/api/openapi-client/categories", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@/api/openapi-client/categories")>()),
  categoryUpdate: vi.fn(),
}));
vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

const initial: CategoryListOKResponse = {
  categories: [
    {
      id: "general",
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
      slug: "general",
      name: "General",
      description: "Community conversation",
      colour: "#8577ce",
      children: [],
      postCount: 0,
      sort: 0,
      cover_image: {
        id: "cover",
        path: "/api/assets/cover",
        filename: "cover.png",
        mime_type: "image/png",
        width: 200,
        height: 50,
      },
    },
  ],
};

it.each([null, undefined])(
  "updates the displayed cover when cover_image_asset_id is %s",
  async (coverID) => {
    const cache = new Map();
    function Wrapper({ children }: PropsWithChildren) {
      return (
        <SWRConfig value={{ provider: () => cache }}>{children}</SWRConfig>
      );
    }

    const { result } = renderHook(
      () => {
        const { data } = useSWR(getCategoryListKey(), async () => initial);
        const { updateCategory } = useCategoryMutations();
        return { data, updateCategory };
      },
      { wrapper: Wrapper },
    );

    await waitFor(() => expect(result.current.data).toBeDefined());
    await act(() =>
      result.current.updateCategory("general", {
        name: "Community",
        cover_image_asset_id: coverID,
      }),
    );

    expect(categoryUpdate).toHaveBeenCalledWith("general", {
      name: "Community",
      cover_image_asset_id: coverID,
    });
    const category = result.current.data?.categories[0];
    expect(category?.name).toBe("Community");
    if (coverID === null) {
      expect(category?.cover_image).toBeUndefined();
    } else {
      expect(category?.cover_image?.id).toBe("cover");
    }
  },
);
