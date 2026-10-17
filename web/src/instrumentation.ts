export async function register() {
  try {
    const { serverEnvironment } = await import("./config");

    serverEnvironment();
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
}
