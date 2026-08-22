import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { FormControl } from "@/components/ui/form-control";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/ui/page-header";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { Grid, HStack, LStack, styled } from "@/styled-system/jsx";

import { PageLayout, type PageLayoutWidth } from ".";

const meta = {
  title: "Foundations/Layout/Page Width",
  component: PageLayout,
  args: {
    width: "content",
  },
  argTypes: {
    width: {
      control: "select",
      options: ["content", "wide", "full"],
    },
  },
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "Selects the workspace width for the current page without changing the global shell. `content` is the default for reading, forms, feeds, and settings; `wide` is for dashboards, tables, Robots, and multi-column tools; `full` is reserved for canvas-like editors and data tools. Prose remains constrained inside wide layouts. The shell retains a bottom gutter after long page content so final form actions do not touch the viewport edge.",
      },
    },
  },
} satisfies Meta<typeof PageLayout>;

export default meta;

type Story = StoryObj<typeof meta>;

function ShellPreview({ width }: { width: PageLayoutWidth }) {
  return (
    <div className="navigation__container">
      <styled.div className="navigation__grid" minH="dvh">
        <styled.aside
          display={{ base: "none", md: "flex" }}
          gridArea="leftbar"
          minW="0"
          flexDirection="column"
          gap="4"
          paddingBlock="2"
        >
          <Text fontWeight="semibold">Storyden</Text>
          <LStack gap="1">
            <Text variant="supporting">Home</Text>
            <Text variant="supporting">Discussion</Text>
            <Text variant="supporting">Library</Text>
          </LStack>
        </styled.aside>

        <div className="navigation__main">
          <styled.main minH="full">
            <PageLayout width={width}>
              <LStack gap="5">
                <PageHeader
                  title="Community workspace"
                  description="The page chooses its workspace width; content inside still owns its readable measure."
                  actions={<Button>New view</Button>}
                />

                <Grid
                  gap="3"
                  gridTemplateColumns={{ base: "1fr", lg: "2fr 1fr" }}
                >
                  <CardBox as="section">
                    <LStack gap="3">
                      <HStack justifyContent="space-between">
                        <SectionHeading>Recent activity</SectionHeading>
                        <Badge>12 updates</Badge>
                      </HStack>
                      <styled.article maxW="layout.readable">
                        <Text>
                          A wide workspace can host tools, dashboards, and
                          structured layouts without forcing prose to span the
                          available width. Readable content remains explicitly
                          constrained where it is authored.
                        </Text>
                      </styled.article>
                    </LStack>
                  </CardBox>

                  <CardBox as="aside">
                    <LStack gap="2">
                      <SectionHeading>Supporting context</SectionHeading>
                      <Text variant="supporting">
                        This belongs to the page composition, not the
                        application shell.
                      </Text>
                    </LStack>
                  </CardBox>
                </Grid>
              </LStack>
            </PageLayout>
          </styled.main>
        </div>
      </styled.div>
    </div>
  );
}

function FormFooterPreview() {
  return (
    <div className="navigation__container">
      <styled.div className="navigation__grid" minH="dvh">
        <div className="navigation__main">
          <styled.main minH="full">
            <PageLayout width="content">
              <styled.form
                display="flex"
                flexDirection="column"
                gap="6"
                minH="dvh"
              >
                <PageHeader title="New automation" />

                <LStack gap="4" maxW="3xl">
                  {Array.from({ length: 10 }, (_, index) => (
                    <FormControl key={index}>
                      <FormLabel>Field {index + 1}</FormLabel>
                      <Input aria-label={`Field ${index + 1}`} />
                    </FormControl>
                  ))}
                </LStack>

                <HStack justifyContent="end" marginTop="auto" maxW="3xl">
                  <Button type="submit">Create</Button>
                </HStack>
              </styled.form>
            </PageLayout>
          </styled.main>
        </div>
      </styled.div>
    </div>
  );
}

export const Playground: Story = {
  render: ({ width = "content" }) => <ShellPreview width={width} />,
};

export const Content: Story = {
  args: {
    width: "content",
  },
  render: ({ width = "content" }) => <ShellPreview width={width} />,
};

export const Wide: Story = {
  args: {
    width: "wide",
  },
  render: ({ width = "wide" }) => <ShellPreview width={width} />,
};

export const Full: Story = {
  args: {
    width: "full",
  },
  render: ({ width = "full" }) => <ShellPreview width={width} />,
};

export const FormFooterSpacing: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Scroll to the final action. The navigation shell keeps a visible bottom gutter after long full-page forms.",
      },
    },
  },
  render: () => <FormFooterPreview />,
};
