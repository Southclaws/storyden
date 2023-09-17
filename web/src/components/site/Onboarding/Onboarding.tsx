"use client";

import {
  Box,
  Button,
  Circle,
  HStack,
  Heading,
  Link,
  ListIcon,
  ListItem,
  OrderedList,
  Text,
  VStack,
} from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { CheckCircle } from "src/components/graphics/CheckCircle";
import { CategoryCreateModal } from "src/components/site/Navigation/Sidebar/components/CategoryCreate/CategoryCreateModal";

import { useOnboarding } from "./useOnboarding";

type Props = {
  status: OnboardingStatus;
  onFinish: () => void;
};

type Step = 1 | 2 | 3 | 4 | 5;

type CardProps = {
  step: Step;
  current: OnboardingStatus;
  title: string;
  url?: string;
  onClick?: () => void;
};

const statusToStep: Record<OnboardingStatus, Step> = {
  requires_first_account: 1,
  requires_category: 2,
  requires_first_post: 3,
  requires_more_accounts: 4,
  complete: 5,
};

function Card(props: PropsWithChildren<CardProps>) {
  const complete = isComplete(props.step, props.current);
  const isCurrent = statusToStep[props.current] === props.step;
  return (
    <ListItem
      p={4}
      borderRadius="2xl"
      bgColor={complete ? "green.200" : "gray.100"}
    >
      <HStack gap={1}>
        <ListIcon
          id="list-icon"
          p={0}
          m={0}
          as={() => (
            <Circle
              id="list-icon-circle"
              as="span"
              bgColor={complete ? "none" : "gray.200"}
              size={7}
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

export function Onboarding(props: Props) {
  const { onOpen, isOpen, onClose, isLoggedIn } = useOnboarding();

  return (
    <VStack height="full" justify="start" pt={4} pb={16}>
      <Heading size="lg">Welcome to Storyden!</Heading>
      <Text p={2} textAlign="center">
        Get your community set up
        <br />
        with the following steps:
      </Text>
      <OrderedList display="flex" flexDir="column" gap={4} listStyleType="none">
        <Card
          step={1}
          current={props.status}
          title="Create an account"
          url="/auth"
        >
          <Text>
            Start by creating an account. The first registration is
            automatically given administrator rights!
          </Text>
        </Card>

        <Card
          step={2}
          current={props.status}
          title="Create a category"
          onClick={onOpen}
        >
          <Text>
            Posts need to live somewhere! So create your first category, give it
            a name and set it up just how you like!
          </Text>
          <CategoryCreateModal onClose={onClose} isOpen={isOpen} />
        </Card>

        <Card
          step={3}
          current={props.status}
          title="Write your first post"
          url="/new"
        >
          <Text>
            An intro, a thesis, a manifesto, a set of rules or just a hi! Get
            started on your first post in your new category!
          </Text>
        </Card>

        <VStack textAlign="center" px="20%">
          <Heading size="md">Invite your people</Heading>
          <Text>
            And you&apos;re ready to go! Spread the word and let the posts flow.{" "}
            <Link color="blue.400" href="https://www.storyden.org/docs">
              Visit the docs
            </Link>{" "}
            for more info if you get stuck.
          </Text>

          {isLoggedIn && (
            <Button onClick={props.onFinish}>
              {isComplete(3, props.status) ? "Finish" : "Skip"}
            </Button>
          )}
        </VStack>
      </OrderedList>
    </VStack>
  );
}

function isComplete(step: Step, status: OnboardingStatus) {
  return statusToStep[status] > step;
}
