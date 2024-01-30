import {
  Box,
  Divider,
  Flex,
  Grid,
  GridItem,
  HStack,
  VStack,
  styled,
} from "@/styled-system/jsx";
import { NextSeo } from "next-seo";
import Image from "next/image";
import Link from "next/link";

import { seo } from "../seo";
import { FullLogo } from "@/components/FullLogo";
import { css, cx } from "@/styled-system/css";
import { Floating, linkButton } from "@/styled-system/patterns";
import { Logo } from "@/components/Logo";

function Nav() {
  const item = css({
    px: "4",
    py: "2",
    flexShrink: "0",
    borderRadius: "md",
    backgroundColor: "unset",
    _hover: {
      backgroundColor: "bg.opaque",
      shadow: "sm",
    },
  });

  return (
    <styled.nav
      position="fixed"
      display="flex"
      justifyContent="center"
      pt="8"
      w="full"
      zIndex="9"
    >
      <styled.ul
        className={Floating()}
        display="flex"
        flexShrink="0"
        justifyContent="center"
        gap="4"
        height="full"
        color="fg.default"
        fontWeight="bold"
        p="1"
      >
        <styled.li display="flex" alignItems="center" px="2">
          <Link href="/">
            <Logo width="2rem" />
          </Link>
        </styled.li>
        <styled.li display="flex">
          <Link className={item} href="/docs">
            Docs
          </Link>
        </styled.li>
        <styled.li display="flex">
          <Link className={item} href="/blog">
            Blog
          </Link>
        </styled.li>
        <styled.li display="flex">
          <Link className={item} href="https://github.com/Southclaws/storyden">
            Source
          </Link>
        </styled.li>
      </styled.ul>
    </styled.nav>
  );
}

function Hero() {
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
            alt="A sun-lit lake sitting before tall snow-covered mountains in the distance."
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
        <VStack gap="0" pt="16">
          <Box width={[36, 40, 40, 48]}>
            <FullLogo />
          </Box>
          <Grid>
            <GridItem
              gridRow="1/2"
              gridColumn="1/2"
              width={["xs", "sm", "2xl"]}
            >
              <Image
                src="/brushything.webp"
                width="716"
                height="214"
                alt=""
                role="presentation"
              />
            </GridItem>

            <GridItem
              gridRow="1/2"
              gridColumn="1/2"
              zIndex={2}
              display="flex"
              justifyContent="center"
              alignItems="center"
            >
              <styled.h1 fontSize={["sm", "md", "2xl"]}>
                A forum for the modern age
              </styled.h1>
            </GridItem>
          </Grid>

          <HStack gap={4}>
            <Link
              className={linkButton({
                backgroundColor: "white",
              })}
              href="/docs"
            >
              Get started
            </Link>
            <Link
              className={cx(
                "story__text-overlay",
                linkButton({
                  backdropBlur: "lg",
                  backdropFilter: "auto",
                  backgroundColor: "rgba(98, 98, 98, 0.5)",
                  color: "white",
                  _hover: {
                    color: "black",
                    background: "white",
                  },
                })
              )}
              href="https://github.com/Southclaws/storyden"
            >
              Source code
            </Link>
          </HStack>
        </VStack>
      </GridItem>
    </Grid>
  );
}

function Story() {
  return (
    <Grid
      maxW="100vw"
      bgColor="black"
      //
      // The grid template rows and columns are defined in order to produce an
      // overlapping effect so that the end of one overlaps with the next start.
      gridTemplateRows={`48px [top] auto [one] 0.4fr [one-text-start] 0.2fr [one-text-end] 0.4fr [two] 90px [one-end] auto [two-text-start] auto [two-text-end] auto [three] 90px [two-end] 0.4fr [three-text-start] 0.2fr [three-text-end] 0.4fr 90px [three-end] auto [bot] 1px`}
      gridTemplateColumns={{
        base: `0em [left] 1fr [left-far] 2fr [left-near] 50% [right-near] 2fr [right-far] 1fr [right] 1em`,
        lg: `10% [left] 2fr [left-far] 2fr [left-near] 25% [right-near] 2fr [right-far] 2fr [right] 10%`,
        xl: `25% [left] 2fr [left-far] 2fr [left-near] 25% [right-near] 2fr [right-far] 2fr [right] 25%`,
      }}
    >
      <GridItem maxW="100vw" gridRow="top / bot" gridColumn="left / right">
        <VStack
          width="full"
          height="full"
          alignItems="center"
          justifyContent="center"
        >
          <picture>
            <source media="(max-width: 768px)" srcSet="tall-bg-stars.webp" />
            <source media="(min-width: 768px)" srcSet="bg-stars.webp" />
            <img
              src="bg-stars.webp"
              alt=""
              role="presentation"
              width={2048}
              height={2048}
            />
          </picture>
        </VStack>
      </GridItem>

      <GridItem gridRow="one / one-end" gridColumn="left-near / right-far">
        <HStack justify="center">
          <Box>
            <Image
              src="/squircle-tree-smol.webp"
              width={512}
              height={512}
              alt="A baby tree sapling about a foot tall"
            />
          </Box>
        </HStack>
      </GridItem>

      <GridItem
        gridRow="one-text-start / one-text-end"
        gridColumn="left-far / right-near"
        zIndex={2}
      >
        <HStack>
          <styled.p
            className="story__text-overlay"
            boxShadow="2px 2px 10px rgba(0, 0, 0, 0.1)"
            backdropFilter="blur(16px)"
            borderRadius="1em"
            width="min-content"
            p={{ base: 2, lg: 3, xl: 4 }}
            fontSize={{ base: "md", sm: "2xl" }}
            fontWeight="medium"
            color="white"
            wordBreak="keep-all"
          >
            Ideas,&nbsp;big&nbsp;or&nbsp;small,&nbsp;always
            <wbr />
            &nbsp;start&nbsp;with&nbsp;people&nbsp;in&nbsp;mind.
          </styled.p>
        </HStack>
      </GridItem>

      <GridItem gridRow="two / two-end" gridColumn="left-far / right-near">
        <HStack justify="center">
          <Box>
            <Image
              src="/squircle-tree-midaf.webp"
              width={512}
              height={512}
              alt="A young growing tree about 5ft tall"
            />
          </Box>
        </HStack>
      </GridItem>

      <GridItem
        gridRow="two-text-start / two-text-end"
        gridColumn="left-near / right-far"
        zIndex={2}
      >
        <HStack justify="end">
          <styled.p
            className="story__text-overlay"
            boxShadow="2px 2px 10px rgba(0, 0, 0, 0.1)"
            backdropFilter="blur(16px)"
            borderRadius="1em"
            width="min-content"
            p={{ base: 2, lg: 3, xl: 4 }}
            fontSize={{ base: "md", sm: "2xl" }}
            fontWeight="medium"
            color="white"
            wordBreak="keep-all"
            textAlign="center"
          >
            Projects,&nbsp;products&nbsp;and
            <wbr />
            &nbsp; people&nbsp;oriented&nbsp;ideas
            <wbr />
            &nbsp; grow&nbsp;into&nbsp;communities.
          </styled.p>
        </HStack>
      </GridItem>

      <GridItem gridRow="three / three-end" gridColumn="left-near / right-far">
        <HStack justify="center">
          <Box>
            <Image
              src="/squircle-tree-bigly.webp"
              width={512}
              height={512}
              alt="A huge magnificent tree in a forest clearing bathed in sunlight"
            />
          </Box>
        </HStack>
      </GridItem>

      <GridItem
        gridRow="three-text-start / three-text-end"
        gridColumn="left-far / right-near"
        zIndex={2}
      >
        <HStack>
          <styled.p
            className="story__text-overlay"
            boxShadow="2px 2px 10px rgba(0, 0, 0, 0.1)"
            backdropFilter="blur(16px)"
            borderRadius="1em"
            width="min-content"
            p={{ base: 2, lg: 3, xl: 4 }}
            fontSize={{ base: "md", sm: "2xl" }}
            fontWeight="medium"
            color="white"
            wordBreak="keep-all"
          >
            Collaboration&nbsp;occurs&nbsp;best&nbsp;when
            <wbr />
            &nbsp;the&nbsp;platform&nbsp;flows&nbsp;with&nbsp;everyone.
          </styled.p>
        </HStack>
      </GridItem>
    </Grid>
  );
}

