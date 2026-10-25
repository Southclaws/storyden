import Link from "next/link";

import { DockerCopyButton } from "@/components/DockerCopyButton";
import { StorydenComputer } from "@/components/StorydenComputer";
import { css } from "@/styled-system/css";
import { Center, VStack, styled } from "@/styled-system/jsx";

export type HomeStats = {
  stars: number | null;
  commits: number | null;
  contributors: number | null;
  loc: number | null;
  apis: number | null;
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
      className="milspec-section"
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
        <caption className="milspec-caption">
          <styled.p
            aria-hidden
            fontFamily="gorton"
            fontSize="2xs"
            textAlign="end"
            p="2"
          >
            CHART 1 of 2
          </styled.p>
          <Center w="full" py="8">
            <styled.h2
              id="knowledge-system-title"
              fontFamily="gorton"
              fontSize="lg"
              textAlign="center"
            >
              STORYDEN&nbsp;HUMAN&nbsp;COMPUTER
              <br />
              KNOWLEDGE&nbsp;SYSTEM
            </styled.h2>
          </Center>
          <StorydenComputer />
        </caption>
        <tbody>
          <styled.tr className={cellFonts}>
            <th scope="row" style={cellStyle} colSpan={2}>
              GITHUB STARS
            </th>
            <td style={cellStyle}>
              {stats.stars?.toLocaleString("en-GB") ?? "—"}
            </td>
            <th id="supported-os" style={cellStyle} rowSpan={3} colSpan={2}>
              SUPPORTED
              <br />
              OPERATING
              <br />
              SYSTEMS
            </th>
            <td
              headers="supported-os"
              style={cellStyle}
              rowSpan={3}
              colSpan={2}
            >
              WINDOWS
              <br />
              MACOS
              <br />
              LINUX
            </td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <th scope="row" style={cellStyle} colSpan={2}>
              GIT COMMITS
            </th>
            <td style={cellStyle}>
              {stats.commits?.toLocaleString("en-GB") ?? "—"}
            </td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <th scope="row" style={cellStyle} colSpan={2}>
              API ENDPOINTS
            </th>
            <td style={cellStyle}>
              {stats.apis?.toLocaleString("en-GB") ?? "—"}
            </td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <th scope="row" style={cellStyle} colSpan={2}>
              CONTRIBUTORS
            </th>
            <td style={cellStyle}>
              {stats.contributors?.toLocaleString("en-GB") ?? "—"}
            </td>
            <th scope="row" style={cellStyle} colSpan={2}>
              MIN MEMORY
            </th>
            <td style={cellStyle} colSpan={2}>
              256 MB
            </td>
          </styled.tr>

          <styled.tr className={cellFonts}>
            <th scope="row" style={cellStyle} colSpan={2}>
              LINES OF CODE
            </th>
            <td style={cellStyle}>
              {stats.loc?.toLocaleString("en-GB") ?? "—"}
            </td>
            <th scope="row" style={cellStyle} colSpan={2}>
              MIN CORES
            </th>
            <td style={cellStyle} colSpan={2}>
              2 CPU
            </td>
          </styled.tr>
        </tbody>
      </styled.table>

      <styled.section
        className="deployment-plate"
        aria-labelledby="deployment-title"
        w="full"
        maxW="breakpoint-lg"
        borderColor="Shades.newspaper"
        borderStyle="solid"
        borderWidth="thin"
        color="Shades.newspaper"
        fontSize="xs"
      >
        <styled.p
          aria-hidden
          fontFamily="gorton"
          fontSize="2xs"
          textAlign="end"
          p="2"
        >
          CHART 2 of 2
        </styled.p>
        <Center w="full" pt="4" pb="8">
          <styled.h2
            id="deployment-title"
            fontFamily="gorton"
            fontSize="lg"
            textAlign="center"
          >
            Up&nbsp;and&nbsp;running&nbsp;before <wbr />
            your&nbsp;coffee&nbsp;gets&nbsp;cold
          </styled.h2>
        </Center>
        <Center pb="8">
          <DockerCopyButton />
        </Center>
        <div className="deployment-links">
          <Link href="/docs/introduction">
            <styled.p p="2" textAlign="center">
              PLEASE SEE SUPPLIED MANUAL FOR OPERATION INSTRUCTIONS
            </styled.p>
          </Link>
          <Link href="https://github.com/Southclaws/storyden" target="_blank">
            <styled.p p="2" textAlign="center">
              OPEN SOURCE SOFTWARE
            </styled.p>
          </Link>
        </div>
      </styled.section>
    </VStack>
  );
}
