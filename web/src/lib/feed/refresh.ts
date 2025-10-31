"use server";

import { revalidatePath, revalidateTag } from "next/cache";

export async function refreshFeed() {
  revalidatePath("/", "layout");
  revalidateTag("api", "max");
}