function Why() {
  return (
    <Flex
      flexDir={{
        base: "column",
        lg: "row",
      }}
      bgColor="hsla(140, 16%, 88%, 1)"
      px={{ base: 12, md: 48, lg: 48, xl: 80, "2xl": 96 }}
      py={{ base: 12, lg: 12 }}
      alignItems="center"
      justifyContent="center"
      gap={4}
      flex="1"
    >
      <Box>
        <styled.h1
          textAlign="right"
          fontFamily="mona"
          fontStyle="normal"
          fontWeight="900"
          fontSize={{
            base: "4xl",
            md: "7xl",
            lg: "6xl",
            xl: "7xl",
            "2xl": "8xl",
          }}
          width={{ base: "max-content", lg: "min-content" }}
          lineHeight={{ base: 1, lg: 1.4 }}
        >
          Why
          <wbr /> Storyden
        </styled.h1>
      </Box>

      <Flex flexDir="column" gap={{ base: 2, lg: 1 }} pt="10.1px">
        <styled.p>Storyden is a platform for building communities.</styled.p>
        <styled.p wordBreak="keep-all">
          But not just another chat app or another forum site.
          <wbr /> Storyden is a modern take on oldschool bulletin board
          <wbr /> forums you may remember from the earlier days of the
          <wbr /> internet.
        </styled.p>
        <styled.p>
          With a fresh perspective and new objectives, Storyden is
          <wbr /> designed to be the community platform for the next era of
          <wbr /> internet culture.
        </styled.p>
      </Flex>
    </Flex>
  );
}

function Features() {
  return (
    <Flex
      w="full"
      flexWrap="wrap"
      alignItems="center"
      justifyContent="center"
      bgColor="hsla(140, 16%, 88%, 1)"
      py={8}
      px={{ base: 4, sm: 12, md: 16, xl: 24 }}
      gap={12}
    >
      <Feature
        image="/accessible.webp"
        alt=""
        heading="Accessible"
        body="Accessibility is non-negotiable and no one can be left behind. WAI and WCAG are a primary focus to ensure great experience for people regardless of a disability."
      />

      <Feature
        image="/secure.webp"
        alt=""
        heading="Secure"
        body="The latest and greatest industry standard security practices as well as new emerging systems such as WebAuthn guarantee the most secure experience for everyone."
      />

      <Feature
        image="/web3.webp"
        alt=""
        heading="Web3"
        body="Love it or hate it, it’s here and it’s staying. So we embrace the new web and provide features such as wallet based login, NFT avatars and more for web3 communities."
      />

      <Feature
        image="/opensource.webp"
        alt=""
        heading="Open source"
        body="The benefits of open source software are impossible to ignore. When it comes to the security, development velocity, and ability to report issues, this is the way forward."
      />

      <Feature
        image="/extensible.webp"
        alt=""
        heading="Extensible"
        body="A fully documented OpenAPI schema means that you can extend the platform with plugins or even build a whole new frontend from scratch if you want to!"
      />

      <Feature
        image="/builttolast.webp"
        alt=""
        heading="Built to last"
        body="Harnessing the power of technology that’s just-modern-enough helps balance stability with longevity. Storyden uses a carefully chosen toolbox with this in mind."
      />
    </Flex>
  );
}

function Feature({ image, alt, heading, body, ...rest }) {
  return (
    <Grid width="full" maxW={394} {...rest}>
      <GridItem
        gridRow="1/2"
        gridColumn="2/3"
        filter="drop-shadow(5px 5px 10px rgba(0, 0, 0, 0.2))"
        zIndex={0}
        bgColor="hsla(0, 0%, 19%, 1)"
        height="90%"
        width="90%"
        m="auto"
      />

      <GridItem
        gridRow="1/2"
        gridColumn="2/3"
        filter="drop-shadow(5px 5px 10px rgba(0, 0, 0, 0.2))"
        zIndex={1}
      >
        <Image src={image} width={394} height={394} alt={alt} />
      </GridItem>

      <GridItem gridRow="1/2" gridColumn="2/3" p={8} zIndex={2}>
        <Flex flexDir="column" justifyContent="end" height="full" color="white">
          <styled.h1 textShadow="4px 4px 8px #000000">{heading}</styled.h1>
          <styled.p textShadow="2px 2px 4px #000000">{body}</styled.p>
        </Flex>
      </GridItem>
    </Grid>
  );
}

function FeatureHeading({ children, ...props }) {
  return (
    <styled.h1 fontSize="lg" {...props}>
      {children}
    </styled.h1>
  );
}

