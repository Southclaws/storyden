"use client";

import { Asset } from "@/api/openapi-schema";
import { buildRequest, buildResult } from "@/api/common";

export async function uploadThemeAsset(file: File): Promise<Asset> {
  const request = buildRequest({
    url: "/info/theme/assets",
    method: "POST",
    params: {
      filename: file.name,
    },
    headers: {
      "Content-Type": "application/octet-stream",
    },
    data: file,
  });

  const response = await fetch(request);

  return buildResult<Asset>(response);
}
