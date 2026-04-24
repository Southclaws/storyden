"use client";

import { formatDate } from "date-fns";

import { EmailQueueItem, EmailQueueListResult } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import {
  Box,
  CardBox,
  HStack,
  LStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { CardBox as cardBox, lstack } from "@/styled-system/patterns";

import { useEmailLogSettingsScreen } from "./useEmailLogSettingsScreen";

export function EmailLogSettingsScreen() {
  const {
    ready,
    data,
    error,
    currentPage,
    refreshEmailLog,
    refreshing,
    retryEmail,
    retryingEmailId,
  } = useEmailLogSettingsScreen();

  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  return (
    <CardBox className={lstack()} gap="4">
      <WStack>
        <Heading size="md">Email Log</Heading>

        <Button
          type="button"
          size="xs"
          variant="subtle"
          loading={refreshing}
          onClick={refreshEmailLog}
        >
          Refresh
        </Button>
      </WStack>
      <styled.p>
        View queued emails and delivery attempts, including failure reasons.
      </styled.p>

      <EmailLogList
        data={data}
        currentPage={currentPage}
        onRetry={retryEmail}
        retryingEmailId={retryingEmailId}
      />
    </CardBox>
  );
}

function EmailLogList({
  data,
  currentPage,
  onRetry,
  retryingEmailId,
}: {
  data: EmailQueueListResult;
  currentPage: number;
  onRetry: (email: EmailQueueItem) => Promise<void>;
  retryingEmailId: string | null;
}) {
  if (data.emails.length === 0) {
    return <EmptyState hideContributionLabel>No emails found.</EmptyState>;
  }

  return (
    <>
      <ul className={lstack({ gap: "3" })}>
        {data.emails.map((email) => (
          <EmailItem
            key={email.id}
            email={email}
            onRetry={onRetry}
            retrying={retryingEmailId === email.id}
          />
        ))}
      </ul>

      <PaginationControls
        path="/admin"
        currentPage={currentPage}
        totalPages={data.total_pages}
        pageSize={data.page_size}
        params={{
          tab: "email",
        }}
      />
    </>
  );
}

function EmailItem({
  email,
  onRetry,
  retrying,
}: {
  email: EmailQueueItem;
  onRetry: (email: EmailQueueItem) => Promise<void>;
  retrying: boolean;
}) {
  const lastAttempt = email.attempts[email.attempts.length - 1] ?? null;
  const error = lastAttempt?.error ?? null;
  const recipient = email.recipient_name
    ? `${email.recipient_name} <${email.recipient_address}>`
    : email.recipient_address;

  const failures = email.attempts.filter((attempt) => attempt.error).length;

  const attemptCountLabel =
    failures === 1 ? "1 failure" : `${failures} failures`;
  const attempts = [...email.attempts].reverse();

  return (
    <li className={cardBox()}>
      <LStack gap="2">
        <WStack gap="2" alignItems="center">
          <HStack>
            <Badge
              variant="solid"
              colorPalette={email.status === "failed" ? "red" : "gray"}
            >
              {email.status}
            </Badge>
          </HStack>

          <styled.time fontSize="xs" color="fg.muted">
            {formatDate(email.queued_at, "PPpp")}
          </styled.time>
        </WStack>

        <WStack>
          <styled.p fontSize="sm" color="fg.subtle">
            “{email.subject}”
          </styled.p>

          <styled.p fontSize="xs" color="fg.subtle">
            {recipient}
          </styled.p>
        </WStack>

        <styled.details
          borderWidth="thin"
          borderColor="border.subtle"
          borderRadius="md"
          bg="bg.subtle"
          overflow="hidden"
        >
          <styled.summary
            listStyle="none"
            cursor="pointer"
            px="3"
            py="2.5"
            fontSize="sm"
            fontWeight="medium"
            color="fg.default"
            _marker={{ display: "none" }}
            css={{
              "&::-webkit-details-marker": {
                display: "none",
              },
            }}
          >
            <LStack gap="1">
              <styled.span fontSize="xs" color="fg.muted" fontWeight="normal">
                {error ? error : "Open to review the full delivery history."}
              </styled.span>
            </LStack>
          </styled.summary>

          <LStack gap="2" px="3" pb="3">
            {attempts.map((attempt, index) => (
              <Box
                key={`${attempt.timestamp}-${index}`}
                borderTopWidth={index === 0 ? "none" : "thin"}
                borderColor="border.subtle"
                pt={index === 0 ? "0" : "2"}
              >
                <Box
                  display="grid"
                  gap="2"
                  alignItems="start"
                  gridTemplateColumns={{
                    base: "1fr",
                    md: "auto minmax(0, 1fr)",
                  }}
                >
                  <WStack gap="0.5">
                    <Badge
                      size="sm"
                      variant="solid"
                      colorPalette={
                        attempt.status === "failed" ? "red" : "gray"
                      }
                    >
                      {attempt.status}
                    </Badge>

                    <styled.time fontSize="xs" color="fg.muted">
                      {formatDate(attempt.timestamp, "PPpp")}
                    </styled.time>
                  </WStack>

                  <styled.pre
                    m="0"
                    minW="0"
                    fontSize="xs"
                    lineHeight="normal"
                    color={attempt.error ? "fg.default" : "fg.muted"}
                    whiteSpace="pre-wrap"
                    overflowWrap="anywhere"
                    wordBreak="break-word"
                  >
                    {attempt.error ?? "No error recorded."}
                  </styled.pre>
                </Box>
              </Box>
            ))}
          </LStack>
        </styled.details>

        <WStack alignItems="end">
          <Badge size="sm">{email.attempts.length} attempts</Badge>

          {email.status === "failed" && (
            <Button
              type="button"
              size="xs"
              variant="subtle"
              loading={retrying}
              loadingText="Retrying..."
              onClick={() => onRetry(email)}
            >
              Retry Now
            </Button>
          )}
        </WStack>
      </LStack>
    </li>
  );
}
