import { Divider, HStack, VStack, styled } from "@/styled-system/jsx";
import { token } from "@/styled-system/tokens";
import Image from "next/image";

export function Footer() {
  return (
    <VStack
      bgColor="Shades.newspaper"
      px={{ base: 6, md: 12, lg: 24, xl: 48, "2xl": 80 }}
      py={{ base: 8, lg: 12 }}
      gap={4}
      flex="1"
    >
      <HStack gap={[2, 3, 4]}>
        <VStack alignItems="end" fontSize="sm" color="Primary.moonlit">
          <styled.a href="https://discord.gg/XF6ZBGF9XF" className="link">
            <HStack gap="2">
              <styled.p>Discord</styled.p>
              <Discord width="1.2em" />
            </HStack>
          </styled.a>

          <styled.a
            href="https://github.com/Southclaws/storyden"
            className="link"
          >
            <HStack gap="2">
              <styled.p>GitHub</styled.p>
              <GitHub width="1.2em" />
            </HStack>
          </styled.a>

          <styled.a href="https://twitter.com/Southclaws" className="link">
            <HStack gap="2">
              <styled.p>X</styled.p>
              <Twitter width="1.2em" />
            </HStack>
          </styled.a>
        </VStack>

        <Divider orientation="vertical" />

        <Image
          src="/brand/fullmark_ink_horizontal.png"
          alt="The Storyden logomark and wordmark"
          width={150 * 1.5}
          height={35 * 1.5}
        />
      </HStack>

      <styled.p textAlign="center" fontSize="xs">
        Storyden&nbsp;brand,&nbsp;logo
        <> and other assets </>
        &copy;&nbsp;Barnaby&nbsp;Keene
      </styled.p>
    </VStack>
  );
}

function Discord(props: any) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 127.14 96.36"
      // fill={token("colors.Primary.moonlit")}
      fill="currentcolor"
      {...props}
    >
      <g data-name="\u56FE\u5C42 2">
        <g data-name="Discord Logos">
          <path
            d="M107.7 8.07A105.15 105.15 0 0 0 81.47 0a72.06 72.06 0 0 0-3.36 6.83 97.68 97.68 0 0 0-29.11 0A72.37 72.37 0 0 0 45.64 0a105.89 105.89 0 0 0-26.25 8.09C2.79 32.65-1.71 56.6.54 80.21a105.73 105.73 0 0 0 32.17 16.15 77.7 77.7 0 0 0 6.89-11.11 68.42 68.42 0 0 1-10.85-5.18c.91-.66 1.8-1.34 2.66-2a75.57 75.57 0 0 0 64.32 0c.87.71 1.76 1.39 2.66 2a68.68 68.68 0 0 1-10.87 5.19 77 77 0 0 0 6.89 11.1 105.25 105.25 0 0 0 32.19-16.14c2.64-27.38-4.51-51.11-18.9-72.15ZM42.45 65.69C36.18 65.69 31 60 31 53s5-12.74 11.43-12.74S54 46 53.89 53s-5.05 12.69-11.44 12.69Zm42.24 0C78.41 65.69 73.25 60 73.25 53s5-12.74 11.44-12.74S96.23 46 96.12 53s-5.04 12.69-11.43 12.69Z"
            data-name="Discord Logo - Large - White"
          />
        </g>
      </g>
    </svg>
  );
}

function Twitter(props: any) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 248 204"
      // fill={token("colors.Primary.moonlit")}
      fill="currentcolor"
      style={{
        enableBackground: "new 0 0 248 204",
      }}
      {...props}
    >
      <path d="M221.95 51.29c.15 2.17.15 4.34.15 6.53 0 66.73-50.8 143.69-143.69 143.69v-.04c-27.44.04-54.31-7.82-77.41-22.64 3.99.48 8 .72 12.02.73 22.74.02 44.83-7.61 62.72-21.66-21.61-.41-40.56-14.5-47.18-35.07a50.338 50.338 0 0 0 22.8-.87C27.8 117.2 10.85 96.5 10.85 72.46v-.64a50.18 50.18 0 0 0 22.92 6.32C11.58 63.31 4.74 33.79 18.14 10.71a143.333 143.333 0 0 0 104.08 52.76 50.532 50.532 0 0 1 14.61-48.25c20.34-19.12 52.33-18.14 71.45 2.19 11.31-2.23 22.15-6.38 32.07-12.26a50.69 50.69 0 0 1-22.2 27.93c10.01-1.18 19.79-3.86 29-7.95a102.594 102.594 0 0 1-25.2 26.16z" />
    </svg>
  );
}

function GitHub(props: any) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 98 96"
      // fill={token("colors.Primary.moonlit")}
      fill="currentcolor"
      {...props}
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0112.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z"
      />
    </svg>
  );
}
