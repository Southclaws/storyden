import {
  Box,
  Button,
  Flex,
  Grid,
  GridItem,
  Heading,
  HStack,
  Link,
  Text,
  VStack,
} from "@chakra-ui/react";
import Image from "next/image";
import localFont from "@next/font/local";

const monasans = localFont({ src: "./Mona-Sans.woff2" });

function Logo() {
  return (
    <svg viewBox="0 0 142 143" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g filter="url(#filter0_d_1118_12433)">
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M66.274 23.2302C67.7837 21.7351 70.216 21.7351 71.7257 23.2302L94.3731 45.6577C95.1078 46.3853 95.5211 47.3763 95.5211 48.4103V78.3796C95.5211 79.8608 94.3204 81.0616 92.8392 81.0616C91.358 81.0616 90.1573 79.8608 90.1573 78.3796V49.5289C90.1573 49.2107 90.0301 48.9058 89.804 48.6819L69.8386 28.9103C69.374 28.4503 68.6256 28.4503 68.1611 28.9103L61.9264 35.0844C60.874 36.1267 59.1759 36.1184 58.1336 35.0659C57.0914 34.0135 57.0997 32.3154 58.1521 31.2732L66.274 23.2302ZM59.9454 51.8012C59.476 52.266 59.4741 53.0239 59.9412 53.491L68.157 61.7067C68.6225 62.1722 69.3772 62.1722 69.8427 61.7067L78.0584 53.491C78.5255 53.0239 78.5237 52.266 78.0543 51.8012L69.8386 43.6652C69.374 43.2052 68.6256 43.2052 68.1611 43.6652L59.9454 51.8012ZM73.6355 67.1853C73.17 66.7198 73.17 65.9651 73.6355 65.4996L83.7569 55.3781C85.275 53.8601 85.269 51.3969 83.7435 49.8863L71.7257 37.9851C70.216 36.4901 67.7837 36.4901 66.274 37.9851L54.2561 49.8863C52.7307 51.3969 52.7247 53.8601 54.2427 55.3781L64.3642 65.4996C64.8296 65.965 64.8296 66.7198 64.3642 67.1853L56.2008 75.3486C55.9773 75.5721 55.6741 75.6977 55.358 75.6977H49.0343C48.376 75.6977 47.8424 75.164 47.8424 74.5057V49.5243C47.8424 49.207 47.9689 48.9028 48.1939 48.6791L52.6352 44.2632C53.6856 43.2189 53.6905 41.5208 52.6461 40.4704C51.6018 39.4201 49.9037 39.4152 48.8533 40.4595L43.621 45.6618C42.8897 46.3889 42.4785 47.3776 42.4785 48.4089V77.1876C42.4785 79.3271 44.2129 81.0616 46.3524 81.0616H56.4689C57.4963 81.0616 58.4816 80.6534 59.2081 79.9269L68.157 70.9781C68.6225 70.5126 69.3772 70.5126 69.8427 70.9781L79.1406 80.276C80.188 81.3234 81.8861 81.3234 82.9334 80.276C83.9808 79.2287 83.9808 77.5306 82.9334 76.4832L73.6355 67.1853Z"
          fill="white"
        />
      </g>
      <g filter="url(#filter1_d_1118_12433)">
        <path
          d="M31.7397 98.8247C31.2845 98.4549 30.8292 98.1846 30.374 98.0139C29.9187 97.8289 29.4777 97.7364 29.0509 97.7364C28.5103 97.7364 28.0693 97.8645 27.7279 98.1206C27.3865 98.3766 27.2158 98.7109 27.2158 99.1235C27.2158 99.408 27.3011 99.6428 27.4718 99.8277C27.6425 100.013 27.8631 100.176 28.1333 100.318C28.4179 100.447 28.7308 100.56 29.0723 100.66C29.4279 100.76 29.7765 100.866 30.1179 100.98C31.4836 101.435 32.4795 102.047 33.1054 102.815C33.7456 103.569 34.0657 104.558 34.0657 105.781C34.0657 106.606 33.9234 107.353 33.6389 108.022C33.3686 108.691 32.9631 109.267 32.4225 109.75C31.8962 110.22 31.2418 110.583 30.4593 110.839C29.6911 111.109 28.8162 111.244 27.8346 111.244C25.8003 111.244 23.9153 110.64 22.1797 109.43L23.9722 106.059C24.5981 106.614 25.217 107.026 25.8287 107.296C26.4404 107.567 27.045 107.702 27.6425 107.702C28.3254 107.702 28.8304 107.545 29.1576 107.232C29.4991 106.919 29.6698 106.564 29.6698 106.165C29.6698 105.924 29.6271 105.717 29.5417 105.547C29.4564 105.362 29.3141 105.198 29.115 105.056C28.9158 104.899 28.6526 104.757 28.3254 104.629C28.0124 104.501 27.6283 104.359 27.1731 104.202C26.6325 104.032 26.099 103.847 25.5726 103.647C25.0605 103.434 24.5981 103.157 24.1856 102.815C23.7872 102.474 23.46 102.047 23.204 101.535C22.9621 101.008 22.8412 100.347 22.8412 99.5503C22.8412 98.7536 22.9692 98.0352 23.2253 97.395C23.4956 96.7406 23.8655 96.1858 24.335 95.7305C24.8186 95.2611 25.4019 94.8983 26.0848 94.6422C26.7819 94.3862 27.5572 94.2581 28.4108 94.2581C29.2074 94.2581 30.0397 94.3719 30.9075 94.5996C31.7753 94.813 32.6075 95.133 33.4042 95.5598L31.7397 98.8247Z"
          fill="white"
        />
        <path
          d="M39.3642 103.071V110.796H35.5018V103.071H34.2214V99.849H35.5018V96.5628H39.3642V99.849H41.5622V103.071H39.3642Z"
          fill="white"
        />
        <path
          d="M45.9951 105.269C45.9951 105.639 46.0591 105.981 46.1872 106.293C46.3294 106.592 46.5073 106.855 46.7207 107.083C46.9483 107.311 47.2115 107.488 47.5102 107.617C47.8232 107.745 48.1504 107.809 48.4918 107.809C48.8332 107.809 49.1533 107.745 49.4521 107.617C49.7651 107.488 50.0282 107.311 50.2416 107.083C50.4693 106.855 50.6471 106.592 50.7751 106.293C50.9174 105.981 50.9885 105.646 50.9885 105.291C50.9885 104.949 50.9174 104.629 50.7751 104.33C50.6471 104.017 50.4693 103.747 50.2416 103.519C50.0282 103.292 49.7651 103.114 49.4521 102.986C49.1533 102.858 48.8332 102.794 48.4918 102.794C48.1504 102.794 47.8232 102.858 47.5102 102.986C47.2115 103.114 46.9483 103.292 46.7207 103.519C46.5073 103.747 46.3294 104.01 46.1872 104.309C46.0591 104.608 45.9951 104.928 45.9951 105.269ZM41.9193 105.227C41.9193 104.416 42.0829 103.662 42.4101 102.965C42.7373 102.253 43.1926 101.642 43.7758 101.129C44.3591 100.603 45.0491 100.19 45.8457 99.8917C46.6566 99.593 47.5387 99.4436 48.4918 99.4436C49.4307 99.4436 50.2985 99.593 51.0952 99.8917C51.9061 100.176 52.6032 100.582 53.1865 101.108C53.784 101.62 54.2463 102.239 54.5735 102.965C54.9007 103.676 55.0643 104.465 55.0643 105.333C55.0643 106.201 54.8936 106.998 54.5522 107.723C54.225 108.435 53.7697 109.053 53.1865 109.58C52.6032 110.092 51.899 110.49 51.0739 110.775C50.263 111.059 49.381 111.202 48.4278 111.202C47.4889 111.202 46.6211 111.059 45.8244 110.775C45.0277 110.49 44.3378 110.085 43.7545 109.558C43.1854 109.032 42.7373 108.406 42.4101 107.681C42.0829 106.941 41.9193 106.123 41.9193 105.227Z"
          fill="white"
        />
        <path
          d="M55.7951 99.849H59.6575V101.642C60.07 100.987 60.5751 100.489 61.1726 100.148C61.7701 99.7921 62.4672 99.6143 63.2638 99.6143C63.3634 99.6143 63.4701 99.6143 63.5839 99.6143C63.712 99.6143 63.8542 99.6285 64.0107 99.657V103.349C63.4986 103.093 62.9437 102.965 62.3462 102.965C61.45 102.965 60.7742 103.235 60.319 103.775C59.878 104.302 59.6575 105.077 59.6575 106.101V110.796H55.7951V99.849Z"
          fill="white"
        />
        <path
          d="M69.7477 109.302L64.4556 99.849H68.9368L71.8603 105.461L74.6984 99.849H79.137L70.3452 116.259H66.0133L69.7477 109.302Z"
          fill="white"
        />
        <path
          d="M81.6828 105.291C81.6828 105.646 81.7468 105.981 81.8749 106.293C82.0029 106.592 82.1736 106.855 82.387 107.083C82.6146 107.311 82.8778 107.488 83.1766 107.617C83.4896 107.745 83.8239 107.809 84.1795 107.809C84.521 107.809 84.841 107.745 85.1398 107.617C85.4528 107.488 85.716 107.311 85.9293 107.083C86.157 106.855 86.3348 106.592 86.4628 106.293C86.6051 105.995 86.6762 105.675 86.6762 105.333C86.6762 104.992 86.6051 104.672 86.4628 104.373C86.3348 104.06 86.157 103.79 85.9293 103.562C85.716 103.334 85.4528 103.157 85.1398 103.029C84.841 102.901 84.521 102.837 84.1795 102.837C83.8381 102.837 83.5109 102.901 83.1979 103.029C82.8992 103.157 82.636 103.334 82.4084 103.562C82.195 103.79 82.0171 104.053 81.8749 104.352C81.7468 104.636 81.6828 104.949 81.6828 105.291ZM86.5482 93.0845H90.4319V110.796H86.5482V109.58C85.7231 110.618 84.6063 111.138 83.1979 111.138C82.4013 111.138 81.6686 110.988 81 110.689C80.3313 110.391 79.7481 109.978 79.2501 109.452C78.7522 108.925 78.361 108.307 78.0765 107.595C77.8062 106.884 77.671 106.116 77.671 105.291C77.671 104.508 77.7991 103.768 78.0551 103.071C78.3254 102.36 78.7024 101.741 79.1861 101.215C79.6698 100.688 80.246 100.276 80.9146 99.9771C81.5975 99.6641 82.3443 99.5076 83.1552 99.5076C84.521 99.5076 85.6519 99.9842 86.5482 100.937V93.0845Z"
          fill="white"
        />
        <path
          d="M99.5274 103.733C99.3994 103.192 99.1362 102.758 98.7379 102.431C98.3395 102.104 97.8558 101.94 97.2868 101.94C96.6893 101.94 96.1985 102.097 95.8144 102.41C95.4445 102.723 95.2098 103.164 95.1102 103.733H99.5274ZM95.0035 105.973C95.0035 107.638 95.7859 108.47 97.3508 108.47C98.1901 108.47 98.8232 108.129 99.25 107.446H102.984C102.23 109.95 100.345 111.202 97.3295 111.202C96.4048 111.202 95.5583 111.066 94.7901 110.796C94.0219 110.512 93.3604 110.113 92.8055 109.601C92.2649 109.089 91.8453 108.477 91.5465 107.766C91.2478 107.055 91.0984 106.258 91.0984 105.376C91.0984 104.465 91.2407 103.647 91.5252 102.922C91.8097 102.182 92.2151 101.556 92.7415 101.044C93.2679 100.532 93.901 100.141 94.6407 99.8704C95.3947 99.5858 96.2412 99.4436 97.1801 99.4436C98.1048 99.4436 98.937 99.5858 99.6768 99.8704C100.417 100.141 101.043 100.539 101.555 101.065C102.067 101.592 102.458 102.239 102.728 103.007C102.999 103.761 103.134 104.615 103.134 105.568V105.973H95.0035Z"
          fill="white"
        />
        <path
          d="M104.189 99.849H108.051V101.236C108.577 100.596 109.111 100.169 109.651 99.9557C110.192 99.7281 110.825 99.6143 111.551 99.6143C112.319 99.6143 112.973 99.7423 113.514 99.9984C114.069 100.24 114.538 100.589 114.922 101.044C115.235 101.414 115.449 101.826 115.562 102.282C115.676 102.737 115.733 103.256 115.733 103.839V110.796H111.871V105.269C111.871 104.729 111.828 104.295 111.743 103.968C111.672 103.626 111.536 103.356 111.337 103.157C111.167 102.986 110.974 102.865 110.761 102.794C110.548 102.723 110.32 102.687 110.078 102.687C109.424 102.687 108.919 102.886 108.563 103.285C108.222 103.669 108.051 104.224 108.051 104.949V110.796H104.189V99.849Z"
          fill="white"
        />
      </g>
      <defs>
        <filter
          id="filter0_d_1118_12433"
          x="20.4785"
          y="0.108887"
          width="101.042"
          height="106.953"
          filterUnits="userSpaceOnUse"
          colorInterpolationFilters="sRGB"
        >
          <feFlood floodOpacity="0" result="BackgroundImageFix" />
          <feColorMatrix
            in="SourceAlpha"
            type="matrix"
            values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
            result="hardAlpha"
          />
          <feOffset dx="2" dy="2" />
          <feGaussianBlur stdDeviation="12" />
          <feComposite in2="hardAlpha" operator="out" />
          <feColorMatrix
            type="matrix"
            values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.75 0"
          />
          <feBlend
            mode="normal"
            in2="BackgroundImageFix"
            result="effect1_dropShadow_1118_12433"
          />
          <feBlend
            mode="normal"
            in="SourceGraphic"
            in2="effect1_dropShadow_1118_12433"
            result="shape"
          />
        </filter>
        <filter
          id="filter1_d_1118_12433"
          x="-0.678711"
          y="66.5361"
          width="143.357"
          height="76.6074"
          filterUnits="userSpaceOnUse"
          colorInterpolationFilters="sRGB"
        >
          <feFlood floodOpacity="0" result="BackgroundImageFix" />
          <feColorMatrix
            in="SourceAlpha"
            type="matrix"
            values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
            result="hardAlpha"
          />
          <feOffset dx="2" dy="2" />
          <feGaussianBlur stdDeviation="12" />
          <feComposite in2="hardAlpha" operator="out" />
          <feColorMatrix
            type="matrix"
            values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.75 0"
          />
          <feBlend
            mode="normal"
            in2="BackgroundImageFix"
            result="effect1_dropShadow_1118_12433"
          />
          <feBlend
            mode="normal"
            in="SourceGraphic"
            in2="effect1_dropShadow_1118_12433"
            result="shape"
          />
        </filter>
      </defs>
    </svg>
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
        <VStack>
          <Box width={[20, 36, 40, 48]}>
            <Logo />
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
              <Heading as="h1" size={["xs", "sm", "lg"]}>
                A forum for the modern age
              </Heading>
            </GridItem>
          </Grid>

          <HStack gap={4}>
            <Link href="https://airtable.com/shrLY0jDp9CuXPB2X">
              <Button size="lg">Get started</Button>
            </Link>
            <Link href="https://github.com/Southclaws/storyden">
              <Button
                variant="outline"
                colorScheme="whiteAlpha"
                size="lg"
                backdropFilter="blur(3px)"
                className="story__text-overlay"
              >
                Source code
              </Button>
            </Link>
          </HStack>
        </VStack>
      </GridItem>
    </Grid>
  );
}

