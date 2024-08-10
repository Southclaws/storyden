import { Options, buildRequest, buildResult } from "./common";

export const fetcher = async <T>(opts: Options): Promise<T> => {
  const response = await fetch(buildRequest(opts));

  return buildResult<T>(response);
};
