import {
  Box,
  Circle,
  HStack,
  Heading,
  ListIcon,
  ListItem,
  OrderedList,
  Text,
  VStack,
} from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { CheckCircle } from "src/components/graphics/CheckCircle";

type Props = {
  status: OnboardingStatus;
};

type Step = 1 | 2 | 3 | 4 | 5;

type CardProps = { n: Step; c: boolean };

const statusToStep: Record<OnboardingStatus, Step> = {
  requires_first_account: 1,
  requires_category: 2,
  requires_more_accounts: 3,
  requires_first_post: 4,
  complete: 5,
};

function Card(props: PropsWithChildren<CardProps>) {
  return (
    <ListItem
      p={4}
      borderRadius="2xl"
      bgColor={props.c ? "green.200" : "gray.100"}
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
              bgColor={props.c ? "none" : "gray.200"}
              size={7}
            >
              {props.c ? (
                <CheckCircle width="2em" height="2em" />
              ) : (
                <Text fontWeight="bold">{props.n}</Text>
              )}
            </Circle>
          )}
          fontSize="3xl"
        />
        <Box>{props.children}</Box>
      </HStack>
    </ListItem>
  );
}

export function Onboarding(props: Props) {
  return (
    <VStack height="full" justify="center" pb={16}>
      <Heading size="lg">Welcome to Storyden!</Heading>
      <Text p={2} textAlign="center">
        Get your community set up
        <br />
        with the following steps:
      </Text>
      <OrderedList display="flex" flexDir="column" gap={4} listStyleType="none">
        <Card n={1} c={isComplete(1, props.status)}>
          <Heading size="md">Create an account</Heading>
          <Text>
            Start by creating an account. The first registration is
            automatically given administrator rights!
          </Text>
        </Card>

        <Card n={2} c={isComplete(2, props.status)}>
          <Heading size="md">Create a category</Heading>
          <Text>
            Posts need to live somewhere! So create your first category, give it
            a name and set it up just how you like!
          </Text>
        </Card>

        <Card n={3} c={isComplete(3, props.status)}>
          <Heading size="md">Invite your people</Heading>
          <Text>
            And you&apos;re ready to go! Spread the word and let the posts flow.
            Visit the docs for more info if you get stuck.
          </Text>
        </Card>
      </OrderedList>
    </VStack>
  );
}

function isComplete(step: Step, status: OnboardingStatus) {
  return statusToStep[status] > step;
}
