import { UIDataTypes, UIMessagePart, isToolUIPart } from "ai";
import { match } from "ts-pattern";

import { useNodeList } from "@/api/openapi-client/nodes";
import { NodeWithChildren } from "@/api/openapi-schema";
import {
  StorydenTools,
  ToolLibraryRequestPageOutput,
  ToolRobotDeleteInput,
} from "@/api/robots";
import {
  type StorydenUIMessage,
  getToolName,
  isToolType,
} from "@/api/robots-types";
import { Button } from "@/components/ui/button";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { Box, HStack, LStack, styled } from "@/styled-system/jsx";
import { wstack } from "@/styled-system/patterns";

import { useRobotChat } from "./RobotChatContext";

type Props = {
  part: UIMessagePart<UIDataTypes, StorydenTools>;
};

export type ConfirmationPart = {
  type: `tool-${string}`;
  toolCallId: string;
  toolName?: string;
  state?: string;
  input?: unknown;
  output?: unknown;
  approval?: {
    id: string;
    approved?: boolean;
    reason?: string;
  };
};

export function RobotToolCall({ part }: Props) {
  const toolName = getToolName(part);

  return (
    <LStack
      role="group"
      aria-label={`${toolName} tool call`}
      className="group"
      gap="1.5"
      pl="3"
      pr="2"
      py="1.5"
      bg="transparent"
      borderLeftWidth="medium"
      borderLeftColor="border.subtle"
      borderLeftRadius="none"
      fontSize="sm"
      w="full"
      alignSelf="flex-start"
    >
      <RobotToolCallTitle part={part} toolName={toolName} />
      <RobotToolConfirmation part={part} />
      <RobotLibraryPageRequest part={part} />
      <RobotToolCallContent part={part} />
    </LStack>
  );
}

