import { useSession } from "@/auth";

import { EmptyState } from "../site/EmptyState";

export function LibraryEmptyState() {
  const session = useSession();

  const contributionLabel = session
    ? "Be the first to contribute!"
    : "Please log in to contribute.";

  return (
    <EmptyState>
      <p>This community library is empty.</p>
      <p>{contributionLabel}</p>
    </EmptyState>
  );
}
