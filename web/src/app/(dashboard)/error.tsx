"use client";

import { useEffect } from "react";

import { GenericError } from "@/screens/errors/GenericError";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error("(dashboard) error caught:", error);
  }, [error]);

  return <GenericError reset={reset} />;
}
