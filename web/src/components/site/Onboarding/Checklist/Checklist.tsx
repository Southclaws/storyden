"use client";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { CategoryCreateModal } from "src/components/category/CategoryCreate/CategoryCreateModal";
import { Button, Heading, Link, OrderedList, Text } from "src/theme/components";

import { VStack } from "@/styled-system/jsx";

import { Card } from "./Card";
import { isComplete, useChecklist } from "./useChecklist";

type Props = {
  onboardingStatus: OnboardingStatus;
  onFinish: () => void;
};

export function Checklist({ onboardingStatus, onFinish }: Props) {
  const { onOpen, isOpen, onClose, isLoggedIn } = useChecklist();

  return (
    <VStack width="full" height="full" justify="start" pt="4" pb="16">
      <Heading size="lg">Welcome to Storyden!</Heading>
      <Text p={2} textAlign="center">
        Get your community set up
        <br />
        with the following steps:
      </Text>
      <OrderedList
        display="flex"
        flexDir="column"
        gap={4}
        listStyleType="none"
        m={0}
      >
        <Card
          step={1}
          current={onboardingStatus}
          title="Create an account"
          url="/register"
        >
          <Text>
            Start by creating an account. The first registration is
            automatically given administrator rights!
          </Text>
        </Card>

        <Card
          step={2}
          current={onboardingStatus}
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
          current={onboardingStatus}
          title="Write your first post"
          url="/new"
        >
          <Text>
            An intro, a thesis, a manifesto, a set of rules or just a hi! Get
            started on your first post in your new category!
          </Text>
        </Card>

        <VStack textAlign="center" px="2">
          <Heading size="md">Invite your people</Heading>
          <Text>
            And you&apos;re ready to go! Spread the word and let the posts flow.{" "}
            <Link color="blue.400" href="https://www.storyden.org/docs">
              Visit the docs
            </Link>{" "}
            for more info if you get stuck.
          </Text>

          {isLoggedIn && (
            <Button onClick={onFinish}>
              {isComplete(3, onboardingStatus) ? "Finish" : "Skip"}
            </Button>
          )}
        </VStack>
      </OrderedList>
    </VStack>
  );
}
