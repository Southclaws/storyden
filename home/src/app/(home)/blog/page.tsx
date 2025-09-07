import { blog } from "@/lib/source";
import { css } from "@/styled-system/css";
import {
  Card,
  Center,
  Grid,
  GridItem,
  styled,
  VStack,
} from "@/styled-system/jsx";
import { formatDate, formatDistanceToNow } from "date-fns";
import Image from "next/image";

const heroImageStyles = css({
  position: "relative",
  objectFit: "cover",
  objectPosition: "top",
  maxHeight: {
    base: "64",
    md: "xs",
    lg: "sm",
  },
  height: "full",
  width: "full",
});

const logoImageStyles = css({
  height: { base: "16", sm: "20", md: "24", lg: "32" },
  width: { base: "16", sm: "20", md: "24", lg: "32" },
});

export default function Page() {
  const posts = [...blog.getPages()].sort(
    (a, b) => new Date(b.data.date).getTime() - new Date(a.data.date).getTime()
  );

  return (
    <VStack w="full">
      <Grid
        w="full"
        gridTemplateRows="1fr"
        gridTemplateColumns="1fr"
        height={{
          base: "64",
          md: "xs",
          lg: "sm",
        }}
      >
        <GridItem gridRow="1/2" gridColumn="1/2">
          <Image
            alt="Storyden mountains"
            src="/brand/Storyden mountains.png"
            width="1920"
            height="1080"
            className={heroImageStyles}
          />
        </GridItem>
        <GridItem
          zIndex="1"
          gridRow="1/2"
          gridColumn="1/2"
          background="linear-gradient(0deg, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0.0) 50%)"
        />
        <GridItem zIndex="2" gridRow="1/2" gridColumn="1/2">
          <Center height="full">
            <VStack pt={{ base: "12", sm: "10", md: "8", lg: "6" }}>
              <Image
                alt="Storyden mountains"
                src="/brand/logomark_newspaper_600.png"
                width="600"
                height="600"
                className={logoImageStyles}
              />
              <styled.h1
                bgColor="black/32"
                borderRadius={{ base: "lg", md: "xl", lg: "2xl" }}
                px="4"
                backdropBlur="md"
                backdropFilter="auto"
                color="Shades.newspaper"
                fontSize={{
                  base: "xl",
                  sm: "2xl",
                  md: "3xl",
                  lg: "4xl",
                }}
              >
                Storyden blog
              </styled.h1>
            </VStack>
          </Center>
        </GridItem>
      </Grid>

      <VStack maxW="prose" px="2">
        <Grid
          gridTemplateColumns={{
            base: "1fr",
            md: "1fr 1fr",
          }}
          gridAutoRows="1fr"
        >
          {posts.map((post) => {
            const date = formatDate(post.data.date, "yyyy-MM-dd");
            const postedAt = formatDistanceToNow(post.data.date, {
              addSuffix: true,
            });

            return (
              <Card key={post.url}>
                <VStack
                  height="full"
                  alignItems="start"
                  justifyContent="space-between"
                  display="flex"
                >
                  <VStack alignItems="start" gap="1">
                    <a href={post.url}>
                      <styled.h2 fontSize="lg">{post.data.title}</styled.h2>
                    </a>
                    <styled.p lineClamp="2">{post.data.description}</styled.p>
                  </VStack>

                  <styled.time color="slate.500" title={date}>
                    {postedAt}
                  </styled.time>
                </VStack>
              </Card>
            );
          })}
        </Grid>
      </VStack>
    </VStack>
  );
}
