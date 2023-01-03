import Link from "../site/Link";
import { StorydenLogo } from "../StorydenLogo";

export function HomeLink() {
  return (
    <Link href="/" _hover={{ cursor: "pointer" }}>
      <StorydenLogo />
    </Link>
  );
}