function Story() {
  return (
    <>
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
            <Text
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
            </Text>
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
            <Text
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
            </Text>
          </HStack>
        </GridItem>

        <GridItem
          gridRow="three / three-end"
          gridColumn="left-near / right-far"
        >
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
            <Text
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
            </Text>
          </HStack>
        </GridItem>
      </Grid>

      <style jsx global>{`
        @supports not (
          (-webkit-backdrop-filter: none) or (backdrop-filter: none)
        ) {
          .story__text-overlay {
            background-color: rgba(8, 8, 8, 0.4);
          }
        }
      `}</style>
    </>
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
        <Heading
          textAlign="right"
          fontFamily={monasans}
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
        </Heading>
      </Box>

      <Flex flexDir="column" gap={{ base: 2, lg: 1 }} pt="10.1px">
        <Text>Storyden is a platform for building communities.</Text>
        <Text wordBreak="keep-all">
          But not just another chat app or another forum site.
          <wbr /> Storyden is a modern take on oldschool bulletin board
          <wbr /> forums you may remember from the earlier days of the
          <wbr /> internet.
        </Text>
        <Text>
          With a fresh perspective and new objectives, Storyden is
          <wbr /> designed to be the community platform for the next era of
          <wbr /> internet culture.
        </Text>
      </Flex>
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
          <Heading textShadow="4px 4px 8px #000000">{heading}</Heading>
          <Text textShadow="2px 2px 4px #000000">{body}</Text>
        </Flex>
      </GridItem>
    </Grid>
  );
}

