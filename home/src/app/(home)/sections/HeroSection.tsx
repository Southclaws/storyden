import { Box, Grid, GridItem, HStack, VStack } from "@/styled-system/jsx";
import Image from "next/image";
import Link from "next/link";

import { cx } from "@/styled-system/css";
import { linkButton } from "@/styled-system/patterns";

type Props = {
  locale?: "en" | "zh";
};

export function HeroSection({ locale = "en" }: Props) {
  const copy =
    locale === "zh"
      ? {
          getStarted: "开始使用",
          liveDemo: "在线演示",
          docsHref: "/zh/docs/introduction",
          heroAlt: "阳光照耀的湖泊，远处是覆雪的高山。",
          logoAlt: "Storyden 标志",
        }
      : {
          getStarted: "Get Started",
          liveDemo: "Live Demo",
          docsHref: "/docs/introduction",
          heroAlt:
            "A sun-lit lake sitting before tall snow-covered mountains in the distance.",
          logoAlt: "The Storyden logo",
        };

  return (
    <Grid>
      <GridItem gridRow="1/2" gridColumn="1/2">
        <picture>
          <source
            media="(max-width: 768px)"
            srcSet="square-nice-lake.webp"
            width={1024}
            height={1024}
          />
          <source
            media="(min-width: 768px)"
            srcSet="wide-nice-lake.webp"
            width={3456}
            height={1728}
          />
          <img
            src="wide-nice-lake.webp"
            role="presentation"
            alt={copy.heroAlt}
            width={3456}
            height={1728}
          />
        </picture>
      </GridItem>

      <GridItem
        gridRow="1/2"
        gridColumn="1/2"
        zIndex={2}
        display="flex"
        justifyContent="center"
        alignItems="center"
        background="linear-gradient(180deg, rgba(59, 83, 111, 0) 49.48%, #000000 97.92%)"
      >
        <VStack gap="6" pt="16">
          <Box width={[36, 40, 40, 48]}>
            <Image
              src="/brand/fullmark_newspaper_vertical_large.png"
              width="1790"
              height="1170"
              alt={copy.logoAlt}
            />
          </Box>

          <HStack gap={4}>
            <Link
              className={linkButton({
                backgroundColor: "white",
                boxShadow: "xl",
              })}
              href={copy.docsHref}
            >
              {copy.getStarted}
            </Link>
            <Link
              target="_blank"
              className={cx(
                "story__text-overlay",
                linkButton({
                  backdropBlur: "lg",
                  backdropFilter: "auto",
                  backgroundColor: "rgba(98, 98, 98, 0.5)",
                  boxShadow: "xl",
                  color: "white",
                  _hover: {
                    color: "black",
                    background: "white",
                  },
                })
              )}
              href="https://makeroom.club"
            >
              {copy.liveDemo}
            </Link>
          </HStack>
        </VStack>
      </GridItem>
    </Grid>
  );
}
