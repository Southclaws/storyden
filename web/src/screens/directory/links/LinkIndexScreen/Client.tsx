"use client";

import { Unready } from "src/components/site/Unready";

import { LinkIndexView } from "./components/LinkIndexView/LinkIndexView";
import { Props, useLinkIndexScreen } from "./useLinkIndexScreen";

export function Client(props: Props) {
  const { ready, data, error } = useLinkIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  return <LinkIndexView links={data} query={props.query} />;
}