function Features() {
  return (
    <VStack bgColor="hsla(140, 16%, 88%, 1)" py={8}>
      <VStack maxW="container.lg" gap={8}>
        <Flex
          flexWrap="wrap"
          px={{ base: 6, md: 9, lg: 12, xl: 12 }}
          pb={12}
          alignItems="center"
          justifyContent="center"
          gap={8}
        >
          <Feature
            image="/accessible.webp"
            alt=""
            heading="Accessible"
            body="Accessibility is non-negotiable and no one can be left behind. WAI and WCAG are a primary focus to ensure great experience for people regardless of a disability."
          />

          <Feature
            pt={24}
            image="/secure.webp"
            alt=""
            heading="Secure"
            body="The latest and greatest industry standard security practices as well as new emerging systems such as WebAuthn guarantee the most secure experience for everyone."
          />
        </Flex>
        <Pair big heading="For community leaders">
          <Text>
            Fearless <b>futurism</b>, radical <b>accessibility</b>, endless
            extensibility. Every modern service, product and movement has
            community at the centre. Communities often grow out of their humble
            beginnings on walled-garden platforms. In an era of growing
            awareness of personal <b>privacy</b>, tech <b>monopoly</b> and{" "}
            <b>decentralisation</b>, communities of all sizes are affected.
          </Text>
        </Pair>
        <Pair heading="Flexible by nature">
          <Text>
            Whether you operate a product <b>support forum</b>, a{" "}
            <b>gaming community</b> or the next big cryptocurrency, web3 or DAO
            project, the <b>people</b> at the centre of whatever you’re doing
            deserve a platform that fades into the background and brings{" "}
            <b>what matters</b> front & centre
          </Text>
        </Pair>
        <Pair heading="Accessible by design">
          <Text>
            And by people, that means{" "}
            <em>
              <b>all</b>
            </em>{" "}
            people. A truly welcoming, inclusive and <b>diverse</b> community
            warrants a platform that takes accessibility seriously,{" "}
            <b>no questions asked</b>.
          </Text>
        </Pair>
        <Flex
          flexWrap="wrap"
          bgColor="hsla(140, 16%, 88%, 1)"
          px={{ base: 6, md: 9, lg: 12, xl: 12 }}
          pb={12}
          alignItems="center"
          justifyContent="center"
          gap={8}
        >
          <Feature
            image="/web3.webp"
            alt=""
            heading="Web3"
            body="Love it or hate it, it’s here and it’s staying. So we embrace the new web and provide features such as wallet based login, NFT avatars and more for web3 communities."
          />

          <Feature
            pt={24}
            image="/opensource.webp"
            alt=""
            heading="Open source"
            body="The benefits of open source software are impossible to ignore. When it comes to the security, development velocity, and ability to report issues, this is the way forward."
          />
        </Flex>
        <Pair big heading="For dev-ops heroes">
          <Text>
            Simple <b>deployment</b>, simple <b>maintenance</b> and simple{" "}
            <b>updates</b>. That’s what matters when you’re self-hosting, so you
            can spend time where it brings the most value.
          </Text>
          <Text>
            The choice of technologies behind the Storyden platform are all{" "}
            <b>meticulously</b> intentional to fit those values of simplicity in
            order to <b>get out of the way</b>.
          </Text>
        </Pair>
        <Pair heading="Container first">
          <Text>
            Zero <b>installation</b> steps. Set some environment variables and
            spin up a container image. Behaving like all other modern
            server-side software is the key to <b>simplicity</b>.
          </Text>
        </Pair>
        <Pair heading="Bring your own frontend">
          <Text>
            Not a fan of themes? That’s fine, <b>headless</b> mode is for you.
            Storyden is, at the core, a powerful API service with which you can{" "}
            <b>wire up</b> anything you want.
          </Text>
          <Text>
            From web to mobile apps and everything in between, the{" "}
            <b>OpenAPI</b> specification provides a fast integration path.
          </Text>
        </Pair>
        <Flex
          flexWrap="wrap"
          px={{ base: 6, md: 9, lg: 12, xl: 12 }}
          pb={12}
          alignItems="center"
          justifyContent="center"
          gap={8}
        >
          <Feature
            image="/extensible.webp"
            alt=""
            heading="Extensible"
            body="A fully documented OpenAPI schema means that you can extend the platform with plugins or even build a whole new frontend from scratch if you want to!"
          />

          <Feature
            pt={24}
            image="/builttolast.webp"
            alt=""
            heading="Built to last"
            body="Harnessing the power of technology that’s just-modern-enough helps balance stability with longevity. Storyden uses a carefully chosen toolbox with this in mind."
          />
        </Flex>
        <Pair big heading="For you and your friends">
          <Text>
            Optimised for <b>humans</b>, ready for the web <b>renaissance</b>.
            Storyden is built to be a stable foundation for the future decades
            of internet citizens and the <b>networks</b> they build.
          </Text>
        </Pair>
        <Pair heading="Sign-in your way">
          <Text>
            Rough relationship with <b>email</b>? Just don’t enable it then.
            Sign in with <b>Passkey</b>, WebAuthn, Web3 <b>wallet</b>, or choose
            from a variety of popular OAuth2 and <b>SSO</b> providers.
          </Text>
        </Pair>
        <Pair heading="Take part">
          <Text>
            Staying <b>closed source</b> is pointless in today’s internet. Fork
            it, hack on it, provide hosting, use as a basis for other apps,{" "}
            <b>contribute</b> back to the community. Not sure where to start?
            Use Storyden to <b>learn</b> about building!
          </Text>
        </Pair>
      </VStack>
    </VStack>
  );
}

