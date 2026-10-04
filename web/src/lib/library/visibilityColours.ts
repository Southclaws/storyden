import { Visibility } from "@/api/openapi-schema";
import type { SystemProperties } from "@/styled-system/types";

export function visibilityColour(
  v: Visibility,
): SystemProperties["colorPalette"] {
  switch (v) {
    case Visibility.published:
      return "visibility.published";
    case Visibility.review:
      return "visibility.review";
    case Visibility.draft:
      return "visibility.draft";
    case Visibility.unlisted:
      return "visibility.unlisted";
  }
}
