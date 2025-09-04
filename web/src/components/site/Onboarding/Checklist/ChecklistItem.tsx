"use client";

import { PropsWithChildren } from "react";

import { OnboardingStatus } from "src/api/openapi-schema";

import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { LinkButton } from "@/components/ui/link-button";
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
      bgColor={complete ? "bg.success" : "bg.subtle"}
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
                color="fg.success"
                fill="bg.success"
              />
            ) : (
              <styled.p fontWeight="bold">{props.step}</styled.p>
            )}
          </Circle>
        </Box>

        <Box w="full">
          <HStack justify="space-between">
            <Heading size="md">{props.title}</Heading>

            {!complete &&
              isCurrent &&
              (props.url ? (
                <LinkButton href={props.url} colorPalette="green" size="xs">
                  Complete
                </LinkButton>
              ) : (
                <Button colorPalette="green" size="xs" onClick={props.onClick}>
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
