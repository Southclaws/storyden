import { RequestError } from "@/api/common";

const ErrUnexpected = "An unexpected error occurred";

/**
 * Derives an end-user message from an error/exception value.
 * @param e exception or any error type
 */
export function deriveError(e: unknown): string {
  if (e === null || e === undefined) {
    return "";
  }

  if (typeof e === "string") {
    return e;
  }

  if (e instanceof Error) {
    if (e instanceof RequestError) {
      return e.message;
    }

    if (e instanceof TypeError) {
      console.error(e);
      return ErrUnexpected;
    }

    if (e.message.includes("React")) {
      // React prints these by default.
      return "Something went wrong while rendering data.";
    }

    return e.message ?? ErrUnexpected;
  }

  console.error("unable to derive error text:", e);
  return "An unknown error occurred";
}
