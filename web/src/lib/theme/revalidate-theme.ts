"use server";

import { revalidatePath, updateTag } from "next/cache";

import { Permission } from "@/api/openapi-schema";
import { getServerSession } from "@/auth/server-session";
import { hasPermission } from "@/utils/permissions";

import { THEME_CACHE_TAG } from "./theme-server";

export async function revalidateTheme() {
  const session = await getServerSession();
  if (!hasPermission(session, Permission.ADMINISTRATOR)) {
    throw new Error("Administrator permission is required.");
  }

  updateTag(THEME_CACHE_TAG);
  revalidatePath("/", "layout");
}
