"use client";

import { css } from "@/styled-system/css";
import { Box, Center, Grid, styled, VStack } from "@/styled-system/jsx";
import { AnimatePresence, motion } from "framer-motion";
import Image from "next/image";
import { useState } from "react";

// @ts-ignore - react-scrollama doesn't have types
import { Scrollama, Step } from "react-scrollama";

const zBack = 0;
const zHighlight = 1;
const zInfo = 2;

// sections: smart, sort, curate
type SectionID = "smart" | "sort" | "curate";

type ScrollamaEnterEvent = {
  data: number;
  direction: "up" | "down";
  element: Element;
  entry: IntersectionObserverEntry;
  scrollamaId: SectionID;
};

type ScrollamaExitEvent = {
  data: number;
  direction: "up" | "down";
  element: Element;
  entry: IntersectionObserverEntry;
  scrollamaId: SectionID;
};

type ScrollamaProgressEvent = {
  data: number;
  direction: "up" | "down"; // NOTE: Always the same direction as when entered.
  element: Element;
  entry: IntersectionObserverEntry;
  progress: number;
  scrollamaId: SectionID;
};

const gridStyles = css({
  gridTemplateColumns: {
    base: "[left] 100dvw [right]",
    md: "[left] 50% [middle] 50% [right]",
  },
  gridTemplateRows: "[top] 100dvh 100dvh 100dvh [bottom]",
  position: "relative",
});

const sectionLeftLayoutStyles = css({
  gridColumn: {
    base: "left / right",
    md: "left / middle",
  },
  // gridRow: "top / bottom",
  height: {
    base: "min",
    md: "dvh",
  },
  zIndex: zInfo,
});

const sectionMoonlitShadowStyles = css({
  boxShadow: {
    base: "0 32px 32px 8px {colors.Primary.moonlit}",
    md: "none",
  },
});

const sectionForestShadowStyles = css({
  boxShadow: {
    base: "0 32px 32px 8px {colors.Primary.forest}",
    md: "none",
  },
});

const sectionCampfireShadowStyles = css({
  boxShadow: {
    base: "0 32px 32px 8px {colors.Primary.campfire}",
    md: "none",
  },
});

const sectionRightLayoutStyles = css({
  zIndex: zHighlight,
  gridColumn: {
    base: "left / right",
    md: "middle / right",
  },
  gridRow: "top / bottom",
  height: "full",
  position: "sticky",
  top: 0,
  h: "100vh",
  bg: "colors.Mono.slush",
});

const sectionRightImageContainerStyles = css({
  w: "full",
  h: "full",
  p: "4",
  pb: {
    base: "24",
    md: "4",
  },
  alignItems: {
    base: "end",
    md: "center",
  },
});

const sectionHeadingStyles = css({
  fontSize: {
    base: "3xl",
    md: "4xl",
    lg: "5xl",
  },
  fontWeight: "bold",
  textWrap: "balance",
});

const sectionParagraphStyles = css({
  fontSize: {
    base: "md",
    sm: "lg",
    md: "xl",
  },
  lineHeight: "relaxed",
});

