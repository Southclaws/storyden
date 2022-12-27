import { Box, Spinner } from "@chakra-ui/react";
import { APIError } from "src/api/openapi/schemas";

export function Unready(props: Partial<APIError>) {
  if (!props.error) return <Spinner />;

  return <Box>{props.message ?? props.error}</Box>;
}
