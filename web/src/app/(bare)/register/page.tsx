import { RegisterScreen } from "src/screens/auth/RegisterScreen/RegisterScreen";
import { getInfo } from "src/utils/info";

export default function Page() {
  return <RegisterScreen />;
}

export async function generateMetadata() {
  const info = await getInfo();
  return {
    title: `Join the community at ${info.title}`,
    description: `Log in or sign up to ${info.title} - powered by Storyden`,
  };
}
