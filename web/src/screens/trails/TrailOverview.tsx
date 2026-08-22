import { Trail } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { CalendarIcon } from "@/components/ui/icons/Calendar";
import { Text } from "@/components/ui/text";
import { describeTrailSchedule, formatOccurrence } from "@/lib/trails";
import { Box, HStack, WStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

import { TrailStatusBadge } from "./TrailStatusBadge";

export function TrailOverview({ trail }: { trail: Trail }) {
  return (
    <styled.section aria-label="Trail details">
      <WStack alignItems="center" gap="4" flexWrap="wrap" marginTop="4">
        <HStack gap="4" minWidth="0" flexWrap="wrap">
          <MemberBadge
            profile={trail.created_by}
            size="xs"
            name="handle"
            as="link"
          />

          <Text as="span" variant="metadata" flexShrink="0">
            {trail.actions.length} {pluralise(trail.actions.length, "action")}
          </Text>
        </HStack>

        <TrailStatusBadge status={trail.status} />
      </WStack>

      {trail.description && (
        <Text
          variant="supporting"
          display="block"
          overflowWrap="anywhere"
          marginBottom="4"
        >
          {trail.description}
        </Text>
      )}

      <styled.div
        display="grid"
        gridTemplateColumns={{
          base: "minmax(0, 1fr) minmax(0, 1fr)",
          md: "minmax(0, 1fr) auto auto",
        }}
        columnGap={{ base: "4", md: "8" }}
        rowGap="3"
        alignItems="start"
      >
        <Box minWidth="0" gridColumn={{ base: "1 / -1", md: "auto" }}>
          <HStack gap="1" alignItems="center">
            <CalendarIcon
              width="4"
              height="4"
              color="text.muted"
              flexShrink="0"
            />
            <Text
              as="span"
              variant="metadata"
              fontWeight="bold"
              textTransform="uppercase"
            >
              Schedule
            </Text>
          </HStack>
          <Text
            as="span"
            fontWeight="semibold"
            overflowWrap="anywhere"
            display="block"
          >
            {describeTrailSchedule(trail.trigger.schedule)}
          </Text>
          <Text as="span" variant="metadata" display="block">
            {trail.trigger.schedule.timezone}
          </Text>
        </Box>

        <Occurrence
          label="Next run"
          value={trail.next_occurrence_at}
          empty="No future run"
          emphasis
        />
      </styled.div>
    </styled.section>
  );
}

function Occurrence({
  label,
  value,
  empty,
  emphasis = false,
}: {
  label: string;
  value?: string;
  empty: string;
  emphasis?: boolean;
}) {
  return (
    <Box minWidth="0">
      <Text
        as="span"
        variant="metadata"
        fontWeight="bold"
        textTransform="uppercase"
        display="block"
      >
        {label}
      </Text>
      {value ? (
        <styled.time
          dateTime={value}
          fontSize="sm"
          fontWeight={emphasis ? "semibold" : "normal"}
          display="block"
        >
          {formatOccurrence(value)}
        </styled.time>
      ) : (
        <Text as="span" variant="supporting" display="block">
          {empty}
        </Text>
      )}
    </Box>
  );
}
