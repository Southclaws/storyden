import {
  CategoryListOKResponse,
  NodeGetOKResponse,
  NodeListResult,
  ThreadListResult,
} from "@/api/openapi-schema";

export type InitialData = {
  initialPage?: number;
  initialThreadList?: ThreadListResult;
  initialThreadSource?: "all" | "uncategorised";
  initialLibraryNodeList?: NodeListResult;
  initialLibraryNode?: NodeGetOKResponse;
  initialCategoryList?: CategoryListOKResponse;
};
