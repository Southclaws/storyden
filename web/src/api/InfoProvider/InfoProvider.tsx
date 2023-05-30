import { PropsWithChildren } from "react";
import ErrorBanner from "src/components/ErrorBanner";
import { useGetInfo } from "../openapi/misc";
import { InfoContext } from "./useInfoProvider";

export function InfoProvider({ children }: PropsWithChildren) {
  const { data, error } = useGetInfo();

  if (error) {
    // TODO: Handle outages in a PWA-friendly way twitter-style. swr should
    // cache locally quite well and allow read actions to work fine.
    // Also need to fix Next.js rewrite proxy error handling too...
    return <ErrorBanner {...error} />;
  }

  return <InfoContext.Provider value={data}>{children}</InfoContext.Provider>;
}
