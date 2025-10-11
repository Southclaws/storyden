import { ThreadListParams } from "@/api/openapi-schema";

// Helper function to get thread list params based on category mode
export function getCategoryThreadListParams(
  mode: "none" | "all" | "uncategorised",
  page?: number,
): ThreadListParams {
  if (mode === "none") {
    return {};
  }

  const params: ThreadListParams = {
    page: page?.toString(),
  };

  if (mode === "uncategorised") {
    params.categories = ["null"];
  }

  return params;
}
