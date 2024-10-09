import { ProfileReference } from "src/api/openapi-schema";
import { Timestamp } from "src/components/site/Timestamp";

import { HStack } from "@/styled-system/jsx";

import { MemberBadge } from "../member/MemberBadge/MemberBadge";
import { DotSeparator } from "../site/Dot";

type Props = {
  href: string;
  author: ProfileReference;
  time: Date;
  updated: Date;
  more?: React.ReactElement;
};

export function Byline(props: Props) {
  return (
    <HStack alignItems="start" minWidth="0" w="full">
      <HStack
        alignItems="center"
        gap="0"
        minWidth="0"
        fontSize="sm"
        color="fg.subtle"
      >
        <MemberBadge profile={props.author} size="sm" name="handle" />
        <DotSeparator />
        <Timestamp created={props.time} href={props.href} />
      </HStack>

      {props.more}
    </HStack>
  );
}
