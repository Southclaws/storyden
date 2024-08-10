"use client";

import {
  ReadonlyURLSearchParams,
  useParams,
  useSearchParams,
} from "next/navigation";

import { ThreadListParams } from "src/api/openapi-schema";

import { FilterParamsSchema } from "./filterParams";

export function useFeedParams(): ThreadListParams {
  const sp = useSearchParams();
  const pp = useParams();

  const all = mergeParams(sp, pp);

  return parseThreadListParams(all);
}

export function parseThreadListParams(params: any): ThreadListParams {
  const parsed = FilterParamsSchema.safeParse(params);

  if (!parsed.success) {
    return {};
  }

  return {
    categories: parsed.data.category ? [parsed.data.category] : undefined,
  };
}

interface Params {
  [key: string]: string | string[];
}

function mergeParams(sp: ReadonlyURLSearchParams, pp: Params): Params {
  const search = Object.fromEntries(sp.entries());
  return {
    ...search,
    ...pp,
  };
}
