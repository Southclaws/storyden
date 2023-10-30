import { LoginScreen } from "src/screens/auth/LoginScreen/LoginScreen";
import { getInfo } from "src/utils/info";

export default function Page() {
  return <LoginScreen />;
}

export async function generateMetadata() {
  const info = await getInfo();
  return {
    title: `Login to ${info.title}`,
    description: `Log in or sign up to ${info.title} - powered by Storyden`,
  };
}