function ForCommunityLeaders() {
  return (
    <VStack
      bgColor="hsla(140, 16%, 88%, 1)"
      py={8}
      px={{ base: 4, sm: 12, md: 16, lg: 48, xl: 96 }}
      gap={12}
    >
      <Pair heading="community leaders" headingColour="#808080">
        <styled.p>
          Fearless <b>futurism</b>, radical <b>accessibility</b>, endless
          extensibility. Every modern service, product and movement has
          community at the centre. Communities often grow out of their humble
          beginnings on walled-garden platforms. In an era of growing awareness
          of personal <b>privacy</b>, tech <b>monopoly</b> and{" "}
          <b>decentralisation</b>, communities of all sizes are affected.
        </styled.p>
      </Pair>

      <VStack maxW="container.lg" gap={8}>
        <HStack maxW={{ base: "full", sm: "container.md" }}>
          <Box p={{ base: 0, sm: 4 }}>
            <svg
              width="100%"
              height="100%"
              viewBox="0 0 156 156"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M102.219 104.512C102.219 103.107 103.381 101.945 104.786 101.945C106.191 101.945 107.353 103.059 107.353 104.512V109.646C107.353 111.051 106.191 112.213 104.786 112.213C103.381 112.213 102.219 111.051 102.219 109.646V104.512Z"
                fill="#303030"
              />
              <path
                fillRule="evenodd"
                clipRule="evenodd"
                d="M78 5.06934C37.8731 5.06934 5.34375 37.5986 5.34375 77.7256C5.34375 117.853 37.8731 150.382 78 150.382C118.127 150.382 150.656 117.853 150.656 77.7256C150.656 37.5986 118.127 5.06934 78 5.06934ZM58.7476 17.7544C33.3875 25.8892 15.0312 49.6637 15.0312 77.7256C15.0312 112.502 43.2233 140.694 78 140.694C112.777 140.694 140.969 112.502 140.969 77.7256C140.969 72.8781 140.421 68.1586 139.384 63.6256L132.831 73.4633C132.25 74.3352 131.281 74.868 130.216 74.868H128.569C127.745 74.868 127.067 74.1899 127.067 73.3665V70.7509C127.067 64.9384 123.192 55.493 117.38 55.493H109.872C107.305 55.493 106.336 58.6899 107.837 59.9977L108.952 60.918C111.083 62.7102 110.647 66.0524 108.177 67.3118L102.606 70.1212H100.911C99.1672 70.1212 97.7141 68.668 97.7141 66.9243V62.8071C97.7141 61.4993 96.6484 60.3852 95.2922 60.3852L95.2637 60.3852C93.8592 60.3849 88.3172 60.3836 88.3172 67.1665V69.5883C88.3172 74.4805 96.7453 78.404 101.541 79.179C102.316 79.2759 102.945 79.954 102.945 80.729V82.4727C102.945 83.5383 102.703 84.604 102.219 85.5727L92.9672 99.3774H91.3203C89.7219 99.3774 88.4141 100.637 88.3656 102.235L88.075 113.957L79.8406 122.191C78.8719 123.209 77.5641 123.741 76.1594 123.741H73.7859C70.8797 123.741 68.5547 121.368 68.5547 118.51V107.612C68.5547 104.899 67.15 104.221 66.1328 104.221H64.4375C62.6938 104.221 61.2406 102.768 61.2406 101.024V92.693C61.2406 88.1883 57.6078 84.5555 53.1031 84.5555H39.25C34.9391 84.5555 29.5625 76.3696 29.5625 72.0587V68.6196C29.5625 66.4399 30.3859 64.3571 31.8391 62.8071L37.0219 57.3337C38.1844 56.1227 39.8313 55.4446 41.5266 55.493H46.6609C48.1625 55.493 49.5672 54.9118 50.6328 53.8462L52.1828 52.2962C53.2484 51.2305 54.7016 50.6493 56.2031 50.6493H61.2406C62.3547 50.6493 63.6141 51.2305 63.5656 52.5384V53.8946C63.4203 56.704 68.8937 55.4446 68.8937 55.4446C68.8937 55.4446 73.5438 51.2305 75.675 50.7462C77.8547 50.2134 82.2141 51.4243 84.7813 52.9743C86.1859 53.7977 87.9297 52.829 87.9297 51.2305C87.9297 48.179 85.4109 45.7087 82.3594 45.7087H80.5672C79.2594 45.7087 78.1937 44.643 78.1937 43.3352V39.8477C78.1937 39.0727 77.5641 38.443 76.7891 38.443H74.9969C73.9797 38.443 73.0109 38.9274 72.4297 39.7993L69.4266 44.2555C68.8453 45.1758 67.8281 45.7087 66.7141 45.7087H66.0359C64.4859 45.7087 63.0328 45.1274 61.9188 44.0133L59.7875 42.3665C58.0437 40.9133 55.8156 40.8649 55.3797 40.8649H50.6812C49.7125 40.8649 48.9859 41.6399 48.9859 42.5602L48.7922 44.0133C48.7922 45.7571 47.5812 48.1305 45.8375 48.1305H43.0281C41.2844 48.1305 39.8797 46.7259 39.8797 45.0305V39.218C39.8797 37.4259 41.3328 36.0212 43.0766 36.0212C44.0453 36.0212 44.7719 35.2462 44.7719 34.3258V32.8243C44.7719 28.9977 47.5328 28.4165 49.5672 28.9008C50.4323 29.1265 51.3339 29.5129 52.2466 29.9039C53.6844 30.52 55.1496 31.1478 56.5422 31.1774H64.1953C67.0531 31.1774 68.9422 29.2399 68.9422 26.4305C68.9422 25.7288 68.9764 24.9978 69.0101 24.2763C69.1605 21.0633 69.302 18.0388 66.375 18.6321L65.0187 19.9883C64.1469 20.9087 63.6141 22.1196 63.6141 23.4274V24.1055C63.6141 25.3165 62.6453 26.3821 61.4344 26.479C60.2719 26.5274 58.7703 26.3337 58.7703 24.8321C58.7703 23.5754 58.7329 22.8427 58.6992 22.1837C58.6473 21.1663 58.6043 20.3248 58.7219 18.0024C58.7265 17.9192 58.7351 17.8364 58.7476 17.7544ZM88.8222 15.6833C86.3875 15.2616 83.9048 14.9797 81.3828 14.8462L83.0859 16.5493H87.9781L88.8222 15.6833Z"
                fill="#303030"
              />
            </svg>
          </Box>

          <VStack alignItems="start">
            <FeatureHeading size="lg">Flexible by nature</FeatureHeading>
            <styled.p>
              Whether you operate a product <b>support forum</b>, a{" "}
              <b>gaming community</b> or the next big cryptocurrency, web3 or
              DAO project, the <b>people</b> at the centre of whatever you’re
              doing deserve a platform that fades into the background and brings{" "}
              <b>what matters</b> front & centre
            </styled.p>
          </VStack>
        </HStack>

        <HStack maxW={{ base: "full", sm: "container.md" }}>
          <VStack alignItems="start">
            <FeatureHeading>Accessible by design</FeatureHeading>
            <styled.p>
              And by people, that means{" "}
              <em>
                <b>all</b>
              </em>{" "}
              people. A truly welcoming, inclusive and <b>diverse</b> community
              warrants a platform that takes accessibility seriously,{" "}
              <b>no questions asked</b>.
            </styled.p>
          </VStack>

          <Box p={{ base: 0, sm: 4 }}>
            <svg
              width="100%"
              height="100%"
              viewBox="0 0 161 161"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M80.5 15.2256C116.423 15.2256 145.5 44.2971 145.5 80.2256C145.5 116.148 116.428 145.226 80.5 145.226C44.5772 145.226 15.5 116.154 15.5 80.2256C15.5 44.3028 44.5716 15.2256 80.5 15.2256ZM80.5 2.72559C37.6978 2.72559 3 37.4234 3 80.2256C3 123.028 37.6978 157.726 80.5 157.726C123.302 157.726 158 123.028 158 80.2256C158 37.4234 123.302 2.72559 80.5 2.72559ZM80.5 20.2256C47.3628 20.2256 20.5 47.0884 20.5 80.2256C20.5 113.363 47.3628 140.226 80.5 140.226C113.637 140.226 140.5 113.363 140.5 80.2256C140.5 47.0884 113.637 20.2256 80.5 20.2256ZM80.5 33.9756C86.7131 33.9756 91.75 39.0125 91.75 45.2256C91.75 51.4387 86.7131 56.4756 80.5 56.4756C74.2869 56.4756 69.25 51.4387 69.25 45.2256C69.25 39.0125 74.2869 33.9756 80.5 33.9756ZM117.294 64.6078C108.322 66.7262 99.9469 68.5915 91.6253 69.5475C91.8913 101.117 95.4709 108.001 99.4494 118.179C100.58 121.073 99.1503 124.335 96.2566 125.465C93.3625 126.595 90.1006 125.166 88.9703 122.272C86.25 115.301 83.6309 109.573 82.0137 97.7256H78.9869C77.3722 109.554 74.7575 115.291 72.03 122.272C70.9003 125.164 67.6394 126.596 64.7441 125.465C61.8503 124.335 60.4209 121.072 61.5513 118.179C65.5241 108.01 69.1091 101.135 69.3753 69.5475C61.0537 68.5918 52.6791 66.7265 43.7062 64.6078C41.0187 63.9731 39.3544 61.2803 39.9891 58.5925C40.6237 55.9046 43.3162 54.2406 46.0044 54.8753C76.2188 62.0093 84.8428 61.995 114.997 54.8753C117.684 54.2409 120.377 55.9046 121.012 58.5925C121.646 61.2803 119.982 63.9734 117.294 64.6078Z"
                fill="#303030"
              />
            </svg>
          </Box>
        </HStack>
      </VStack>
    </VStack>
  );
}