export function Scrolly() {
  const [currentStepIndex, setCurrentStepIndex] = useState(0);

  const onStepEnter = (data: ScrollamaEnterEvent) => {
    setCurrentStepIndex(data.data);
  };

  const onStepExit = ({ data }: ScrollamaExitEvent) => {
    // console.log("exit", data);
  };

  function handleStepProgress(data: ScrollamaProgressEvent) {
    // console.log("progress", data.progress, data.direction);
  }

  return (
    <Grid className={gridStyles} gap="0">
      <Box id="step-right-side" className={sectionRightLayoutStyles}>
        <AnimatePresence mode="wait">
          {currentStepIndex === 0 && (
            <motion.div
              key="step-0"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ position: "absolute", inset: 0 }}
            >
              <Center className={sectionRightImageContainerStyles}>
                <Image
                  priority={true}
                  className={css({
                    // NOTE: Padding because image has drop shadow on bot/right.
                    paddingLeft: "2",
                    paddingTop: "2",
                  })}
                  src="/sort-search.png"
                  alt="A screenshot of Storyden's search and ask feature surfacing a member's vague question looking for a specific thread about design books."
                  width={660 * 2}
                  height={432 * 2}
                />
              </Center>
            </motion.div>
          )}
          {currentStepIndex === 1 && (
            <motion.div
              key="step-1"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ position: "absolute", inset: 0 }}
            >
              <Center className={sectionRightImageContainerStyles}>
                <Image
                  priority={true}
                  className={css({
                    // NOTE: Padding because image has drop shadow on bot/right.
                    paddingLeft: "2",
                    paddingTop: "2",
                  })}
                  src="/organise.png"
                  alt="A graphic showing the flow of a question leading to a directory leading to a specific link to a website."
                  width={660 * 2}
                  height={432 * 2}
                />
              </Center>
            </motion.div>
          )}
          {currentStepIndex === 2 && (
            <motion.div
              key="step-2"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ position: "absolute", inset: 0 }}
            >
              <Center className={sectionRightImageContainerStyles}>
                <Image
                  priority={true}
                  className={css({
                    // NOTE: Padding because image has drop shadow on bot/right.
                    paddingLeft: "2",
                    paddingTop: "2",
                  })}
                  src="/curate.png"
                  alt="A diagram of the link fetching flow going from fetching, assigning tags and creating a page to demonstrate the AI assisted curation capabilities."
                  width={660 * 2}
                  height={432 * 2}
                />
              </Center>
            </motion.div>
          )}
        </AnimatePresence>
      </Box>

      <Box
        gridRow="1/2"
        gridColumn="left / right"
        h="dvh"
        zIndex={zBack}
        bgColor="Primary.moonlit"
      />

      <Box
        gridRow="2/3"
        gridColumn="left / right"
        h="dvh"
        zIndex={zBack}
        bgColor="Primary.forest"
      />

      <Box
        gridRow="3/4"
        gridColumn="left / right"
        h="dvh"
        zIndex={zBack}
        bgColor="Primary.campfire"
      />

      <Scrollama
        onStepEnter={onStepEnter}
        onStepExit={onStepExit}
        onStepProgress={handleStepProgress}
        offset={0.5}
      >
        <Step scrollamaId="smart" id="step-smart" data={0}>
          <Center
            id="step-inner-smart"
            className={`${sectionLeftLayoutStyles} ${sectionMoonlitShadowStyles}`}
            bgColor="Primary.moonlit"
            color="Shades.newspaper"
            gridRow="1/2"
            alignItems={{
              base: "start",
              md: "center",
            }}
          >
            <VStack alignItems="start" p="4">
              <styled.h1 className={sectionHeadingStyles}>
                Built-in&nbsp;smarts.
                <br />
                opt-in&nbsp;brains.
              </styled.h1>
              <styled.p fontSize="xl" lineHeight="relaxed" textWrap="balance">
                Language model integration, at the core.
                <br />
                But only if you want it.
              </styled.p>
            </VStack>
          </Center>
        </Step>

        <Step scrollamaId="sort" id="step-sort" data={1}>
          <Center
            id="step-inner-sort"
            className={`${sectionLeftLayoutStyles} ${sectionForestShadowStyles}`}
            bgColor="Primary.forest"
            color="Shades.newspaper"
            gridRow="2/3"
            alignItems={{
              base: "start",
              md: "center",
            }}
          >
            <VStack
              alignItems="flex-start"
              p="4"
              maxW="prose"
              textWrap="balance"
            >
              <styled.h1 className={sectionHeadingStyles}>
                Sort. Search. Ask.
              </styled.h1>
              <styled.p fontSize="xl" lineHeight="relaxed">
                Allow your collective knowledge to grow without losing great
                ideas to the void of time, banished to the archive.
              </styled.p>
              <p>Looking for something?</p>
              <p>Just ask.</p>
            </VStack>
          </Center>
        </Step>

        <Step scrollamaId="curate" id="step-curate" data={2}>
          <Center
            id="step-inner-curate"
            className={`${sectionLeftLayoutStyles} ${sectionCampfireShadowStyles}`}
            bgColor="Primary.campfire"
            color="Primary.moonlit"
            gridRow="3/4"
            alignItems={{
              base: "start",
              md: "center",
            }}
          >
            <VStack
              alignItems="flex-start"
              p="4"
              maxW="prose"
              textWrap="balance"
            >
              <styled.h1 className={sectionHeadingStyles}>
                Curate, effortlessly.
              </styled.h1>
              <styled.p className={sectionParagraphStyles}>
                It doesn’t stop at threads. Let your community’s Library grow
                with structured and organised pages.
              </styled.p>
              <styled.p className={sectionParagraphStyles}>
                Alexandria-scale corpus? No worries, let Storyden’s intelligence
                sort it out.
              </styled.p>
            </VStack>
          </Center>
        </Step>
      </Scrollama>
    </Grid>
  );
}
