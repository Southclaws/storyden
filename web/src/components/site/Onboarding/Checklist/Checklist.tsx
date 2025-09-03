"use client";

import { OnboardingStatus } from "src/api/openapi-schema";
import { CategoryCreateModal } from "src/components/category/CategoryCreate/CategoryCreateModal";

import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
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
      <Heading size="lg">Welcome to Storyden!</Heading>
      <styled.p p="2" textAlign="center">
        Get your community set up
        <br />
        with the following steps:
      </styled.p>
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
          <styled.p>
            Start by creating an account. The first registration is
            automatically given administrator rights!
          </styled.p>
        </ChecklistItem>

        <ChecklistItem
          step={2}
          current={onboardingStatus}
          title="Create a category"
          onClick={onOpen}
        >
          <styled.p>
            Posts need to live somewhere! So create your first category, give it
            a name and set it up just how you like!
          </styled.p>
          <CategoryCreateModal onClose={onClose} isOpen={isOpen} />
        </ChecklistItem>

        <ChecklistItem
          step={3}
          current={onboardingStatus}
          title="Write your first post"
          url="/new"
        >
          <styled.p>
            An intro, a thesis, a manifesto, a set of rules or just a hi! Get
            started on your first post in your new category!
          </styled.p>
        </ChecklistItem>

        <VStack textAlign="center" px="2">
          <Heading size="md">Invite your people</Heading>
          <styled.p>
            And you&apos;re ready to go! Spread the word and let the posts flow.
          </styled.p>

          <styled.p>
            <LinkButton
              size="xs"
              colorPalette="accent"
              href="https://www.storyden.org/docs"
            >
              Visit the docs
            </LinkButton>{" "}
            for more info if you get stuck.
          </styled.p>

          {isLoggedIn && (
            <Button onClick={onFinish}>
              {isComplete(3, onboardingStatus) ? "Finish" : "Skip"}
            </Button>
          )}
        </VStack>
      </styled.ol>

      <hr />

      <VStack>
        <Heading size="sm">Not an admin?</Heading>

        <p>
          This site is not quite ready to use but you can still browse around!
        </p>

        <Button onClick={onFinish}>Hide this message</Button>
      </VStack>
    </VStack>
  );
}