function ForDevops() {
  return (
    <VStack
      w="full"
      bgColor="#0C0A14"
      py={24}
      px={{ base: 4, sm: 12, md: 16, lg: 48, xl: 96 }}
      gap={12}
      color="hsla(160, 9%, 92%, 1)"
    >
      <Pair
        heading={
          <styled.span color="hsla(160, 9%, 92%, 1)">
            dev-ops heroes
          </styled.span>
        }
        headingColour="hsla(160, 7%, 59%, 1)"
      >
        <styled.p>
          Simple <b>deployment</b>, simple <b>maintenance</b> and simple{" "}
          <b>updates</b>. That’s what matters when you’re self-hosting, so you
          can spend time where it brings the most value.
        </styled.p>
        <styled.p>
          The choice of technologies behind the Storyden platform are all{" "}
          <b>meticulously</b> intentional to fit those values of simplicity in
          order to <b>get out of the way</b>.
        </styled.p>
      </Pair>

      <Flex
        flexDir={{ base: "column-reverse", lg: "row" }}
        alignItems="center"
        gap={4}
        maxW={{ base: "full", sm: "container.lg" }}
      >
        <Box flex="2 0 auto">
          <Image
            alt="a silly fake terminal example"
            src="/terminal.png"
            width="457"
            height="92"
          />
        </Box>

        <VStack flex="1 0 1" alignItems="start">
          <FeatureHeading>Container first</FeatureHeading>
          <styled.p>
            Zero <b>installation</b> steps. Set some environment variables and
            spin up a container image. Behaving like all other modern
            server-side software is the key to <b>simplicity</b>.
          </styled.p>
        </VStack>
      </Flex>

      <HStack maxW={{ base: "full", sm: "container.md" }}>
        <VStack alignItems="start" flex="0 1 auto">
          <FeatureHeading>Bring your own frontend</FeatureHeading>
          <styled.p>
            Not a fan of themes? That’s fine, <b>headless</b> mode is for you.
            Storyden is, at the core, a powerful API service with which you can{" "}
            <b>wire up</b> anything you want.
          </styled.p>
          <styled.p>
            From web to mobile apps and everything in between, the{" "}
            <b>OpenAPI</b> specification provides a fast integration path.
          </styled.p>
        </VStack>

        <Box flex="2 0 auto">
          <Image
            alt="openapi logo"
            src="/openapi.png"
            width="127"
            height="126"
          />
        </Box>
      </HStack>
    </VStack>
  );
}

