import { useNavigation } from "src/components/Navigation/useNavigation";

export function useMenu() {
  const navigation = useNavigation();

  return navigation;
}
