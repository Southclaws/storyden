import { UIDataTypes, UIMessagePart, isToolOrDynamicToolUIPart } from "ai";
import { match } from "ts-pattern";

import { StorydenTools } from "@/api/robots";
import { getToolName, isToolType } from "@/api/robots-types";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { Box, LStack, styled } from "@/styled-system/jsx";
import { wstack } from "@/styled-system/patterns";

type Props = {
  part: UIMessagePart<UIDataTypes, StorydenTools>;
};

export function RobotToolCall({ part }: Props) {
  return (
    <LStack
      gap="2"
      p="2"
      bg="bg.muted"
      borderRadius="md"
      fontSize="sm"
      w="full"
      maxW="5/6"
      alignSelf="flex-start"
    >
      <RobotToolCallTitle part={part} />
      <RobotToolCallContent part={part} />
    </LStack>
  );
}

function RobotToolCallTitle({ part }: Props) {
  const toolName = getToolName(part);
  const toolDetails = JSON.stringify(part, null, 2);

  return (
    <styled.details w="full">
      <styled.summary
        className={wstack()}
        cursor="pointer"
        fontSize="xs"
        color="fg.muted"
      >
        <styled.span
          display="flex"
          alignItems="center"
          _detailsOpen={{
            color: "black",
          }}
        >
          <ToolIcon w="3" h="3" />
          &nbsp;{toolName}
        </styled.span>
        <RobotToolCallStatus part={part} />
      </styled.summary>
      <Box>
        <styled.pre
          fontSize="xs"
          p="2"
          bg="bg.subtle"
          borderRadius="sm"
          overflow="auto"
          maxH="32"
          mt="1"
        >
          {toolDetails}
        </styled.pre>
      </Box>
    </styled.details>
  );
}

function RobotToolCallContent({ part }: Props) {
  if (!isToolOrDynamicToolUIPart(part)) {
    return null;
  }

  if (!part.output) {
    return null;
  }

  if (!isToolType(part.type)) {
    return null;
  }

  switch (part.type) {
    case "tool-search":
      return <p>{part.output.results} results found</p>;

    case "tool-robot_switch":
      return null;

    case "tool-system_all_tool_names":
      return <p>{part.output.tools?.length ?? 0} tools available</p>;

    case "tool-robot_create":
      return <p>Created "{part.output.name}"</p>;

    case "tool-robot_list":
      return <p>{part.output.total} robots</p>;

    case "tool-robot_get":
      return null;

    case "tool-robot_update":
      return part.output.name ? <p>Updated "{part.output.name}"</p> : null;

    case "tool-robot_delete":
      return <p>Deleted robot</p>;

    case "tool-library_page_list":
      return null;

    case "tool-get_library_page":
      return null;

    case "tool-create_library_page":
      return <p>Created "{part.output.name}"</p>;

    case "tool-update_library_page":
      return <p>Updated "{part.output.name}"</p>;

    case "tool-search_library_pages":
      return <p>{part.output.items.length} pages found</p>;

    case "tool-library_page_property_schema_get":
      return null;

    case "tool-library_page_property_schema_update":
      return <p>Updated property schema</p>;

    case "tool-library_page_properties_update":
      return <p>Updated properties</p>;

    case "tool-tag_list":
      return null;

    case "tool-link_create":
      return <p>Created link</p>;

    case "tool-thread_create":
      return <p>Created "{part.output.title}"</p>;

    case "tool-thread_list":
      return null;

    case "tool-thread_get":
      return null;

    case "tool-thread_update":
      return <p>Updated "{part.output.title}"</p>;

    case "tool-thread_reply":
      return <p>Posted reply</p>;

    case "tool-category_list":
      return null;
  }
}

function RobotToolCallStatus({ part }: Props) {
  if (!isToolOrDynamicToolUIPart(part)) {
    return null;
  }

  const label = match(part.state)
    .with("input-available", () => "Running tool")
    .with("input-streaming", () => "Running tool")
    .with("output-available", () => "Tool complete")
    .with("output-error", () => "Error")
    .otherwise(() => "Tool complete");

  return <styled.span>{label}</styled.span>;
}