function ForYou() {
  return (
    <VStack
      w="full"
      bgColor="red.100"
      py={24}
      px={{ base: 4, sm: 12, md: 16, lg: 48, xl: 96 }}
      gap={12}
    >
      <Pair heading="you and your friends" headingColour="#808080">
        <styled.p>
          Optimised for <b>humans</b>, ready for the web <b>renaissance</b>.
          Storyden is built to be a stable foundation for the future decades of
          internet citizens and the <b>networks</b> they build.
        </styled.p>
      </Pair>

      <HStack maxW={{ base: "full", sm: "container.md" }} gap={4}>
        <styled.p textAlign="right" width="60%">
          Rough relationship with <b>email</b>? Just don’t enable it then. Sign
          in with <b>Passkey</b>, WebAuthn, Web3 <b>wallet</b>, or choose from a
          variety of popular OAuth2 and <b>SSO</b> providers.
        </styled.p>

        <Box flex="2 0 auto">
          <svg
            width="129"
            height="129"
            viewBox="0 0 129 129"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M74.9533 21.4923V16.1056C74.9533 13.1612 77.3488 10.7658 80.2931 10.7658H85.6798C86.3524 10.7658 86.8977 10.2205 86.8977 9.54794C86.8977 8.87535 86.3524 8.33008 85.6798 8.33008H80.2931C76.0057 8.33008 72.5176 11.8182 72.5176 16.1056V21.4923C72.5176 22.1649 73.0629 22.7102 73.7354 22.7102C74.408 22.7102 74.9533 22.1649 74.9533 21.4923ZM85.6798 53.8592C86.3524 53.8592 86.8977 54.4045 86.8977 55.0771C86.8977 55.7496 86.3524 56.2949 85.6798 56.2949H80.2931C76.0057 56.2949 72.5176 52.8068 72.5176 48.5194V43.1327C72.5176 42.4601 73.0629 41.9148 73.7354 41.9148C74.408 41.9148 74.9533 42.4601 74.9533 43.1327V48.5194C74.9533 51.4637 77.3488 53.8592 80.2931 53.8592H85.6798ZM120.482 43.1327V48.5194C120.482 52.8068 116.994 56.2949 112.707 56.2949H107.32C106.648 56.2949 106.102 55.7496 106.102 55.0771C106.102 54.4045 106.648 53.8592 107.32 53.8592H112.707C115.651 53.8592 118.047 51.4637 118.047 48.5194V43.1327C118.047 42.4601 118.592 41.9148 119.265 41.9148C119.937 41.9148 120.482 42.4601 120.482 43.1327ZM120.482 16.1056V21.4923C120.482 22.1649 119.937 22.7102 119.265 22.7102C118.592 22.7102 118.047 22.1649 118.047 21.4923V16.1056C118.047 13.1612 115.651 10.7658 112.707 10.7658H107.32C106.648 10.7658 106.102 10.2205 106.102 9.54794C106.102 8.87535 106.648 8.33008 107.32 8.33008H112.707C116.994 8.33008 120.482 11.8182 120.482 16.1056ZM104.905 43.7807C105.418 43.3068 105.449 42.5066 104.975 41.9935C104.501 41.4804 103.701 41.4487 103.188 41.9226C101.363 43.6091 98.9873 44.5378 96.5 44.5378C94.0127 44.5378 91.6374 43.6091 89.8117 41.9226C89.2985 41.4487 88.4985 41.4805 88.0245 41.9935C87.5506 42.5066 87.5823 43.3068 88.0954 43.7807C90.3899 45.9 93.3746 47.0672 96.5 47.0672C99.6253 47.0672 102.61 45.9 104.905 43.7807ZM99.0294 26.4106V35.404C99.0294 37.1861 97.5795 38.636 95.7974 38.636H94.7669C94.0684 38.636 93.5022 38.0698 93.5022 37.3713C93.5022 36.6728 94.0684 36.1066 94.7669 36.1066H95.7974C96.1848 36.1066 96.5 35.7914 96.5 35.404V26.4106C96.5 25.7121 97.0662 25.1459 97.7647 25.1459C98.4632 25.1459 99.0294 25.7121 99.0294 26.4106ZM107.742 29.947V26.3403C107.742 25.6807 107.207 25.1459 106.547 25.1459C105.888 25.1459 105.353 25.6807 105.353 26.3403V29.947C105.353 30.6067 105.888 31.1415 106.547 31.1415C107.207 31.1415 107.742 30.6067 107.742 29.947ZM85.5393 29.947C85.5393 30.6067 86.0741 31.1415 86.7337 31.1415C87.3934 31.1415 87.9282 30.6067 87.9282 29.947V26.3403C87.9282 25.6807 87.3934 25.1459 86.7337 25.1459C86.0741 25.1459 85.5393 25.6807 85.5393 26.3403V29.947Z"
              fill="#303030"
            />
            <path
              d="M44.4949 10.8751C44.9419 10.4281 45.6226 10.4281 46.0697 10.8751L53.9436 18.75C54.38 19.1865 54.3854 19.883 53.9436 20.325L49.8715 24.3975C49.4366 24.8324 49.4366 25.5376 49.8715 25.9725C50.3064 26.4074 51.0114 26.4074 51.4463 25.9725L55.5183 21.8999C56.8402 20.5779 56.8142 18.471 55.5183 17.175L47.6444 9.30018C46.3277 7.98327 44.2368 7.98327 42.9201 9.30018L38.8481 13.3727C38.4132 13.8076 38.4132 14.5128 38.8481 14.9477C39.2829 15.3826 39.988 15.3826 40.4228 14.9477L44.4949 10.8751Z"
              fill="#303030"
            />
            <path
              d="M46.0697 15.6002C46.5045 15.1653 46.5045 14.4601 46.0697 14.0252C45.6348 13.5903 44.9297 13.5903 44.4949 14.0252L42.9201 15.6002C42.4852 16.0351 42.4852 16.7403 42.9201 17.1752C43.355 17.6101 44.06 17.6101 44.4949 17.1752L46.0697 15.6002Z"
              fill="#303030"
            />
            <path
              d="M50.794 18.7502C51.2288 19.1851 51.2288 19.8902 50.794 20.3252L49.2192 21.9001C48.7843 22.335 48.0793 22.335 47.6444 21.9001C47.2095 21.4652 47.2095 20.7601 47.6444 20.3252L49.2192 18.7502C49.6541 18.3153 50.3591 18.3153 50.794 18.7502Z"
              fill="#303030"
            />
            <path
              d="M33.4714 18.7502C33.9063 18.3153 34.6114 18.3153 35.0462 18.7502L46.0697 29.775C46.5045 30.2099 46.5045 30.915 46.0697 31.35C45.6348 31.7849 44.9298 31.7849 44.4949 31.35L33.4714 20.3252C33.0366 19.8902 33.0366 19.1851 33.4714 18.7502Z"
              fill="#303030"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M34.2588 43.162C30.7799 46.6413 25.1395 46.6413 21.6606 43.162C18.1817 39.6827 18.1817 34.0416 21.6606 30.5622C25.1395 27.0829 30.7799 27.0829 34.2588 30.5622C37.7377 34.0416 37.7377 39.6827 34.2588 43.162ZM32.684 41.587C30.0749 44.1965 25.8445 44.1965 23.2354 41.587C20.6262 38.9775 20.6262 34.7467 23.2354 32.1372C25.8445 29.5277 30.0749 29.5277 32.684 32.1372C35.2932 34.7467 35.2932 38.9775 32.684 41.587Z"
              fill="#303030"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M36.621 14.0251C35.3043 12.7082 33.2134 12.7082 31.8967 14.0251L11.4245 34.4998C7.5232 38.4016 7.5283 44.7424 11.421 48.6709L11.4245 48.6745L16.1489 53.3994L16.1525 53.403C20.0667 57.2823 26.404 57.2823 30.3183 53.403L30.3219 53.3994L50.794 32.9248C52.1159 31.6028 52.0899 29.4959 50.794 28.1999L36.621 14.0251ZM33.4714 15.6001C33.9185 15.153 34.5992 15.153 35.0462 15.6001L49.2192 29.7748C49.6557 30.2113 49.6611 30.9079 49.2192 31.3498L28.7491 51.8225C25.7034 54.8393 20.7678 54.8394 17.722 51.8228L13.001 47.1012C9.96652 44.0367 9.97225 39.1022 12.9993 36.0747L33.4714 15.6001Z"
              fill="#303030"
            />
            <path
              d="M7.75 77.6875C7.75 77.0488 8.23632 76.5625 8.875 76.5625H56.125C56.7485 76.5625 57.25 77.0561 57.25 77.6875V88.9825C57.25 89.6038 57.7537 90.1075 58.375 90.1075C58.9963 90.1075 59.5 89.6038 59.5 88.9825V77.6875C59.5 75.7989 57.9765 74.3125 56.125 74.3125H8.875C6.99368 74.3125 5.5 75.8062 5.5 77.6875V93.4375C5.5 95.289 6.98635 96.8125 8.875 96.8125H29.125C29.7463 96.8125 30.25 96.3088 30.25 95.6875C30.25 95.0662 29.7463 94.5625 29.125 94.5625H8.875C8.24365 94.5625 7.75 94.061 7.75 93.4375V77.6875Z"
              fill="#303030"
            />
            <path
              d="M16.75 81.0625C16.75 80.4412 16.2463 79.9375 15.625 79.9375C15.0037 79.9375 14.5 80.4412 14.5 81.0625V83.4569L12.5153 82.1298C11.9988 81.7845 11.3002 81.9232 10.9548 82.4397C10.6094 82.9562 10.7482 83.6549 11.2647 84.0002L13.5906 85.5555L11.2487 87.0981C10.7298 87.4398 10.5862 88.1375 10.928 88.6564C11.2698 89.1753 11.9675 89.3188 12.4863 88.9771L14.5 87.6507V90.0625C14.5 90.6838 15.0037 91.1875 15.625 91.1875C16.2463 91.1875 16.75 90.6838 16.75 90.0625V87.6681L18.7347 88.9952C19.2511 89.3406 19.9498 89.2019 20.2952 88.6854C20.6405 88.1689 20.5018 87.4702 19.9853 87.1248L17.6595 85.5696L20.0013 84.0271C20.5202 83.6853 20.6638 82.9876 20.322 82.4687C19.9802 81.9498 19.2825 81.8063 18.7637 82.1481L16.75 83.4744V81.0625Z"
              fill="#303030"
            />
            <path
              d="M29.125 79.9375C29.7463 79.9375 30.25 80.4412 30.25 81.0625V83.4744L32.2637 82.1481C32.7825 81.8063 33.4802 81.9498 33.822 82.4687C34.1638 82.9876 34.0202 83.6853 33.5013 84.0271L31.1595 85.5696L33.4853 87.1248C34.0018 87.4702 34.1405 88.1689 33.7952 88.6854C33.4498 89.2019 32.7511 89.3406 32.2347 88.9952L30.25 87.6681V90.0625C30.25 90.6838 29.7463 91.1875 29.125 91.1875C28.5037 91.1875 28 90.6838 28 90.0625V87.6507L25.9863 88.9771C25.4675 89.3188 24.7698 89.1753 24.428 88.6564C24.0862 88.1375 24.2298 87.4398 24.7487 87.0981L27.0906 85.5555L24.7647 84.0002C24.2482 83.6549 24.1094 82.9562 24.4548 82.4397C24.8002 81.9232 25.4988 81.7845 26.0153 82.1298L28 83.4569V81.0625C28 80.4412 28.5037 79.9375 29.125 79.9375Z"
              fill="#303030"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M43.2158 106.538C43.9542 105.8 44.9557 105.385 46 105.385C47.0443 105.385 48.0458 105.8 48.7842 106.538C49.5227 107.277 49.9375 108.278 49.9375 109.323C49.9375 110.367 49.5227 111.368 48.7842 112.107C48.0458 112.845 47.0443 113.26 46 113.26C44.9557 113.26 43.9542 112.845 43.2158 112.107C42.4773 111.368 42.0625 110.367 42.0625 109.323C42.0625 108.278 42.4773 107.277 43.2158 106.538ZM46 107.635C45.5524 107.635 45.1232 107.813 44.8068 108.129C44.4903 108.446 44.3125 108.875 44.3125 109.323C44.3125 109.77 44.4903 110.199 44.8068 110.516C45.1232 110.832 45.5524 111.01 46 111.01C46.4476 111.01 46.8768 110.832 47.1932 110.516C47.5097 110.199 47.6875 109.77 47.6875 109.323C47.6875 108.875 47.5097 108.446 47.1932 108.129C46.8768 107.813 46.4476 107.635 46 107.635Z"
              fill="#303030"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M38.125 99.0647C36.261 99.0647 34.75 100.576 34.75 102.44V115.937C34.75 117.801 36.261 119.312 38.125 119.312H53.875C55.739 119.312 57.25 117.801 57.25 115.937V102.44C57.25 100.576 55.739 99.0647 53.875 99.0647H52.75V96.8125C52.75 93.0607 49.7234 90.0625 46 90.0625C42.2512 90.0625 39.25 93.0637 39.25 96.8125V99.0647H38.125ZM37 102.44C37 101.818 37.5037 101.315 38.125 101.315H53.875C54.4963 101.315 55 101.818 55 102.44V115.937C55 116.559 54.4963 117.062 53.875 117.062H38.125C37.5037 117.062 37 116.559 37 115.937V102.44ZM41.5 96.8125C41.5 94.3063 43.4938 92.3125 46 92.3125C48.4866 92.3125 50.5 94.3093 50.5 96.8125V99.0625H41.5V96.8125Z"
              fill="#303030"
            />
            <path
              d="M112.887 75.9815C112.663 75.9815 112.44 75.9275 112.244 75.8195C106.878 73.1465 102.238 72.0125 96.6482 72.0125C91.1421 72.0125 85.8876 73.2815 81.1082 75.8195C80.4375 76.1705 79.599 75.9275 79.2077 75.2795C78.8443 74.6315 79.0959 73.7945 79.7667 73.4435C84.9653 70.6625 90.667 69.3125 96.6482 69.3125C102.629 69.3125 107.828 70.5815 113.53 73.3625C114.228 73.7675 114.48 74.5775 114.117 75.2255C113.865 75.7115 113.418 75.9815 112.887 75.9815ZM72.8911 90.1565C72.6116 90.1565 72.3321 90.0755 72.0806 89.9135C71.4936 89.4815 71.298 88.6445 71.7452 88.0235C74.5122 84.2435 78.0338 81.2735 82.2262 79.1945C91.0583 74.8205 102.238 74.7935 111.042 79.1675C115.235 81.2465 118.756 84.1625 121.523 87.9425C121.97 88.5365 121.803 89.4005 121.188 89.8325C120.545 90.2645 119.679 90.1295 119.231 89.5625C116.716 86.1065 113.53 83.4335 109.756 81.5705C101.735 77.6015 91.4775 77.6015 83.4839 81.5975C79.6828 83.4875 76.4966 86.1875 73.9811 89.5625C73.7575 89.9675 73.3383 90.1565 72.8911 90.1565ZM90.3595 122.746C89.9962 122.746 89.6608 122.611 89.3813 122.34C86.9497 119.992 85.6361 118.479 83.7634 115.213C81.8349 111.892 80.8287 107.841 80.8287 103.494C80.8287 95.4755 87.9279 88.9415 96.6482 88.9415C105.368 88.9415 112.468 95.4755 112.468 103.494C112.468 103.853 112.32 104.196 112.058 104.449C111.796 104.702 111.441 104.845 111.07 104.845C110.699 104.845 110.344 104.702 110.082 104.449C109.82 104.196 109.673 103.853 109.673 103.494C109.673 96.9605 103.831 91.6415 96.6482 91.6415C89.4651 91.6415 83.6237 96.9605 83.6237 103.494C83.6237 107.382 84.5181 110.973 86.223 113.862C88.0118 116.994 89.2415 118.318 91.3937 120.424C91.9247 120.964 91.9247 121.8 91.3937 122.34C91.0583 122.611 90.7229 122.746 90.3595 122.746ZM110.399 117.751C107.073 117.751 104.139 116.94 101.735 115.347C97.5705 112.62 95.083 108.193 95.083 103.494C95.083 103.136 95.2302 102.793 95.4923 102.54C95.7544 102.287 96.1098 102.145 96.4805 102.145C96.8511 102.145 97.2065 102.287 97.4686 102.54C97.7307 102.793 97.8779 103.136 97.8779 103.494C97.8779 107.301 99.8903 110.892 103.3 113.106C105.285 114.403 107.604 115.024 110.399 115.024C111.07 115.024 112.188 114.943 113.306 114.754C114.061 114.618 114.815 115.104 114.927 115.861C115.067 116.562 114.564 117.292 113.781 117.426C112.188 117.723 110.791 117.751 110.399 117.751ZM104.781 123.312C104.67 123.312 104.53 123.312 104.418 123.312C99.9741 122.07 97.0674 120.477 94.0209 117.588C90.108 113.862 87.9559 108.84 87.9559 103.494C87.9559 99.1205 91.8129 95.5565 96.5643 95.5565C101.316 95.5565 105.173 99.1205 105.173 103.494C105.173 106.383 107.828 108.733 110.986 108.733C114.2 108.733 116.8 106.383 116.8 103.494C116.8 93.3155 107.716 85.0535 96.5364 85.0535C88.5987 85.0535 81.2759 89.3195 78.0617 95.9345C76.9717 98.1215 76.4127 100.687 76.4127 103.494C76.4127 105.601 76.6084 108.922 78.2853 113.242C78.5648 113.944 78.2015 114.727 77.4748 114.97C76.7481 115.213 75.9376 114.862 75.686 114.187C74.2886 110.65 73.6457 107.112 73.6457 103.494C73.6457 100.254 74.2886 97.3115 75.5463 94.7465C79.2636 87.2135 87.5087 82.3265 96.5364 82.3265C109.225 82.3265 119.595 91.8035 119.595 103.467C119.595 107.841 115.738 111.406 110.986 111.406C106.235 111.406 102.378 107.841 102.378 103.467C102.378 100.578 99.7785 98.2295 96.5643 98.2295C93.3501 98.2295 90.7508 100.578 90.7508 103.467C90.7508 108.084 92.5955 112.404 95.9774 115.645C98.6326 118.182 101.176 119.587 105.117 120.612C105.871 120.829 106.291 121.585 106.095 122.287C105.955 122.908 105.368 123.312 104.781 123.312Z"
              fill="#303030"
            />
          </svg>
        </Box>
      </HStack>

      <HStack maxW={{ base: "full", sm: "container.sm" }} gap={4}>
        <Box flex="2 0 auto">
          <svg
            width="101"
            height="101"
            viewBox="0 0 101 101"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M50.4997 3.9375C24.0257 3.9375 2.58301 25.3802 2.58301 51.8542C2.58301 73.0573 16.2992 90.9661 35.346 97.3151C37.7419 97.7344 38.6403 96.2969 38.6403 95.0391C38.6403 93.901 38.5804 90.1276 38.5804 86.1146C26.5413 88.3307 23.4268 83.1797 22.4684 80.4844C21.9294 79.1068 19.5934 74.8542 17.557 73.7161C15.8799 72.8177 13.484 70.6016 17.4971 70.5417C21.2705 70.4818 23.9658 74.0156 24.8643 75.4531C29.1768 82.7005 36.0648 80.6641 38.82 79.4062C39.2393 76.2917 40.4971 74.1953 41.8747 72.9974C31.2132 71.7995 20.0726 67.6667 20.0726 49.3385C20.0726 44.1276 21.9294 39.8151 24.984 36.4609C24.5049 35.263 22.8278 30.3516 25.4632 23.763C25.4632 23.763 29.4762 22.5052 38.6403 28.6745C42.4736 27.5964 46.5465 27.0573 50.6195 27.0573C54.6924 27.0573 58.7653 27.5964 62.5986 28.6745C71.7627 22.4453 75.7757 23.763 75.7757 23.763C78.4111 30.3516 76.734 35.263 76.2549 36.4609C79.3096 39.8151 81.1663 44.0677 81.1663 49.3385C81.1663 67.7266 69.9658 71.7995 59.3044 72.9974C61.0413 74.4948 62.5387 77.3698 62.5387 81.862C62.5387 88.2708 62.4788 93.4219 62.4788 95.0391C62.4788 96.2969 63.3773 97.7943 65.7731 97.3151C84.7002 90.9661 98.4163 72.9974 98.4163 51.8542C98.4163 25.3802 76.9736 3.9375 50.4997 3.9375V3.9375Z"
              fill="#24292F"
            />
          </svg>
        </Box>

        <styled.p flex="0 1 auto">
          Staying <b>closed source</b> is pointless in today’s internet. Fork
          it, hack on it, provide hosting, use as a basis for other apps,{" "}
          <b>contribute</b> back to the community. Not sure where to start? Use
          Storyden to <b>learn</b> about building!
        </styled.p>
      </HStack>
    </VStack>
  );
}

