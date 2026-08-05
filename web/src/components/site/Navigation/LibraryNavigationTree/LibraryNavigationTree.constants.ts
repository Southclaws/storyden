import { Visibility } from "@/api/openapi-schema";

export const libraryNavigationVisibility = [
  Visibility.draft,
  Visibility.review,
  Visibility.unlisted,
  Visibility.published,
] satisfies Visibility[];
