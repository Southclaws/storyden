import { RxDatabase, addRxPlugin, createRxDatabase } from "rxdb";
import { getRxStorageDexie } from "rxdb/plugins/storage-dexie";

let dbInstance: RxDatabase | null = null;

export async function getLocalStore() {
  if (dbInstance) {
    return dbInstance;
  }

  const db = await createRxDatabase({
    name: "storyden",
    storage: getRxStorageDexie(),
  });

  await db.addCollections({
    nodes: {
      schema: {
        title: "nodes",
        version: 0,
        primaryKey: "id",
        type: "object",
        properties: {
          id: { type: "string" },
          name: { type: "string" },
          content: { type: "string" },
          tags: { type: "array", items: { type: "string" } },
        },
        required: ["id", "name"],
      },
    },
    threads: {
      schema: {
        title: "threads",
        version: 0,
        primaryKey: "id",
        type: "object",
        properties: {
          id: { type: "string" },
          title: { type: "string" },
          nodeId: { type: "string" }, // Reference to a node
          messages: { type: "array", items: { type: "string" } },
        },
        required: ["id", "title", "nodeId"],
      },
    },
  });

  dbInstance = db;
  return db;
}