function Pair({ heading, headingColour, children }) {
  return (
    <Flex px={{ base: 1, md: 9, lg: 12, xl: 12 }} flexDir="column" gap={6}>
      <styled.h1
        width="min-content"
        textAlign="left"
        whiteSpace="nowrap"
        fontWeight="black"
        fontSize={{
          base: "2xl",
          md: "3xl",
          lg: "3xl",
          xl: "4xl",
          "2xl": "5xl",
        }}
      >
        <styled.span color={headingColour}>For</styled.span> {heading}
      </styled.h1>
      <Box mt={1}>{children}</Box>
    </Flex>
  );
}

function CTA() {
  return (
    <VStack
      bgColor="hsla(160, 9%, 92%, 1)"
      color="hsla(0, 0%, 19%, 1)"
      p={8}
      gap={2}
      w="full"
      textAlign="center"
    >
      <styled.h1 fontWeight="bold" fontSize={{ base: "2xl", lg: "4xl" }}>
        Interested?
      </styled.h1>
      <styled.p>
        Storyden is early in development and is looking for <b>feedback</b> and{" "}
        <b>contributors</b>!
      </styled.p>
      <styled.p>
        If you have <b>opinions</b> about forum software, please click{" "}
        <styled.a
          href="https://airtable.com/shrLY0jDp9CuXPB2X"
          color="hsla(265, 56%, 42%, 1)"
        >
          this link!
        </styled.a>
      </styled.p>
      <styled.p>
        If you're feeling <b>social</b>, fancy chatting with likeminded folks
        and want to steer Storyden's development as an early bird, please click{" "}
        <styled.a
          href="https://discord.gg/XF6ZBGF9XF"
          color="hsla(265, 56%, 42%, 1)"
        >
          this link!
        </styled.a>
      </styled.p>
      <styled.p>
        If you know <b>Golang</b> or <b>React.js</b> and are interested in
        contributing to a <b>high-quality</b> open source project, please click{" "}
        <styled.a
          href="https://github.com/Southclaws/storyden"
          color="hsla(265, 56%, 42%, 1)"
        >
          this link!
        </styled.a>
      </styled.p>
    </VStack>
  );
}

