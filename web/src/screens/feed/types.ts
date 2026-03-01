import {
  CategoryListOKResponse,
  NodeGetOKResponse,
  NodeListResult,
  ThreadListResult,
} from "@/api/openapi-schema";

export type InitialData = {
  initialPage?: number;
  initialThreadList?: ThreadListResult;
  initialLibraryNodeList?: NodeListResult;
  initialLibraryNode?: NodeGetOKResponse;
  initialCategoryList?: CategoryListOKResponse;
};
