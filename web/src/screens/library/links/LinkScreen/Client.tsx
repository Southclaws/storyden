"use client";

import { LinkView } from "src/components/library/links/LinkView";
import { Unready } from "src/components/site/Unready";

import { Props, useLinkScreen } from "./useLinkScreen";

export function Client(props: Props) {
  const { ready, data, error } = useLinkScreen(props);

  if (!ready) return <Unready {...error} />;

  return <LinkView link={data} />;
}
