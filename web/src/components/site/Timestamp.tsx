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
    <styled.span textWrap="nowrap">
      {props.href ? (
        <Anchor href={props.href}>
          {props.large && (
            <styled.span className="fluid-font-size">created</styled.span>
          )}{" "}
          {createdAt}
        </Anchor>
      ) : (
        <styled.span>{createdAt}</styled.span>
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
