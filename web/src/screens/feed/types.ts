import {
  CategoryListResult,
  NodeListResult,
  ThreadListResult,
} from "@/api/openapi-schema";

export type InitialData = {
  threads?: ThreadListResult;
  page?: number;
  library?: NodeListResult;
  categories?: CategoryListResult;
};
