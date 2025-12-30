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

const featureCardStyles = `
  .feature-card {
    border: 1px solid #e0e0e0;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    transition: all 0.3s ease;
    cursor: pointer;
  }
  .feature-card:hover {
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
    border-color: #4CAF50;
  }
`;

export default function LandingPage() {
  return (
    <Box minH="dvh" display="flex" flexDirection="column" style={{ background: '#ffffff' }}>
      <style>{featureCardStyles}</style>
      {/* Hero Section with Mint Gradient */}
      <Box
        position="relative"
        overflow="hidden"
        py={{ base: '16', md: '32' }}
        style={{
          background: 'linear-gradient(135deg, #e8f7f5 0%, #e0f2f0 50%, #f0fffe 100%)',
        }}
      >
        {/* Decorative Circles */}
        <Box
          position="absolute"
          top="10"
          right="8"
          w="72"
          h="72"
          rounded="full"
          style={{
            background: 'rgba(76, 175, 80, 0.15)',
            filter: 'blur(40px)',
          }}
          zIndex="base"
        />
        <Box
          position="absolute"
          bottom="20"
          left="20"
          w="96"
          h="96"
          rounded="full"
          style={{
            background: 'rgba(76, 175, 80, 0.08)',
            filter: 'blur(60px)',
          }}
          zIndex="base"
        />

        <Box position="relative" maxW="6xl" mx="auto" px={{ base: '4', md: '6' }} zIndex="docked">
          <VStack maxW="3xl" mx="auto" alignItems="center" gap="8" textAlign="center">
            {/* Badge */}
            <HStack
              px="4"
              py="2"
              rounded="full"
              style={{
                background: 'rgba(76, 175, 80, 0.12)',
                color: 'rgb(27, 94, 32)',
                fontSize: '14px',
                fontWeight: '600',
              }}
              gap="2"
              justifyContent="center"
            >
              <Sparkles width="16" height="16" />
              <span>Join 10,000+ Hair Warriors</span>
            </HStack>

            {/* Heading */}
            <styled.h1
              style={{
                fontSize: 'clamp(28px, 6vw, 60px)',
                fontWeight: '900',
                lineHeight: '1.2',
                letterSpacing: '-0.02em',
                color: '#1a1a1a',
              }}
            >
              Your Hair Recovery Journey,{' '}
              <span style={{ background: 'linear-gradient(135deg, #4CAF50 0%, #00897b 100%)', backgroundClip: 'text', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                Together
              </span>
            </styled.h1>

            {/* Description */}
            <styled.p
              style={{
                fontSize: 'clamp(16px, 2vw, 20px)',
                color: '#4a4a4a',
                fontWeight: '500',
                lineHeight: '1.6',
                maxWidth: '560px',
              }}
            >
              Join Traya's supportive community to get expert guidance, share your progress, and stay motivated throughout your transformation.
            </styled.p>

            {/* CTA Button */}
            <Box pt="2">
              <Link href="/community">
                <Button
                  size="lg"
                  px="10"
                  h="14"
                  style={{
                    background: 'linear-gradient(135deg, #4CAF50 0%, #45a049 100%)',
                    color: 'white',
                    fontWeight: '700',
                    fontSize: '16px',
                    boxShadow: '0 8px 24px rgba(76, 175, 80, 0.3)',
                    border: 'none',
                  }}
                >
                  <HStack gap="2" justifyContent="center">
                    <span>Enter Community</span>
                    <ArrowRight width="20" height="20" />
                  </HStack>
                </Button>
              </Link>
            </Box>
          </VStack>
        </Box>
      </Box>

      {/* Features Section */}
      <Box py={{ base: '16', md: '24' }} style={{ background: '#fafafa' }}>
        <Box maxW="6xl" mx="auto" px={{ base: '4', md: '6' }}>
          {/* Section Header */}
          <VStack alignItems="center" gap="4" mb={{ base: '12', md: '20' }} textAlign="center">
            <styled.h2
              style={{
                fontSize: 'clamp(24px, 4vw, 48px)',
                fontWeight: '900',
                color: '#1a1a1a',
              }}
            >
              Everything You Need to Succeed
            </styled.h2>
            <styled.p
              style={{
                fontSize: '18px',
                color: '#666666',
                fontWeight: '500',
                maxWidth: '560px',
              }}
            >
              Our community is designed to support you at every step of your hair recovery journey
            </styled.p>
          </VStack>

          {/* Features Grid */}
          <Box
            display="grid"
            gridTemplateColumns={{ base: '1fr', sm: '1fr 1fr', lg: '1fr 1fr 1fr 1fr' }}
            gap={{ base: '6', md: '8' }}
          >
            {features.map((feature) => (
              <Box
                key={feature.title}
                className="feature-card"
                bg="white"
                rounded="2xl"
                p={{ base: '6', md: '8' }}
              >
                <Box
                  w="14"
                  h="14"
                  rounded="xl"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  mb="6"
                  style={{
                    background: 'linear-gradient(135deg, #4CAF50 0%, #00897b 100%)',
                  }}
                >
                  <feature.icon width="28" height="28" style={{ color: 'white' }} />
                </Box>
                <styled.h3
                  style={{
                    fontWeight: '700',
                    fontSize: '18px',
                    color: '#1a1a1a',
                    marginBottom: '12px',
                  }}
                >
                  {feature.title}
                </styled.h3>
                <styled.p
                  style={{
                    fontSize: '14px',
                    lineHeight: '1.6',
                    color: '#666666',
                  }}
                >
                  {feature.description}
                </styled.p>
              </Box>
            ))}
          </Box>
        </Box>
      </Box>

      {/* Stats Section */}
      <Box
        py={{ base: '16', md: '24' }}
        style={{
          background: 'linear-gradient(135deg, #4CAF50 0%, #2e7d32 100%)',
          color: 'white',
        }}
      >
        <Box maxW="6xl" mx="auto" px={{ base: '4', md: '6' }}>
          <VStack alignItems="center" gap="12">
            <styled.h2
              style={{
                fontSize: 'clamp(24px, 4vw, 48px)',
                fontWeight: '900',
                textAlign: 'center',
                color: 'white',
              }}
            >
              Join a Thriving Community
            </styled.h2>

            {/* Stats Grid */}
            <Box
              display="grid"
              gridTemplateColumns={{ base: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' }}
              gap={{ base: '8', md: '12' }}
              w="full"
            >
              {[
                { value: '10K+', label: 'Active Members' },
                { value: '50K+', label: 'Posts Shared' },
                { value: '85%', label: 'See Results' },
                { value: '4.8', label: 'Community Rating' },
              ].map((stat) => (
                <VStack key={stat.label} alignItems="center" gap="2">
                  <styled.div
                    style={{
                      fontSize: 'clamp(24px, 5vw, 48px)',
                      fontWeight: '900',
                      color: 'white',
                    }}
                  >
                    {stat.value}
                  </styled.div>
                  <styled.div
                    style={{
                      fontSize: '14px',
                      fontWeight: '600',
                      color: 'rgba(255, 255, 255, 0.85)',
                    }}
                  >
                    {stat.label}
                  </styled.div>
                </VStack>
              ))}
            </Box>
          </VStack>
        </Box>
      </Box>

      {/* Final CTA Section */}
      <Box
        py={{ base: '16', md: '24' }}
        style={{
          background: 'linear-gradient(135deg, #e8f7f5 0%, #e0f2f0 100%)',
        }}
      >
        <Box maxW="6xl" mx="auto" px={{ base: '4', md: '6' }}>
          <VStack alignItems="center" gap="8" textAlign="center">
            <styled.h2
              style={{
                fontSize: 'clamp(24px, 4vw, 48px)',
                fontWeight: '900',
                color: '#1a1a1a',
              }}
            >
              Ready to Start Your Journey?
            </styled.h2>
            <styled.p
              style={{
                fontSize: '18px',
                color: '#666666',
                fontWeight: '500',
                maxWidth: '640px',
              }}
            >
              Connect with fellow warriors, get expert advice, and transform your hair health together.
            </styled.p>
            <Link href="/community">
              <Button
                size="lg"
                px="10"
                h="14"
                style={{
                  background: 'linear-gradient(135deg, #4CAF50 0%, #45a049 100%)',
                  color: 'white',
                  fontWeight: '700',
                  fontSize: '16px',
                  boxShadow: '0 8px 24px rgba(76, 175, 80, 0.3)',
                  border: 'none',
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
          borderTop: '1px solid #e0e0e0',
        }}
      >
        <Box maxW="6xl" mx="auto" px={{ base: '4', md: '6' }}>
          <HStack
            justifyContent={{ base: 'center', md: 'space-between' }}
            alignItems="center"
            gap="4"
            flexDirection={{ base: 'column', md: 'row' }}
          >
            <HStack gap="2" alignItems="center">
              <Box
                w="10"
                h="10"
                rounded="lg"
                display="flex"
                alignItems="center"
                justifyContent="center"
                fontSize="xl"
                style={{
                  background: 'linear-gradient(135deg, #4CAF50 0%, #00897b 100%)',
                }}
              >
                💚
              </Box>
              <styled.span
                style={{
                  fontWeight: '700',
                  fontSize: '18px',
                  color: '#1a1a1a',
                }}
              >
                Traya Community
              </styled.span>
            </HStack>
            <styled.p
              style={{
                fontSize: '14px',
                color: '#999999',
              }}
            >
              © 2024 Traya Health. All rights reserved.
            </styled.p>
          </HStack>
        </Box>
      </Box>
    </Box>
  );
}
