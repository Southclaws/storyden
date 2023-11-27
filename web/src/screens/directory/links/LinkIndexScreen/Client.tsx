"use client";

import { Unready } from "src/components/site/Unready";

import { LinkResultList } from "./components/LinkResultList/LinkResultList";
import { Props, useLinkIndexScreen } from "./useLinkIndexScreen";

export function Client(props: Props) {
  const { ready, data, error } = useLinkIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  return <LinkResultList links={data} />;
}
