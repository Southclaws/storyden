"use client";

import { PropsWithChildren } from "react";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { CheckCircle } from "src/components/graphics/CheckCircle";
import {
  Button,
  Heading,
  ListIcon,
  ListItem,
  Text,
} from "src/theme/components";

import { Box, Circle, HStack } from "@/styled-system/jsx";

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
    <ListItem
      p={4}
      borderRadius="2xl"
      bgColor={complete ? "green.200" : "gray.100"}
    >
      <HStack gap="1">
        <ListIcon
          id="list-icon"
          p={0}
          m={0}
          as={() => (
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
                <Text fontWeight="bold">{props.step}</Text>
              )}
            </Circle>
          )}
          fontSize="3xl"
        />

        <Box>
          <HStack justify="space-between">
            <Heading size="md">{props.title}</Heading>

            {!complete && isCurrent && (
              <Button
                as="a"
                href={props.url}
                bgColor="green.200"
                size="xs"
                onClick={props.onClick}
              >
                Complete
              </Button>
            )}
          </HStack>
          {props.children}
        </Box>
      </HStack>
    </ListItem>
  );
}