function Pair({ big, heading, children }) {
  return (
    <Flex
      px={{ base: 6, md: 9, lg: 12, xl: 12 }}
      flexDir={{ base: "column", md: "row" }}
      alignItems="flex-start"
      gap={6}
      mt={big ? 10 : 2}
    >
      <Heading
        width={{ base: "100%", sm: "66%", xl: "50%" }}
        textAlign={{ base: "left", md: "right" }}
        // whiteSpace="nowrap"
        fontWeight="black"
        fontSize={
          big
            ? {
                base: "3xl",
                md: "2xl",
                lg: "2xl",
                xl: "3xl",
                "2xl": "4xl",
              }
            : {
                base: "xl",
                md: "lg",
                lg: "lg",
                xl: "xl",
                "2xl": "2xl",
              }
        }
      >
        {heading}
      </Heading>
      <Box width={{ base: "100%", sm: "67%", xl: "50%" }} mt={1}>
        {children}
      </Box>
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
      <Heading fontWeight="bold" fontSize={{ base: "2xl", lg: "4xl" }}>
        Interested?
      </Heading>
      <Text>
        Storyden is early in development and is looking for <b>feedback</b> and{" "}
        <b>contributors</b>!
      </Text>
      <Text>
        If you have <b>opinions</b> about forum software, please click{" "}
        <Link
          isExternal
          href="https://airtable.com/shrLY0jDp9CuXPB2X"
          color="hsla(265, 56%, 42%, 1)"
        >
          this link!
        </Link>
      </Text>
      <Text>
        If you know <b>Golang</b> or <b>React.js</b> and are interested in
        contributing to a <b>high-quality</b> open source project, please click{" "}
        <Link
          isExternal
          href="https://github.com/Southclaws/storyden"
          color="hsla(265, 56%, 42%, 1)"
        >
          this link!
        </Link>
      </Text>
    </VStack>
  );
}

