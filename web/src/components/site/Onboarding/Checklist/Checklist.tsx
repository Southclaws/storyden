"use client";

import { OnboardingStatus } from "@/api/openapi-schema";
import { CategoryCreateModal } from "@/components/category/CategoryCreate/CategoryCreateModal";
import { Button } from "@/components/ui/button";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { VStack, styled } from "@/styled-system/jsx";

import { ChecklistItem } from "./ChecklistItem";
import { isComplete, useChecklist } from "./useChecklist";

type Props = {
  onboardingStatus: OnboardingStatus;
  onFinish: () => void;
};

export function Checklist({ onboardingStatus, onFinish }: Props) {
  const { onOpen, isOpen, onClose, isLoggedIn } = useChecklist();

  return (
    <VStack width="full" height="full" justify="start" pt="4" pb="16">
      <PageHeading>Welcome to Storyden!</PageHeading>
      <Text variant="supporting" p="2" textAlign="center">
        Get your community set up
        <br />
        with the following steps:
      </Text>
      <styled.ol
        display="flex"
        flexDir="column"
        gap="4"
        listStyleType="none"
        m="0"
      >
        <ChecklistItem
          step={1}
          current={onboardingStatus}
          title="Create an account"
          url="/register"
        >
          <Text variant="supporting">
            Start by creating an account. The first registration is
            automatically given administrator rights!
          </Text>
        </ChecklistItem>

        <ChecklistItem
          step={2}
          current={onboardingStatus}
          title="Create a category"
          onClick={onOpen}
        >
          <Text variant="supporting">
            Posts need to live somewhere! So create your first category, give it
            a name and set it up just how you like!
          </Text>
          <CategoryCreateModal onClose={onClose} isOpen={isOpen} />
        </ChecklistItem>

        <ChecklistItem
          step={3}
          current={onboardingStatus}
          title="Write your first post"
          url="/new"
        >
          <Text variant="supporting">
            An intro, a thesis, a manifesto, a set of rules or just a hi! Get
            started on your first post in your new category!
          </Text>
        </ChecklistItem>

        <VStack textAlign="center" px="2">
          <SectionHeading>Invite your people</SectionHeading>
          <Text variant="supporting">
            And you&apos;re ready to go! Spread the word and let the posts flow.
          </Text>

          <Text variant="supporting">
            <LinkButton
              colorPalette="accent"
              href="https://www.storyden.org/docs"
            >
              Visit the docs
            </LinkButton>{" "}
            for more info if you get stuck.
          </Text>

          {isLoggedIn && (
            <Button onClick={onFinish}>
              {isComplete(3, onboardingStatus) ? "Finish" : "Skip"}
            </Button>
          )}
        </VStack>
      </styled.ol>

      <hr />

      <VStack>
        <SectionHeading>Not an admin?</SectionHeading>

        <p>
          This site is not quite ready to use but you can still browse around!
        </p>

        <Button onClick={onFinish}>Hide this message</Button>
      </VStack>
    </VStack>
  );
}
