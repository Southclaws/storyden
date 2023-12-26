"use client";

import { Unready } from "src/components/site/Unready";

import { LinkIndexView } from "./components/LinkIndexView/LinkIndexView";
import { Props, useLinkIndexScreen } from "./useLinkIndexScreen";

export function Client(props: Props) {
  const { ready, data, mutate, error } = useLinkIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <LinkIndexView
      links={data}
      mutate={mutate}
      query={props.query}
      page={props.page}
    />
  );
}
