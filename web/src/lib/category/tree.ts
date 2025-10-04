import { Category } from "@/api/openapi-schema";

export type CategoryTree = Category & {
  children: CategoryTree[];
};

export function buildCategoryTree(categories: Category[]): CategoryTree[] {
  const nodes = new Map<string, CategoryTree>();
  const roots: CategoryTree[] = [];

  categories.forEach((category) => {
    nodes.set(category.id, { ...category, children: [] });
  });

  categories.forEach((category) => {
    const node = nodes.get(category.id);
    if (!node) return;

    const parentId = category.parent;

    if (parentId && nodes.has(parentId)) {
      nodes.get(parentId)?.children.push(node);
    } else {
      roots.push(node);
    }
  });

  const sortTree = (items: CategoryTree[]) => {
    items.sort((a, b) => a.sort - b.sort);
    items.forEach((child) => sortTree(child.children));
  };

  sortTree(roots);

  return roots;
}

export function isDescendant(
  nodes: CategoryTree[],
  ancestorId: string,
  descendantId: string,
): boolean {
  const ancestor = findNode(nodes, ancestorId);
  if (!ancestor) {
    return false;
  }

  return containsNode(ancestor, descendantId);
}

function findNode(nodes: CategoryTree[], id: string): CategoryTree | null {
  for (const node of nodes) {
    if (node.id === id) {
      return node;
    }

    const found = findNode(node.children, id);
    if (found) {
      return found;
    }
  }

  return null;
}

function containsNode(node: CategoryTree, id: string): boolean {
  for (const child of node.children) {
    if (child.id === id) {
      return true;
    }

    if (containsNode(child, id)) {
      return true;
    }
  }

  return false;
}