function Discord(props) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 127.14 96.36"
      fill="hsla(265, 56%, 42%, 1)"
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

function Twitter(props) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 248 204"
      fill="hsla(265, 56%, 42%, 1)"
      style={{
        enableBackground: "new 0 0 248 204",
      }}
      {...props}
    >
      <path d="M221.95 51.29c.15 2.17.15 4.34.15 6.53 0 66.73-50.8 143.69-143.69 143.69v-.04c-27.44.04-54.31-7.82-77.41-22.64 3.99.48 8 .72 12.02.73 22.74.02 44.83-7.61 62.72-21.66-21.61-.41-40.56-14.5-47.18-35.07a50.338 50.338 0 0 0 22.8-.87C27.8 117.2 10.85 96.5 10.85 72.46v-.64a50.18 50.18 0 0 0 22.92 6.32C11.58 63.31 4.74 33.79 18.14 10.71a143.333 143.333 0 0 0 104.08 52.76 50.532 50.532 0 0 1 14.61-48.25c20.34-19.12 52.33-18.14 71.45 2.19 11.31-2.23 22.15-6.38 32.07-12.26a50.69 50.69 0 0 1-22.2 27.93c10.01-1.18 19.79-3.86 29-7.95a102.594 102.594 0 0 1-25.2 26.16z" />
    </svg>
  );
}

