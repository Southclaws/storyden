import { Visibility } from "@/api/openapi-schema";
import { ColorPalette } from "@/styled-system/tokens";

export function visibilityColour(v: Visibility): ColorPalette {
  switch (v) {
    case Visibility.published:
      return "bg";
    case Visibility.review:
      return "blue";
    case Visibility.draft:
      return "green";
    case Visibility.unlisted:
      return "pink";
  }
}
