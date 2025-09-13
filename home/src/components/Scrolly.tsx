"use client";

// @ts-ignore - react-scrollama doesn't have types
import { Scrollama, Step } from "react-scrollama";
import { useState } from "react";
import { Box, VStack, Center, styled, Grid } from "@/styled-system/jsx";
import { css } from "@/styled-system/css";
import Image from "next/image";
import { motion, AnimatePresence } from "framer-motion";

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
              <Center
                w="full"
                h="full"
                p="4"
                pb={{
                  base: "24",
                  md: "4",
                }}
                alignItems={{
                  base: "end",
                  md: "center",
                }}
              >
                <Image
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
              <Center w="full" h="full" rounded="2xl" p="4">
                <Image
                  className={css({
                    // NOTE: Padding because image has drop shadow on bot/right.
                    paddingLeft: "2",
                    paddingTop: "2",
                  })}
                  src="/organise.png"
                  alt="A screenshot of Storyden's organization feature showing how to structure content in libraries."
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
              <Center w="full" h="full" rounded="2xl" p="4">
                <Image
                  className={css({
                    // NOTE: Padding because image has drop shadow on bot/right.
                    paddingLeft: "2",
                    paddingTop: "2",
                  })}
                  src="/curate.png"
                  alt="A screenshot of Storyden's curation features for organizing community knowledge."
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
        offset={0.3}
        debug
      >
        <Step scrollamaId="smart" id="step-smart" data={0}>
          <Center
            id="step-inner-smart"
            className={sectionLeftLayoutStyles}
            bgColor="Primary.moonlit"
            color="Shades.newspaper"
            gridRow="1/2"
            alignItems={{
              base: "start",
              md: "center",
            }}
          >
            <VStack alignItems="start" p="4">
              <styled.h1
                fontSize={{
                  base: "3xl",
                  md: "4xl",
                  lg: "5xl",
                }}
                fontWeight="bold"
                textWrap="nowrap"
              >
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
            className={sectionLeftLayoutStyles}
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
              <styled.h1
                fontSize={{
                  base: "3xl",
                  md: "4xl",
                  lg: "5xl",
                }}
                fontWeight="bold"
              >
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
            className={sectionLeftLayoutStyles}
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
              <styled.h1 fontSize="6xl" fontWeight="bold">
                Curate, effortlessly.
              </styled.h1>
              <styled.p fontSize="xl" lineHeight="relaxed">
                It doesn’t stop at threads. Let your community’s Library grow
                with structured and organised pages.
              </styled.p>
              <styled.p fontSize="xl" lineHeight="relaxed">
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
