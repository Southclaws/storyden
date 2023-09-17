import { Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";

import { PublicProfile } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";

export function Metadata(props: PublicProfile) {
  return (
    <Text size="md" color="gray.500">
      {"Registered "}
      <Timestamp
        created={formatDistanceToNow(new Date(props.createdAt), {
          addSuffix: true,
        })}
      />
    </Text>
  );
}
