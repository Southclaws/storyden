export const LibraryPageEditMode = {
  direct: "direct",
  proposal: "proposal",
} as const;

export type LibraryPageEditMode =
  (typeof LibraryPageEditMode)[keyof typeof LibraryPageEditMode];

export function normaliseLibraryPageEditMode(
  value: string | null | undefined,
): LibraryPageEditMode {
  if (value === LibraryPageEditMode.proposal) {
    return LibraryPageEditMode.proposal;
  }

  return LibraryPageEditMode.direct;
}
