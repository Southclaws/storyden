import { timestamp } from "src/utils/date";

import { styled } from "@/styled-system/jsx";
import { JsxStyleProps } from "@/styled-system/types";

import { Anchor } from "./Anchor";

type Props = {
  created: string | Date;
  href?: string;
  large?: boolean;
};

export function Timestamp(props: Props & JsxStyleProps) {
  const { created, href, large, ...rest } = props;

  const createdDate = normaliseDate(created);

  const createdAt = timestamp(createdDate, !large);

  return (
    <styled.span
      className="timestamp__container"
      textWrap="nowrap"
      minW="fit"
      flexShrink="0"
      overflow="hidden"
      {...rest}
    >
      {href ? (
        <Anchor className="timestamp__anchor" href={href}>
          {large && (
            <styled.span className="timestamp__label fluid-font-size">
              created
            </styled.span>
          )}{" "}
          {createdAt}
        </Anchor>
      ) : (
        <styled.time
          className="timestamp__time"
          dateTime={createdDate.toISOString()}
        >
          {createdAt}
        </styled.time>
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
