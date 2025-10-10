// Helper function to get thread list params based on category mode
export function getCategoryThreadListParams(
  mode: "none" | "all" | "uncategorised",
  page?: number,
) {
  if (mode === "none") {
    return {};
  }

  const params: { page?: string; categories?: string[] } = {
    page: page?.toString(),
  };

  if (mode === "uncategorised") {
    params.categories = ["null"];
  }

  return params;
}
