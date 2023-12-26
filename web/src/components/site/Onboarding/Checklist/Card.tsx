"use client";

import { PropsWithChildren } from "react";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { CheckCircle } from "src/components/graphics/CheckCircle";
import { Heading1 } from "src/theme/components/Heading/Index";
import { Link } from "src/theme/components/Link";

import { Box, Circle, HStack, styled } from "@/styled-system/jsx";

import { Step, isComplete, statusToStep } from "./useChecklist";

type CardProps = {
  step: Step;
  current: OnboardingStatus;
  title: string;
  url?: string;
  onClick?: () => void;
};

export function Card(props: PropsWithChildren<CardProps>) {
  const complete = isComplete(props.step, props.current);
  const isCurrent = statusToStep[props.current] === props.step;
  return (
    <styled.li
      p="4"
      borderRadius="2xl"
      bgColor={complete ? "green.200" : "gray.100"}
    >
      <HStack gap="1">
        <Box>
          <Circle
            id="list-icon-circle"
            size="7"
            style={{
              backgroundColor: complete ? "none" : "gray.200",
            }}
          >
            {complete ? (
              <CheckCircle width="2em" height="2em" />
            ) : (
              <styled.p fontWeight="bold">{props.step}</styled.p>
            )}
          </Circle>
        </Box>

        <Box>
          <HStack justify="space-between">
            <Heading1 size="md">{props.title}</Heading1>

            {!complete && isCurrent && (
              <Link
                href={props.url ?? ""}
                bgColor="green.200"
                size="xs"
                onClick={props.onClick}
              >
                Complete
              </Link>
            )}
          </HStack>
          {props.children}
        </Box>
      </HStack>
    </styled.li>
  );
}
