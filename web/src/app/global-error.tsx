"use client";

import { useEffect } from "react";

import { GenericError } from "@/screens/errors/GenericError";
import { deriveError } from "@/utils/error";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error("Global error caught:", error);
  }, [error]);

  return (
    <html>
      <body>
        <GenericError reset={reset} message={deriveError(error)} />
      </body>
    </html>
  );
}