function Footer() {
  return (
    <Flex
      flexDir="column"
      bgColor="hsla(140, 16%, 88%, 1)"
      px={{ base: 12, md: 12, lg: 24, xl: 48, "2xl": 80 }}
      py={{ base: 12, lg: 12 }}
      alignItems="center"
      justifyContent="center"
      gap={2}
      flex="1"
    >
      <Image
        src="/mark.png"
        alt="The Storyden logomark and wordmark"
        width={150}
        height={50}
      />

      <Link href="https://twitter.com/Southclaws">Twitter</Link>
      <Link href="https://github.com/Southclaws/storyden">GitHub</Link>

      <svg
        width="24"
        height="25"
        viewBox="0 0 24 25"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M12 1.95117C9.10051 1.95117 6.75 4.30168 6.75 7.20117V10.2012C5.09315 10.2012 3.75 11.5443 3.75 13.2012V19.9512C3.75 21.608 5.09315 22.9512 6.75 22.9512H17.25C18.9069 22.9512 20.25 21.608 20.25 19.9512V13.2012C20.25 11.5443 18.9069 10.2012 17.25 10.2012V7.20117C17.25 4.30168 14.8995 1.95117 12 1.95117ZM15.75 10.2012V7.20117C15.75 5.1301 14.0711 3.45117 12 3.45117C9.92893 3.45117 8.25 5.1301 8.25 7.20117V10.2012H15.75Z"
          fill="#0F172A"
        />
      </svg>
      <Text>Storyden brand, logo and other assets &copy; Barnaby Keene</Text>
    </Flex>
  );
}

export default function Home() {
  return (
    <Box>
      <Hero />
      <Story />
      <Why />
      <Features />
      <CTA />
      <Footer />
    </Box>
  );
}
