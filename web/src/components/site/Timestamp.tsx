import { timestamp } from "src/utils/date";

import { styled } from "@/styled-system/jsx";

import { Anchor } from "./Anchor";

type Props = {
  created: string | Date;
  href?: string;
  large?: boolean;
};

export function Timestamp(props: Props) {
  const { created } = props;

  const createdDate = normaliseDate(created);

  const createdAt = timestamp(createdDate, !props.large);

  return (
    <styled.span
      className="timestamp__container"
      textWrap="nowrap"
      minW="fit"
      flexShrink="0"
      overflow="hidden"
    >
      {props.href ? (
        <Anchor className="timestamp__anchor" href={props.href}>
          {props.large && (
            <styled.span className="timestamp__label fluid-font-size">
              created
            </styled.span>
          )}{" "}
          {createdAt}
        </Anchor>
      ) : (
        <styled.span className="timestamp__time">{createdAt}</styled.span>
      )}
    </styled.span>
  );
}

function normaliseDate(date: string | Date | undefined) {
  if (date === undefined) {
    return new Date();
  }

  return date instanceof Date ? date : new Date(date);
}
