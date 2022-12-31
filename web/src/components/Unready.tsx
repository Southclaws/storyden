import { Box, Flex, Spinner } from "@chakra-ui/react";
import { APIError } from "src/api/openapi/schemas";

export function Unready(props: Partial<APIError>) {
  if (!props.error)
    return (
      <Flex width="full" justifyContent="center" p={12}>
        <Spinner
          thickness="4px"
          speed="0.65s"
          color="hsl(0, 0%, 75%)"
          size="xl"
        />
      </Flex>
    );

  return <Box>{props.message ?? props.error}</Box>;
}
