import { Box } from "@chakra-ui/react";

import { PostContent } from "src/api/openapi/schemas";

import { Markdown } from "../Markdown";

type Props = {
  value: PostContent;
};

export function ContentViewer({ value }: Props) {
  switch (value.type) {
    case "text/markdown":
      return (
        <Box>
          <Markdown>{value.value ?? ""}</Markdown>
        </Box>
      );

    case "application/json":
      return <Box>TODO: Rich text renderer... {JSON.stringify(value)}</Box>;

    default:
      console.error(`ContentViewer unexpected content type: ${value.type}`);
      return <></>;
  }
}