function RobotToolCallTitle({ part, toolName }: Props & { toolName: string }) {
  const toolDetails = JSON.stringify(part, null, 2);

  return (
    <styled.details w="full" aria-label={`${toolName} tool call details`}>
      <styled.summary
        className={wstack()}
        aria-label={`${toolName} tool call details`}
        cursor="pointer"
        fontSize="xs"
        color="fg.muted"
        _groupHover={{ color: "fg.default" }}
      >
        <styled.span
          display="flex"
          alignItems="center"
          _detailsOpen={{
            color: "fg.default",
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
  if (!isToolUIPart(part)) {
    return null;
  }

  if (!part.output) {
    return null;
  }

  if (!isToolType(part.type)) {
    return null;
  }

  switch (part.type) {
    case "tool-content_search":
    case "tool-thread_search":
    case "tool-reply_search":
    case "tool-post_search":
      return <p>{part.output.results} results found</p>;

    case "tool-member_search":
      return <p>{part.output.results} members found</p>;

    case "tool-library_search_pages":
      return <p>{part.output.results} pages found</p>;

    case "tool-robot_switch":
      return null;

    case "tool-system_robot_tool_catalog":
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
      if (isConfirmationDenied(part.output)) {
        return <p>Delete cancelled</p>;
      }
      return <p>Deleted robot</p>;

    case "tool-library_page_list":
      return null;

    case "tool-library_request_page":
      return <p>Selected "{part.output.name}"</p>;

    case "tool-get_library_page":
      return null;

    case "tool-create_library_page":
      return <p>Created "{part.output.name}"</p>;

    case "tool-update_library_page":
      return <p>Updated "{part.output.name}"</p>;

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
  const { messages } = useRobotChat();

  if (!isToolUIPart(part)) {
    return null;
  }

  const confirmationResolution = getToolConfirmationResolution(part, messages);

  const label = match(part.state)
    .with("approval-requested", () =>
      confirmationResolution === "denied"
        ? "Denied"
        : confirmationResolution === "approved"
          ? "Tool complete"
          : "Needs approval",
    )
    .with("approval-responded", () => "Tool complete")
    .with("input-available", () =>
      part.type === "tool-library_request_page"
        ? "Needs selection"
        : "Running tool",
    )
    .with("input-streaming", () => "Running tool")
    .with("output-available", () => "Tool complete")
    .with("output-error", () => "Error")
    .with("output-denied", () => "Denied")
    .otherwise(() => "Tool complete");

  return <styled.span>{label}</styled.span>;
}

function RobotToolConfirmation({ part }: Props) {
  const { messages, resolveToolConfirmation } = useRobotChat();

  if (!isToolUIPart(part)) {
    return null;
  }

  if (
    part.state !== "approval-requested" ||
    !part.approval?.id ||
    getToolConfirmationResolution(part, messages)
  ) {
    return null;
  }

  return (
    <HStack gap="2" justifyContent="flex-start">
      <Button
        aria-label="Approve"
        size="xs"
        variant="solid"
        onClick={() =>
          resolveToolConfirmation({
            approvalId: part.approval.id,
            toolName: getToolPartName(part),
            approved: true,
          })
        }
      >
        Approve
      </Button>
      <Button
        aria-label="Deny"
        size="xs"
        variant="outline"
        onClick={() =>
          resolveToolConfirmation({
            approvalId: part.approval.id,
            toolName: getToolPartName(part),
            approved: false,
          })
        }
      >
        Deny
      </Button>
    </HStack>
  );
}

export function RobotToolConfirmationBatch({
  parts,
}: {
  parts: ConfirmationPart[];
}) {
  const { messages, resolveToolConfirmation, robots } = useRobotChat();
  const pendingParts = parts.filter(
    (part) =>
      part.state === "approval-requested" &&
      !getToolConfirmationResolution(part, messages),
  );
  const hasPending = pendingParts.length > 0;
  const isPartiallyResolved = hasPending && pendingParts.length !== parts.length;

  const resolvePart = (part: ConfirmationPart, approved: boolean) =>
    part.approval?.id
      ? resolveToolConfirmation({
          approvalId: part.approval.id,
          toolName: getToolPartName(part),
          approved,
        })
      : Promise.resolve();

  const resolveAll = async (approved: boolean) => {
    for (const part of pendingParts) {
      await resolvePart(part, approved);
    }
  };

  return (
    <LStack
      role="group"
      aria-label={isPartiallyResolved ? "Partial approvals" : "Confirmation batch"}
      className="group"
      gap="2"
      pl="3"
      pr="2"
      py="1.5"
      bg="transparent"
      borderLeftWidth="medium"
      borderLeftColor="border.subtle"
      borderLeftRadius="none"
      fontSize="sm"
      w="full"
      alignSelf="flex-start"
    >
      <styled.div className={wstack()} fontSize="xs" color="fg.muted">
        <styled.span display="flex" alignItems="center">
          <ToolIcon w="3" h="3" />
          &nbsp;{hasPending ? `Approve these ${parts.length} actions?` : `${parts.length} actions resolved`}
        </styled.span>
        <styled.span>
          {hasPending ? `${pendingParts.length} pending` : "Tool complete"}
        </styled.span>
      </styled.div>

      <LStack as="ul" gap="2" alignItems="stretch" w="full">
        {parts.map((part, index) => (
          <ConfirmationBatchRow
            key={part.toolCallId}
            part={part}
            resolution={getToolConfirmationResolution(part, messages)}
            label={formatConfirmationAction(part, index, robots)}
            onApprove={() => resolvePart(part, true)}
            onDeny={() => resolvePart(part, false)}
          />
        ))}
      </LStack>

      {hasPending ? (
        <HStack gap="2" justifyContent="flex-start" pt="1">
          <Button
            aria-label="Approve all confirmations"
            size="xs"
            variant="solid"
            onClick={() => resolveAll(true)}
          >
            Approve all
          </Button>
          <Button
            aria-label="Deny all confirmations"
            size="xs"
            variant="outline"
            onClick={() => resolveAll(false)}
          >
            Deny all
          </Button>
        </HStack>
      ) : null}
    </LStack>
  );
}

function ConfirmationBatchRow({
  part,
  resolution,
  label,
  onApprove,
  onDeny,
}: {
  part: ConfirmationPart;
  resolution?: ConfirmationResolution;
  label: string;
  onApprove: () => void;
  onDeny: () => void;
}) {
  const pending = part.state === "approval-requested" && !resolution;
  const denied =
    resolution === "denied" ||
    part.state === "output-denied" ||
    part.approval?.approved === false ||
    isConfirmationDenied(part.output);

  return (
    <HStack
      as="li"
      gap="2"
      justifyContent="space-between"
      alignItems="center"
      w="full"
      listStyle="none"
    >
      <HStack gap="1.5" minW="0" color="fg.muted">
        <ToolIcon w="3" h="3" flexShrink="0" />
        <styled.span overflow="hidden" textOverflow="ellipsis">
          {label}
        </styled.span>
      </HStack>

      {pending ? (
        <HStack gap="1.5" flexShrink="0">
          <Button
            aria-label={`Approve ${label}`}
            size="xs"
            variant="solid"
            onClick={onApprove}
          >
            Approve
          </Button>
          <Button
            aria-label={`Deny ${label}`}
            size="xs"
            variant="outline"
            onClick={onDeny}
          >
            Deny
          </Button>
        </HStack>
      ) : (
        <styled.span fontSize="xs" color={denied ? "fg.muted" : "fg.default"}>
          {denied ? "Denied" : "Approved"}
        </styled.span>
      )}
    </HStack>
  );
}

function formatConfirmationAction(
  part: ConfirmationPart,
  index: number,
  robots: readonly { id: string; name: string }[],
) {
  if (part.type === "tool-robot_delete") {
    const input = part.input as ToolRobotDeleteInput | undefined;
    const robot = robots.find((robot) => robot.id === input?.id);
    if (robot) {
      return `Delete ${robot.name}`;
    }
    return input?.id ? `Delete Robot ${input.id}` : `Delete Robot ${index + 1}`;
  }

  return `Action ${index + 1}`;
}

export function isConfirmationToolPart(
  part: UIMessagePart<UIDataTypes, StorydenTools>,
): boolean {
  return (
    isToolUIPart(part) &&
    (part.state === "approval-requested" ||
      part.state === "approval-responded" ||
      part.state === "output-denied")
  );
}

function getToolPartName(part: { type: string; toolName?: string }) {
  if (part.toolName) {
    return part.toolName;
  }
  if (part.type.startsWith("tool-")) {
    return part.type.slice("tool-".length);
  }
  return undefined;
}

type ConfirmationResolution = "approved" | "denied";

function getToolConfirmationResolution(
  part: { toolCallId?: string; approval?: { id?: string } },
  messages: readonly StorydenUIMessage[],
): ConfirmationResolution | undefined {
  const approvalID = part.approval?.id ?? part.toolCallId;
  if (!approvalID) {
    return undefined;
  }

  for (const message of messages) {
    for (const messagePart of message.parts ?? []) {
      if (!isToolUIPart(messagePart) || messagePart.toolCallId !== approvalID) {
        continue;
      }

      if (messagePart.state === "approval-responded") {
        if (messagePart.approval?.approved === false) {
          return "denied";
        }
        if (messagePart.approval?.approved === true) {
          return "approved";
        }
      }

      if (
        isAdkRequestConfirmationPart(messagePart) &&
        messagePart.state === "output-available" &&
        "output" in messagePart
      ) {
        const confirmed = getConfirmationOutputDecision(messagePart.output);
        if (confirmed === false) {
          return "denied";
        }
        if (confirmed === true) {
          return "approved";
        }
      }

      if (messagePart.state === "output-denied") {
        return "denied";
      }
    }
  }

  return undefined;
}

function isAdkRequestConfirmationPart(part: { type: string }) {
  return part.type === "tool-adk_request_confirmation";
}

function getConfirmationOutputDecision(output: unknown): boolean | undefined {
  if (!output || typeof output !== "object" || !("confirmed" in output)) {
    return undefined;
  }

  const confirmed = output.confirmed;
  return typeof confirmed === "boolean" ? confirmed : undefined;
}

function RobotLibraryPageRequest({ part }: Props) {
  const { resolveLibraryPageRequest } = useRobotChat();
  const { data, isLoading, error } = useNodeList({ depth: "-1" });

  if (!isToolUIPart(part)) {
    return null;
  }

  if (
    part.type !== "tool-library_request_page" ||
    part.state !== "input-available"
  ) {
    return null;
  }

  if (isLoading) {
    return <styled.p color="fg.muted">Loading pages...</styled.p>;
  }

  if (error) {
    return <styled.p color="fg.error">Could not load Library pages.</styled.p>;
  }

  const pages = flattenLibraryPages(data?.nodes ?? []);

  if (pages.length === 0) {
    return <styled.p color="fg.muted">No Library pages found.</styled.p>;
  }

  return (
    <LStack
      role="group"
      aria-label="Choose a Library page"
      gap="1"
      alignItems="stretch"
      maxH="48"
      overflowY="auto"
      w="full"
    >
      {pages.map((page) => {
        const output: ToolLibraryRequestPageOutput = {
          id: page.id,
          slug: page.slug,
          name: page.name,
          description: page.description,
        };

        return (
          <Button
            key={page.id}
            aria-label={`Select Library page ${page.name}`}
            size="xs"
            variant="outline"
            justifyContent="flex-start"
            h="auto"
            py="1.5"
            px="2"
            textAlign="left"
            onClick={() =>
              resolveLibraryPageRequest({
                toolCallId: part.toolCallId,
                page: output,
              })
            }
          >
            <LStack gap="0" alignItems="flex-start" overflow="hidden">
              <styled.span
                maxW="full"
                overflow="hidden"
                textOverflow="ellipsis"
              >
                {page.name}
              </styled.span>
              <styled.span
                color="fg.muted"
                fontSize="xs"
                maxW="full"
                overflow="hidden"
                textOverflow="ellipsis"
              >
                /{page.slug}
              </styled.span>
            </LStack>
          </Button>
        );
      })}
    </LStack>
  );
}

function flattenLibraryPages(nodes: NodeWithChildren[]): NodeWithChildren[] {
  return nodes.flatMap((node) => [node, ...flattenLibraryPages(node.children)]);
}

function isConfirmationDenied(output: unknown) {
  if (!output || typeof output !== "object") {
    return false;
  }
  const maybe = output as { _storyden_confirmation?: { approved?: boolean } };
  return maybe._storyden_confirmation?.approved === false;
}
