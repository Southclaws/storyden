"use client";

import { MemberIndexView } from "src/components/directory/members/MemberIndexView/MemberIndexView";
import { Unready } from "src/components/site/Unready";

import { Props, useMemberIndexScreen } from "./useMemberIndexScreen";

export function Client(props: Props) {
  const { ready, data, error } = useMemberIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  return <MemberIndexView profiles={data} query={props.query} />;
}
