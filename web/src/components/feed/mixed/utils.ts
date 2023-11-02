import { ThreadReference } from "src/api/openapi/schemas";

export function getPostType(thread: ThreadReference) {
  if (thread.link) {
    return "link";
  }

  return "text";
}
