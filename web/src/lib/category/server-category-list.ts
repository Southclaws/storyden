"use server";

import { cache } from "react";

import { categoryList } from "@/api/openapi-server/categories";

export const categoryListCached = cache(async () => {
  return await categoryList({
    cache: "default",
  });
});
