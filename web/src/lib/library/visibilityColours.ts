import { Visibility } from "@/api/openapi-schema";
import { ColorPalette } from "@/styled-system/tokens";
import { UtilityValues } from "@/styled-system/types/prop-type";

export function visibilityColour(v: Visibility): UtilityValues["colorPalette"] {
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
