import { AuthScreen } from "src/screens/auth/AuthScreen";
import { getInfo } from "src/utils/info";

export default function Page() {
  return <AuthScreen method={null} />;
}

export async function generateMetadata() {
  const info = await getInfo();
  return {
    title: `Authenticate with ${info.title}`,
    description: `Log in or sign up to ${info.title} - powered by Storyden`,
  };
}
