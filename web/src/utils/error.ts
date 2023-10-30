import { ZodError } from "zod";

/**
 * Derives an end-user message from an error/exception value.
 * @param e exception or any error type
 */
export function deriveError(e: unknown): string {
  if (e instanceof Error) {
    return e.message;
  }

  console.error("unhandled error", e, typeof e);

  return "unknown error occurred";
}

export function zodFormError(error: ZodError) {
  return {
    message: error.issues.reduce((i, c) => [c.message, ...i], [""]).join(", "),
  };
}
