import { ProfileReference } from "src/api/openapi-schema";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Timestamp } from "src/components/site/Timestamp";

import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  href: string;
  author: ProfileReference;
  time: Date;
  updated: Date;
  more?: React.ReactElement;
};

export function Byline(props: Props) {
  return (
    <Flex alignItems="start" justify="space-between" minWidth="0">
      <Flex
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
      </Flex>

      {props.more}
    </Flex>
  );
}
