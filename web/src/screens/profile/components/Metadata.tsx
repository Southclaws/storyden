import { formatDistanceToNow } from "date-fns";

import { PublicProfile } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";

import { styled } from "@/styled-system/jsx";

export function Metadata(props: PublicProfile) {
  return (
    <>
      <styled.p color="fg.muted">
        Registered&nbsp;
        <Timestamp
          created={formatDistanceToNow(new Date(props.createdAt), {
            addSuffix: true,
          })}
        />
      </styled.p>

      {props.deletedAt && (
        <styled.p color="fg.destructive" wordBreak="keep-all">
          Suspended&nbsp;
          <styled.time textWrap="nowrap">
            {formatDistanceToNow(new Date(props.deletedAt), {
              addSuffix: true,
            })}
          </styled.time>
        </styled.p>
      )}
    </>
  );
}
