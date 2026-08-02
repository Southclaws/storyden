"use client";

import { PropsWithChildren } from "react";

import { OnboardingStatus } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import { Box, Circle, HStack, styled } from "@/styled-system/jsx";

import { Step, isComplete, statusToStep } from "./useChecklist";

type CardProps = {
  step: Step;
  current: OnboardingStatus;
  title: string;
  url?: string;
  onClick?: () => void;
};

export function ChecklistItem(props: PropsWithChildren<CardProps>) {
  const complete = isComplete(props.step, props.current);
  const isCurrent = statusToStep[props.current] === props.step;
  return (
    <styled.li
      p="4"
      w="full"
      maxW="prose"
      borderRadius="2xl"
      bgColor={complete ? "status.success.surface" : "background.inset"}
    >
      <HStack w="full" gap="2">
        <Box>
          <Circle
            id="list-icon-circle"
            size="7"
            style={{
              backgroundColor: complete ? "none" : "gray.200",
            }}
          >
            {complete ? (
              <CheckCircleIcon
                width="8"
                height="8"
                color="status.success.content"
                fill="status.success.surface"
              />
            ) : (
              <Text fontWeight="bold">{props.step}</Text>
            )}
          </Circle>
        </Box>

        <Box w="full">
          <HStack justify="space-between">
            <Text
              variant="supporting"
              color="text.default"
              fontWeight="semibold"
            >
              {props.title}
            </Text>

            {!complete &&
              isCurrent &&
              (props.url ? (
                <LinkButton href={props.url} colorPalette="green">
                  Complete
                </LinkButton>
              ) : (
                <Button colorPalette="green" onClick={props.onClick}>
                  Complete
                </Button>
              ))}
          </HStack>
          {props.children}
        </Box>
      </HStack>
    </styled.li>
  );
}
