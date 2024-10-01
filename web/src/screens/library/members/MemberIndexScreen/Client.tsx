"use client";

import { MemberIndexView } from "src/components/library/members/MemberIndexView/MemberIndexView";
import { Unready } from "src/components/site/Unready";

import { Props, useMemberIndexScreen } from "./useMemberIndexScreen";

export function Client(props: Props) {
  const { ready, data, mutate, error } = useMemberIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <MemberIndexView
      profiles={data}
      mutate={mutate}
      query={props.query}
      page={props.page}
    />
  );
}
