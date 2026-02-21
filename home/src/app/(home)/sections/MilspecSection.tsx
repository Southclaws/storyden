import Link from "next/link";

import { css } from "@/styled-system/css";
import { Center, VStack, styled } from "@/styled-system/jsx";
import { DockerCopyButton } from "@/components/DockerCopyButton";
import { StorydenComputer } from "@/components/StorydenComputer";

export type HomeStats = {
  stars: number;
  commits: number;
  contributors: number;
  loc: number;
  apis: number;
};

const cellStyle = { border: "1px solid currentColor", padding: "4px" };

const cellFonts = css({
  fontSize: {
    base: "2xs",
    sm: "xs",
    md: "sm",
    lg: "md",
  },
});

type Props = {
  stats: HomeStats;
};

export function MilspecSection({ stats }: Props) {
  return (
    <VStack
      bgColor="Mono.ink"
      py={{
        base: "4",
        sm: "12",
        md: "16",
        lg: "20",
      }}
      px={{
        base: "4",
        sm: "8",
        md: "12",
        lg: "16",
      }}
      gap={{
        base: "4",
        md: "8",
      }}
      fontFamily="gorton"
      letterSpacing="widest"
    >
      <styled.table
        w="full"
        maxW="breakpoint-lg"
        borderColor="Shades.newspaper"
        borderStyle="solid"
        borderWidth="thin"
        color="Shades.newspaper"
        fontSize="xs"
      >
        <tbody>
          <styled.tr>
            <td colSpan={6}>
              <styled.p
                aria-hidden
                fontFamily="gorton"
                fontSize="2xs"
                textAlign="end"
                p="2"
              >
                CHART 1 of 2
              </styled.p>
            </td>
          </styled.tr>

          <styled.tr>
            <td colSpan={6}>
              <Center w="full" py="8">
                <styled.h2 fontFamily="gorton" fontSize="lg" textAlign="center">
                  STORYDEN&nbsp;HUMAN&nbsp;COMPUTER
                  <br />
                  KNOWLEDGE&nbsp;SYSTEM
                </styled.h2>
              </Center>
            </td>
          </styled.tr>

          <styled.tr>
            <td colSpan={6}>
              <StorydenComputer />
            </td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <td style={cellStyle} colSpan={2}>
              GITHUB STARS
            </td>
            <td style={cellStyle}>{stats.stars}</td>
            <td style={cellStyle} rowSpan={3} colSpan={2}>
              SUPPORTED
              <br />
              OPERATING
              <br />
              SYSTEMS
            </td>
            <td style={cellStyle} rowSpan={3} colSpan={2}>
              WINDOWS
              <br />
              MACOS
              <br />
              LINUX
            </td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <td style={cellStyle} colSpan={2}>
              GIT COMMITS
            </td>
            <td style={cellStyle}>{stats.commits}</td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <td style={cellStyle} colSpan={2}>
              API ENDPOINTS
            </td>
            <td style={cellStyle}>89</td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <td style={cellStyle} colSpan={2}>
              CONTRIBUTORS
            </td>
            <td style={cellStyle}>{stats.contributors}</td>
            <td style={cellStyle} colSpan={2}>
              MIN MEMORY
            </td>
            <td style={cellStyle}>100 MB</td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <td style={cellStyle} colSpan={2}>
              LINES OF CODE
            </td>
            <td style={cellStyle}>{stats.loc}</td>
            <td style={cellStyle} colSpan={2}>
              MIN CORES
            </td>
            <td style={cellStyle}>1 CPU</td>
          </styled.tr>
        </tbody>
      </styled.table>

      <styled.table
        w="full"
        maxW="breakpoint-lg"
        borderColor="Shades.newspaper"
        borderStyle="solid"
        borderWidth="thin"
        color="Shades.newspaper"
        fontSize="xs"
      >
        <tbody>
          <styled.tr>
            <td colSpan={5}>
              <styled.p
                aria-hidden
                fontFamily="gorton"
                fontSize="2xs"
                textAlign="end"
                p="2"
              >
                CHART 2 of 2
              </styled.p>
            </td>
          </styled.tr>

          <styled.tr>
            <td colSpan={5}>
              <Center w="full" pt="4" pb="8">
                <styled.h2 fontFamily="gorton" fontSize="lg" textAlign="center">
                  Up&nbsp;and&nbsp;running&nbsp;before
                  <wbr />
                  your&nbsp;coffee&nbsp;gets&nbsp;cold
                </styled.h2>
              </Center>
            </td>
          </styled.tr>

          <styled.tr
            display={{
              base: "none",
              md: "table-row",
            }}
          >
            <td colSpan={5}>
              <Center pb="8">
                <DockerCopyButton />
              </Center>
            </td>
          </styled.tr>

          <styled.tr>
            <styled.td style={cellStyle} colSpan={2}>
              <Link href="/docs/introduction">
                <styled.p p="2" textAlign="center">
                  PLEASE SEE SUPPLIED MANUAL FOR OPERATION INSTRUCTIONS
                </styled.p>
              </Link>
            </styled.td>
            <styled.td style={cellStyle} colSpan={3} textWrap="balance">
              <Link
                href="https://github.com/Southclaws/storyden"
                target="_blank"
              >
                <styled.p p="2" textAlign="center">
                  OPEN SOURCE SOFTWARE RELEASED TO THE PUBLIC UNDER THE MIT
                  LICENSE
                </styled.p>
              </Link>
            </styled.td>
          </styled.tr>
        </tbody>
      </styled.table>
    </VStack>
  );
}
