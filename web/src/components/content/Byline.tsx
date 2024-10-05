import { ProfileReference } from "src/api/openapi-schema";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Timestamp } from "src/components/site/Timestamp";

import { HStack, styled } from "@/styled-system/jsx";

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
        <ProfilePill profileReference={props.author} />
        <styled.span pr="2">•</styled.span>
        <Timestamp
          created={props.time}
          updated={props.updated}
          href={props.href}
        />
      </HStack>

      {props.more}
    </HStack>
  );
}
