import type { WebMCPToolName } from "./tools.generated";

export async function waitForWebMCPTool(
  expectedTool: WebMCPToolName,
  unavailableMessage: string,
) {
  const modelContext = document.modelContext;
  if (!modelContext) {
    throw new Error("WebMCP is unavailable in this browser.");
  }

  const isRegistered = async () =>
    (await modelContext.getTools()).some((tool) => tool.name === expectedTool);

  if (await isRegistered()) {
    return;
  }

  await new Promise<void>((resolve, reject) => {
    const timeout = window.setTimeout(() => {
      cleanup();
      reject(new Error(unavailableMessage));
    }, 2_000);

    const handleToolChange = () => {
      void isRegistered()
        .then((ready) => {
          if (!ready) return;
          cleanup();
          resolve();
        })
        .catch((error) => {
          cleanup();
          reject(error);
        });
    };

    const cleanup = () => {
      window.clearTimeout(timeout);
      modelContext.removeEventListener("toolchange", handleToolChange);
    };

    modelContext.addEventListener("toolchange", handleToolChange);
    handleToolChange();
  });
}
