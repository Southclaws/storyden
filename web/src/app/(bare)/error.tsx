"use client";

import { UnreadyBanner } from "@/components/site/Unready";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return <UnreadyBanner error={error} />;
}
