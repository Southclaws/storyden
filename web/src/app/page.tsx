'use client';

import Link from 'next/link';
import {
  Users,
  MessageCircle,
  TrendingUp,
  Heart,
  ArrowRight,
  Sparkles,
} from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Box, VStack, HStack, styled } from '@/styled-system/jsx';

const features = [
  {
    icon: Users,
    title: 'Cohort-Based Journey',
    description:
      'Connect with warriors at the same stage of their hair recovery journey',
  },
  {
    icon: MessageCircle,
    title: 'Expert Guidance',
    description: 'Get answers from hair experts and community managers',
  },
  {
    icon: TrendingUp,
    title: 'Track Progress',
    description: 'Share milestones, celebrate wins, and stay motivated',
  },
  {
    icon: Heart,
    title: 'Peer Support',
    description: 'Find encouragement from people who understand your journey',
  },
];

export default function LandingPage() {
  return (
    <Box minH="dvh" display="flex" flexDirection="column" style={{ background: 'hsl(240, 20%, 99%)' }}>
      {/* Hero Section */}
      <Box position="relative" overflow="hidden" py={{ base: '12', md: '20' }} bg="white">
        <Box maxW="6xl" mx="auto" px="4">
          <VStack maxW="3xl" mx="auto" alignItems="center" gap="6" textAlign="center">
            {/* Badge */}
            <HStack
              px="4"
              py="2"
              rounded="full"
              style={{ background: 'hsl(255, 50%, 90%)' }}
              fontSize="sm"
              fontWeight="medium"
              gap="2"
              justifyContent="center"
            >
              <Sparkles width="16" height="16" />
              <span>Join 10,000+ Hair Warriors</span>
            </HStack>

            {/* Heading */}
            <styled.h1 fontSize={{ base: '4xl', md: '5xl', lg: '6xl' }} fontWeight="bold" lineHeight="tight">
              Your Hair Recovery Journey, <span style={{ color: '#0090ff' }}>Together</span>
            </styled.h1>

            {/* Description */}
            <styled.p
              fontSize={{ base: 'lg', md: 'xl' }}
              maxW="2xl"
              style={{ color: 'hsl(220, 5.9%, 40%)' }}
            >
              Join Traya's supportive community to get expert guidance, share your progress, and stay
              motivated throughout your transformation.
            </styled.p>

            {/* CTA Button */}
            <Link href="/community">
              <Button
                size="lg"
                px="8"
                h="14"
                style={{
                  background: '#0090ff',
                  color: 'white',
                }}
              >
                <HStack gap="2" justifyContent="center">
                  <span>Enter Community</span>
                  <ArrowRight width="20" height="20" />
                </HStack>
              </Button>
            </Link>
          </VStack>
        </Box>
      </Box>

      {/* Features Section */}
      <Box py="20" style={{ background: 'hsl(240, 20%, 98%)' }}>
        <Box maxW="6xl" mx="auto" px="4">
          {/* Section Header */}
          <VStack alignItems="center" gap="4" mb="16" textAlign="center">
            <styled.h2 fontSize={{ base: '3xl', md: '4xl' }} fontWeight="bold">
              Everything You Need to Succeed
            </styled.h2>
            <styled.p fontSize="lg" maxW="2xl" style={{ color: 'hsl(220, 5.9%, 40%)' }}>
              Our community is designed to support you at every step of your hair recovery journey
            </styled.p>
          </VStack>

          {/* Features Grid */}
          <Box
            display="grid"
            gridTemplateColumns={{ base: '1fr', md: '1fr 1fr', lg: '1fr 1fr 1fr 1fr' }}
            gap="6"
          >
            {features.map((feature) => (
              <Box
                key={feature.title}
                bg="white"
                rounded="2xl"
                p="6"
                style={{
                  border: '1px solid hsl(240, 10.1%, 86.5%)',
                }}
              >
                <Box
                  w="12"
                  h="12"
                  rounded="xl"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  mb="4"
                  bg="blue.9"
                >
                  <feature.icon width="24" height="24" style={{ color: 'white' }} />
                </Box>
                <styled.h3 fontWeight="semibold" fontSize="lg" mb="2">
                  {feature.title}
                </styled.h3>
                <styled.p fontSize="sm" style={{ color: 'hsl(220, 5.9%, 40%)' }}>
                  {feature.description}
                </styled.p>
              </Box>
            ))}
          </Box>
        </Box>
      </Box>

      {/* Stats Section */}
      <Box py="20" bg="blue.9" color="white">
        <Box maxW="6xl" mx="auto" px="4">
          <VStack alignItems="center" gap="8">
            <styled.h2 fontSize={{ base: '2xl', md: '3xl' }} fontWeight="bold" textAlign="center" color="white">
              Join a Thriving Community
            </styled.h2>

            {/* Stats Grid */}
            <Box
              display="grid"
              gridTemplateColumns={{ base: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' }}
              gap="8"
              w="full"
            >
              {[
                { value: '10K+', label: 'Active Members' },
                { value: '50K+', label: 'Posts Shared' },
                { value: '85%', label: 'See Results' },
                { value: '4.8', label: 'Community Rating' },
              ].map((stat) => (
                <VStack key={stat.label} alignItems="center" gap="1">
                  <styled.div fontSize="3xl" fontWeight="bold" color="white">
                    {stat.value}
                  </styled.div>
                  <styled.div fontSize="sm" color="blue.1">
                    {stat.label}
                  </styled.div>
                </VStack>
              ))}
            </Box>
          </VStack>
        </Box>
      </Box>

      {/* Final CTA Section */}
      <Box py="20" style={{ background: 'hsl(240, 11.1%, 94.7%)' }}>
        <Box maxW="6xl" mx="auto" px="4">
          <VStack alignItems="center" gap="8" textAlign="center">
            <styled.h2 fontSize={{ base: '3xl', md: '4xl' }} fontWeight="bold">
              Ready to Start Your Journey?
            </styled.h2>
            <styled.p fontSize="lg" maxW="xl" style={{ color: 'hsl(220, 5.9%, 40%)' }}>
              Connect with fellow warriors, get expert advice, and transform your hair health together.
            </styled.p>
            <Link href="/community">
              <Button
                size="lg"
                px="8"
                h="14"
                style={{
                  background: '#0090ff',
                  color: 'white',
                }}
              >
                <HStack gap="2" justifyContent="center">
                  <span>Join the Community</span>
                  <ArrowRight width="20" height="20" />
                </HStack>
              </Button>
            </Link>
          </VStack>
        </Box>
      </Box>

      {/* Footer */}
      <Box
        py="8"
        mt="auto"
        bg="white"
        style={{
          borderTop: '1px solid hsl(240, 10.1%, 86.5%)',
        }}
      >
        <Box maxW="6xl" mx="auto" px="4">
          <HStack
            justifyContent={{ base: 'center', md: 'space-between' }}
            alignItems="center"
            gap="4"
            flexDirection={{ base: 'column', md: 'row' }}
          >
            <HStack gap="2" alignItems="center">
              <Box
                w="8"
                h="8"
                rounded="lg"
                display="flex"
                alignItems="center"
                justifyContent="center"
                fontSize="lg"
                bg="blue.9"
              >
                💚
              </Box>
              <styled.span fontWeight="bold">Traya Community</styled.span>
            </HStack>
            <styled.p fontSize="sm" style={{ color: 'hsl(220, 5.9%, 40%)' }}>
              © 2024 Traya Health. All rights reserved.
            </styled.p>
          </HStack>
        </Box>
      </Box>
    </Box>
  );
}