function GitHub(props) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 98 96"
      fill="hsla(265, 56%, 42%, 1)"
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

function Footer() {
  return (
    <VStack
      bgColor="hsla(140, 16%, 88%, 1)"
      px={{ base: 6, md: 12, lg: 24, xl: 48, "2xl": 80 }}
      py={{ base: 8, lg: 12 }}
      gap={4}
      flex="1"
    >
      <HStack gap={[2, 3, 4]}>
        <VStack alignItems="end" fontSize="sm" color="hsla(265, 56%, 42%, 1)">
          <styled.a href="https://discord.gg/bFkJDvpcmg">
            <HStack>
              <styled.p>Discord</styled.p>
              <Discord width="1.2em" />
            </HStack>
          </styled.a>

          <styled.a href="https://github.com/Southclaws/storyden">
            <HStack>
              <styled.p>GitHub</styled.p>
              <GitHub width="1.2em" />
            </HStack>
          </styled.a>

          <styled.a href="https://twitter.com/Southclaws">
            <HStack>
              <styled.p>Twitter</styled.p>
              <Twitter width="1.2em" />
            </HStack>
          </styled.a>
        </VStack>

        <Divider orientation="vertical" />

        <Image
          src="/mark.png"
          alt="The Storyden logomark and wordmark"
          width={150}
          height={50}
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

export default function Home() {
  return (
    <Box>
      <NextSeo title="Storyden: A forum for the modern age." {...seo} />
      <Nav />
      <Hero />
      <Story />
      <Why />
      <Features />
      <ForCommunityLeaders />
      <ForDevops />
      <ForYou />
      <CTA />
      <Footer />
    </Box>
  );
}
