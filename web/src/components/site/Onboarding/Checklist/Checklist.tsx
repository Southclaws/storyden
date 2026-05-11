"use client";

import { OnboardingStatus } from "src/api/openapi-schema";
import { CategoryCreateModal } from "src/components/category/CategoryCreate/CategoryCreateModal";

import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { useI18n } from "@/i18n/provider";
import { VStack, styled } from "@/styled-system/jsx";

import { ChecklistItem } from "./ChecklistItem";
import { isComplete, useChecklist } from "./useChecklist";

type Props = {
  onboardingStatus: OnboardingStatus;
  onFinish: () => void;
};

export function Checklist({ onboardingStatus, onFinish }: Props) {
  const { t } = useI18n();
  const { onOpen, isOpen, onClose, isLoggedIn } = useChecklist();

  return (
    <VStack width="full" height="full" justify="start" pt="4" pb="16">
      <Heading size="lg">{t("Welcome to Storyden!")}</Heading>
      <styled.p p="2" textAlign="center">
        {t("Get your community set up")}
        <br />
        {t("with the following steps:")}
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
          title={t("Create an account")}
          url="/register"
        >
          <styled.p>
            {t(
              "Start by creating an account. The first registration is automatically given administrator rights!",
            )}
          </styled.p>
        </ChecklistItem>

        <ChecklistItem
          step={2}
          current={onboardingStatus}
          title={t("Create a category")}
          onClick={onOpen}
        >
          <styled.p>
            {t(
              "Posts need to live somewhere! So create your first category, give it a name and set it up just how you like!",
            )}
          </styled.p>
          <CategoryCreateModal onClose={onClose} isOpen={isOpen} />
        </ChecklistItem>

        <ChecklistItem
          step={3}
          current={onboardingStatus}
          title={t("Write your first post")}
          url="/new"
        >
          <styled.p>
            {t(
              "An intro, a thesis, a manifesto, a set of rules or just a hi! Get started on your first post in your new category!",
            )}
          </styled.p>
        </ChecklistItem>

        <VStack textAlign="center" px="2">
          <Heading size="md">{t("Invite your people")}</Heading>
          <styled.p>
            {t(
              "And you're ready to go! Spread the word and let the posts flow.",
            )}
          </styled.p>

          <styled.p>
            <LinkButton
              size="xs"
              colorPalette="accent"
              href="https://www.storyden.org/docs"
            >
              {t("Visit the docs")}
            </LinkButton>{" "}
            {t("for more info if you get stuck.")}
          </styled.p>

          {isLoggedIn && (
            <Button onClick={onFinish}>
              {isComplete(3, onboardingStatus) ? t("Finish") : t("Skip")}
            </Button>
          )}
        </VStack>
      </styled.ol>

      <hr />

      <VStack>
        <Heading size="sm">{t("Not an admin?")}</Heading>

        <p>
          {t(
            "This site is not quite ready to use but you can still browse around!",
          )}
        </p>

        <Button onClick={onFinish}>{t("Hide this message")}</Button>
      </VStack>
    </VStack>
  );
}
