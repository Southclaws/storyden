import { timestamp } from "src/utils/date";

import { styled } from "@/styled-system/jsx";

import { Anchor } from "./Anchor";

type Props = {
  created: string | Date;
  updated?: string | Date | undefined;
  href?: string;
  large?: boolean;
};

export function Timestamp(props: Props) {
  const { created, updated } = props;

  const createdDate = normaliseDate(created);
  const updatedDate = normaliseDate(updated);

  const showUpdated = isUpdatedLongerThanADay(createdDate, updatedDate);

  const createdAt = timestamp(createdDate, !props.large);
  const updatedAt = updated ? timestamp(updated, !props.large) : null;

  return (
    <styled.span>
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
      {showUpdated && <> (updated {updatedAt})</>}
    </styled.span>
  );
}

function normaliseDate(date: string | Date | undefined) {
  if (date === undefined) {
    return new Date();
  }

  return date instanceof Date ? date : new Date(date);
}

function isUpdatedLongerThanADay(created: Date, updated?: Date) {
  if (!updated) return false;

  return updated.getTime() - created.getTime() > 1000 * 60 * 60 * 24;
}
