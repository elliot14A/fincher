import { keyframes, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

const haloAnimation = keyframes({
  '0%': { opacity: '0.5', transform: 'scale(1)' },
  '80%, 100%': { opacity: '0', transform: 'scale(2.5)' },
})

const floatGlow = keyframes({
  '0%, 100%': { transform: 'translateY(0px) scale(1)' },
  '50%': { transform: 'translateY(-12px) scale(1.06)' },
})

export const landingContainer = style({
  minHeight: '100vh',
  width: '100%',
  backgroundColor: vars.color.background,
  color: vars.color.textPrimary,
  display: 'flex',
  flexDirection: 'column',
  position: 'relative',
})

// Sticky Top Navigation - Gaur / Apple Ergonomics
export const navHeader = style({
  position: 'sticky',
  top: 0,
  zIndex: 50,
  width: '100%',
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: 'rgba(10, 10, 12, 0.85)',
  backdropFilter: 'blur(16px)',
  WebkitBackdropFilter: 'blur(16px)',
})

export const navInner = style({
  maxWidth: '1360px',
  width: '100%',
  margin: '0 auto',
  height: '64px',
  padding: `0 ${vars.space['2xl']}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
})

export const navLogoArea = style({
  display: 'flex',
  alignItems: 'center',
  textDecoration: 'none',
  color: 'inherit',
})

export const navLinks = style({
  display: 'none',
  alignItems: 'center',
  gap: vars.space.xl,
  '@media': {
    '(min-width: 768px)': {
      display: 'flex',
    },
  },
})

export const navLinkItem = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textSecondary,
  textDecoration: 'none',
  fontWeight: 500,
  transition: 'color 0.15s ease',
  cursor: 'pointer',
  ':hover': {
    color: vars.color.textPrimary,
  },
})

export const navRight = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
})

export const navCtaBtn = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  height: '36px',
  padding: `0 ${vars.space.xl}`,
  borderRadius: vars.radii.sm,
  backgroundColor: vars.color.primary,
  color: vars.color.textInverse,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  textDecoration: 'none',
  transition: 'all 0.15s ease',
  boxShadow: '0 2px 8px rgba(94, 106, 210, 0.3)',
  ':hover': {
    backgroundColor: vars.color.primaryHover,
    transform: 'translateY(-1px)',
    boxShadow: '0 4px 14px rgba(94, 106, 210, 0.45)',
  },
})

// Hero Section
export const heroSection = style({
  position: 'relative',
  paddingTop: '96px',
  paddingBottom: '80px',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  textAlign: 'center',
  width: '100%',
})

export const heroGlow = style({
  position: 'absolute',
  top: '-140px',
  left: '50%',
  transform: 'translateX(-50%)',
  width: '800px',
  height: '460px',
  borderRadius: vars.radii.full,
  background: `radial-gradient(ellipse at center, ${vars.color.primaryMuted} 0%, rgba(94, 106, 210, 0.05) 50%, transparent 75%)`,
  filter: 'blur(60px)',
  pointerEvents: 'none',
  zIndex: 0,
  animation: `${floatGlow} 9s ease-in-out infinite`,
})

export const heroContainer = style({
  maxWidth: '1080px',
  width: '100%',
  margin: '0 auto',
  padding: `0 ${vars.space.xl}`,
  position: 'relative',
  zIndex: 1,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
})

export const heroBadge = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: '6px 14px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderStrong}`,
  fontSize: vars.fontSize.xs,
  fontFamily: fonts.mono,
  color: vars.color.textPrimary,
  marginBottom: vars.space.xl,
  boxShadow: '0 2px 10px rgba(0,0,0,0.2)',
})

export const heroBadgeDot = style({
  position: 'relative',
  width: '6px',
  height: '6px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.success,
  flexShrink: 0,
})

export const heroBadgePulse = style({
  position: 'absolute',
  inset: 0,
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.success,
  animation: `${haloAnimation} 2s ease-out infinite`,
})

export const heroTitle = style({
  fontSize: '2.75rem',
  lineHeight: '1.05',
  fontWeight: 700,
  letterSpacing: '-0.04em',
  color: vars.color.textPrimary,
  margin: `0 0 ${vars.space.lg} 0`,
  '@media': {
    '(min-width: 640px)': {
      fontSize: '3.75rem',
    },
    '(min-width: 1024px)': {
      fontSize: '4.75rem',
    },
  },
})

export const heroTitleGradient = style({
  background: `linear-gradient(135deg, ${vars.color.textPrimary} 35%, ${vars.color.primary} 100%)`,
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent',
})

export const heroSubtitle = style({
  fontSize: vars.fontSize.lg,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  maxWidth: '740px',
  margin: `0 0 ${vars.space['2xl']} 0`,
  '@media': {
    '(min-width: 768px)': {
      fontSize: '1.25rem',
      lineHeight: '1.5',
    },
  },
})

export const heroActions = style({
  display: 'flex',
  flexWrap: 'wrap',
  alignItems: 'center',
  justifyContent: 'center',
  gap: vars.space.md,
  marginBottom: vars.space['3xl'],
})

export const primaryCtaBtn = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.sm,
  height: '46px',
  padding: `0 ${vars.space['2xl']}`,
  borderRadius: vars.radii.sm,
  backgroundColor: vars.color.primary,
  color: vars.color.textInverse,
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  textDecoration: 'none',
  transition: 'all 0.15s ease',
  boxShadow: '0 2px 14px rgba(94, 106, 210, 0.4)',
  ':hover': {
    backgroundColor: vars.color.primaryHover,
    transform: 'translateY(-1px)',
    boxShadow: '0 4px 20px rgba(94, 106, 210, 0.55)',
  },
})

export const secondaryCtaBtn = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.sm,
  height: '46px',
  padding: `0 ${vars.space.xl}`,
  borderRadius: vars.radii.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderStrong}`,
  color: vars.color.textPrimary,
  fontSize: vars.fontSize.base,
  fontWeight: 500,
  textDecoration: 'none',
  transition: 'all 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    borderColor: vars.color.textSecondary,
  },
})

// Architecture Visual Graphic in Hero
export const architectureCard = style({
  width: '100%',
  maxWidth: '1080px',
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.md,
  boxShadow: '0 28px 70px -15px rgba(0, 0, 0, 0.75)',
  overflow: 'hidden',
  textAlign: 'left',
  marginTop: vars.space.md,
})

export const architectureCardHeader = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: `${vars.space.sm} ${vars.space.lg}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surfaceElevated,
})

export const archHeaderLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const archDot = style({
  width: '9px',
  height: '9px',
  borderRadius: vars.radii.full,
})

export const dotRed = style({
  backgroundColor: vars.color.danger,
})

export const dotYellow = style({
  backgroundColor: vars.color.warning,
})

export const dotGreen = style({
  backgroundColor: vars.color.success,
})

export const archTitle = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textSecondary,
  marginLeft: vars.space.xs,
})

export const archLivePill = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.success,
  backgroundColor: vars.color.successMuted,
  padding: '2px 8px',
  borderRadius: vars.radii.full,
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
})

export const archGrid = style({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: vars.space.lg,
  padding: vars.space.xl,
  '@media': {
    '(min-width: 860px)': {
      gridTemplateColumns: '1.2fr auto 1.4fr',
      alignItems: 'center',
    },
  },
})

export const archEventBox = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  backgroundColor: vars.color.background,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  padding: vars.space.md,
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
})

export const archEventBadge = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  color: vars.color.warning,
  fontWeight: 600,
  fontSize: vars.fontSize['2xs'],
  marginBottom: vars.space['2xs'],
})

export const archCodePre = style({
  margin: 0,
  lineHeight: '1.5',
  color: vars.color.textPrimary,
  overflowX: 'auto',
})

export const archConnectorCol = style({
  display: 'none',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  gap: vars.space.xs,
  color: vars.color.primary,
  '@media': {
    '(min-width: 860px)': {
      display: 'flex',
    },
  },
})

export const archConnectorText = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
})

export const archNodeList = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
})

export const archNodeCard = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: `${vars.space.sm} ${vars.space.md}`,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  fontSize: vars.fontSize.xs,
  gap: vars.space.md,
})

export const archNodeLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  minWidth: 0,
})

export const archNodeIcon = style({
  color: vars.color.primary,
  flexShrink: 0,
})

export const archNodeText = style({
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const archNodeSubtext = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
})

// General Section Styling
export const sectionContainer = style({
  maxWidth: '1240px',
  margin: '0 auto',
  padding: `${vars.space['3xl']} ${vars.space.xl}`,
  width: '100%',
})

export const sectionHeader = style({
  maxWidth: '740px',
  marginBottom: vars.space['3xl'],
})

export const sectionBadge = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize['2xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  letterSpacing: '0.1em',
  color: vars.color.primary,
  marginBottom: vars.space.sm,
})

export const sectionHeading = style({
  fontSize: '2.1rem',
  fontWeight: 700,
  letterSpacing: '-0.03em',
  lineHeight: vars.lineHeight.tight,
  color: vars.color.textPrimary,
  margin: `0 0 ${vars.space.sm} 0`,
  '@media': {
    '(min-width: 768px)': {
      fontSize: '2.75rem',
    },
  },
})

export const sectionSubheading = style({
  fontSize: vars.fontSize.base,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  margin: 0,
})

// Problem Cards Grid
export const problemGrid = style({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: vars.space.lg,
  '@media': {
    '(min-width: 768px)': {
      gridTemplateColumns: 'repeat(3, 1fr)',
    },
  },
})

export const problemCard = style({
  display: 'flex',
  flexDirection: 'column',
  padding: vars.space.xl,
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  gap: vars.space.md,
  transition: 'border-color 0.2s ease, transform 0.2s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
    transform: 'translateY(-2px)',
  },
})

export const problemIconBox = style({
  width: '40px',
  height: '40px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.dangerMuted,
  border: `1px solid ${vars.color.dangerBorder}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.danger,
})

export const problemTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const problemText = style({
  fontSize: vars.fontSize.sm,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  margin: 0,
})

// 4 Pillars Feature Grid
export const featureGrid = style({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: vars.space.xl,
  '@media': {
    '(min-width: 768px)': {
      gridTemplateColumns: 'repeat(2, 1fr)',
    },
  },
})

export const featureCard = style({
  display: 'flex',
  flexDirection: 'column',
  padding: vars.space['2xl'],
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.md,
  gap: vars.space.lg,
  position: 'relative',
  overflow: 'hidden',
  transition: 'border-color 0.2s ease, transform 0.2s ease',
  ':hover': {
    borderColor: vars.color.primaryBorder,
    transform: 'translateY(-2px)',
  },
})

export const featureIconBox = style({
  width: '44px',
  height: '44px',
  borderRadius: vars.radii.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderStrong}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.primary,
})

export const featureTitle = style({
  fontSize: vars.fontSize.lg,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const featureText = style({
  fontSize: vars.fontSize.sm,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  margin: 0,
})

export const featurePillList = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.xs,
  marginTop: 'auto',
})

export const featurePill = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  padding: '3px 8px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textSecondary,
})

// Conversational AI Intelligence Showcase Section
export const chatShowcaseContainer = style({
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.lg,
  overflow: 'hidden',
  display: 'grid',
  gridTemplateColumns: '1fr',
  boxShadow: '0 24px 60px -15px rgba(0,0,0,0.65)',
  '@media': {
    '(min-width: 900px)': {
      gridTemplateColumns: '1fr 1.35fr',
    },
  },
})

export const chatPromptsCol = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space['2xl']} ${vars.space.xl}`,
  gap: vars.space.xs,
  borderRight: 'none',
  backgroundColor: 'transparent',
  '@media': {
    '(min-width: 900px)': {
      borderRight: `1px solid ${vars.color.borderSubtle}`,
    },
  },
})

export const chatPromptsHeader = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
  marginBottom: vars.space.md,
  padding: `0 ${vars.space.xs}`,
})

export const chatPromptsTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const chatPromptsSubtitle = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  margin: 0,
  lineHeight: vars.lineHeight.normal,
})

export const chatPromptButton = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.md,
  padding: `${vars.space.md} ${vars.space.md}`,
  backgroundColor: 'transparent',
  border: 'none',
  borderLeft: '3px solid transparent',
  borderRadius: vars.radii.xs,
  textAlign: 'left',
  cursor: 'pointer',
  transition: 'all 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
  },
})

export const chatPromptButtonActive = style({
  backgroundColor: vars.color.surfaceHover,
  borderLeft: `3px solid ${vars.color.primary}`,
})

export const chatPromptLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  minWidth: 0,
})

export const chatPromptIcon = style({
  color: vars.color.primary,
  flexShrink: 0,
})

export const chatPromptText = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textPrimary,
  lineHeight: vars.lineHeight.snug,
})

export const chatPromptTag = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  flexShrink: 0,
})

export const chatWindowCard = style({
  display: 'flex',
  flexDirection: 'column',
  backgroundColor: 'transparent',
  padding: `${vars.space['2xl']} ${vars.space['2xl']}`,
  gap: vars.space.lg,
  minHeight: '400px',
})

export const chatWindowHeader = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  paddingBottom: vars.space.md,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
})

export const chatWindowHeaderLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
})

export const chatWindowAvatar = style({
  width: '26px',
  height: '26px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.primaryMuted,
  color: vars.color.primary,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
})

export const chatWindowTitle = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const chatWindowStatus = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.success,
  display: 'flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
})

export const chatWindowBody = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xl,
  flex: 1,
})

export const chatUserBubble = style({
  alignSelf: 'flex-end',
  maxWidth: '85%',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: `${vars.radii.sm} ${vars.radii.sm} 0 ${vars.radii.sm}`,
  padding: `${vars.space.sm} ${vars.space.lg}`,
  fontSize: vars.fontSize.xs,
  color: vars.color.textPrimary,
  lineHeight: vars.lineHeight.normal,
})

export const chatAiBubble = style({
  alignSelf: 'flex-start',
  maxWidth: '100%',
  backgroundColor: 'transparent',
  border: 'none',
  padding: 0,
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.md,
})

export const chatAiAnswerText = style({
  fontSize: vars.fontSize.xs,
  lineHeight: '1.65',
  color: vars.color.textPrimary,
  margin: 0,
})

export const chatCitationsWrapper = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  paddingTop: vars.space.md,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
})

export const chatCitationsLabel = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  color: vars.color.textTertiary,
  letterSpacing: '0.06em',
})

export const chatCitationList = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.xs,
})

export const chatCitationPill = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.teal,
  backgroundColor: vars.color.tealMuted,
  padding: '4px 10px',
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.tealBorder}`,
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const chatFooterActionRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  marginTop: 'auto',
  paddingTop: vars.space.md,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
})

export const chatFooterNote = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
})

export const chatGoToChatLink = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize.xs,
  color: vars.color.primary,
  fontWeight: 600,
  textDecoration: 'none',
  ':hover': {
    textDecoration: 'underline',
  },
})

// Interactive Live Playground Preview
export const demoContainer = style({
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.lg,
  overflow: 'hidden',
  boxShadow: '0 24px 60px -15px rgba(0,0,0,0.65)',
})

export const demoNav = style({
  display: 'flex',
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: 'transparent',
  padding: `${vars.space.xs} ${vars.space.xl} 0 ${vars.space.xl}`,
  overflowX: 'auto',
  gap: vars.space.xs,
})

export const demoTabBtn = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  padding: `${vars.space.md} ${vars.space.lg}`,
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textTertiary,
  background: 'none',
  border: 'none',
  borderBottom: '2px solid transparent',
  cursor: 'pointer',
  whiteSpace: 'nowrap',
  transition: 'all 0.15s ease',
  ':hover': {
    color: vars.color.textPrimary,
  },
})

export const demoTabBtnActive = style({
  color: vars.color.primary,
  fontWeight: 600,
  borderBottom: `2px solid ${vars.color.primary}`,
})

export const demoBody = style({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: vars.space['2xl'],
  padding: `${vars.space['2xl']} ${vars.space['2xl']}`,
  '@media': {
    '(min-width: 900px)': {
      gridTemplateColumns: '1.2fr 1fr',
    },
  },
})

export const demoWaterfallCol = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const demoWaterfallHeader = style({
  fontSize: vars.fontSize['2xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  color: vars.color.textTertiary,
  marginBottom: vars.space.sm,
  letterSpacing: '0.04em',
})

export const demoStepRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: `${vars.space.sm} ${vars.space.md}`,
  backgroundColor: 'transparent',
  border: 'none',
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  borderRadius: vars.radii.xs,
  fontSize: vars.fontSize.xs,
  transition: 'background-color 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
  },
})

export const demoStepLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
})

export const stepIndexText = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  width: '18px',
})

export const demoStepName = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const demoStepBadge = style({
  fontSize: vars.fontSize['3xs'],
  padding: '2px 6px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.tealMuted,
  color: vars.color.teal,
})

export const stepRightBox = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
})

export const demoStepLatency = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const stepCheckIcon = style({
  color: vars.color.success,
})

export const demoScorecardCol = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.md,
  backgroundColor: 'transparent',
  border: 'none',
  padding: 0,
})

export const demoScorecardTitle = style({
  fontSize: vars.fontSize['2xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  color: vars.color.primary,
  fontWeight: 600,
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const demoScorecardVerdict = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize.xs,
})

export const demoScorecardText = style({
  fontSize: vars.fontSize.xs,
  lineHeight: '1.65',
  color: vars.color.textSecondary,
  backgroundColor: 'transparent',
  border: 'none',
  borderLeft: `2px solid ${vars.color.primaryBorder}`,
  borderRadius: 0,
  padding: `0 0 0 ${vars.space.md}`,
  margin: 0,
})

export const executedActionsWrapper = style({
  marginTop: 'auto',
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  paddingTop: vars.space.sm,
})

export const executedActionsTitle = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
})

export const executedActionPill = style({
  fontSize: vars.fontSize.xs,
  fontFamily: fonts.mono,
  color: vars.color.textPrimary,
  backgroundColor: 'transparent',
  padding: `${vars.space['2xs']} 0`,
  border: 'none',
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const executedActionIcon = style({
  color: vars.color.primary,
  flexShrink: 0,
})

// Tech Ecosystem Strip
export const techStrip = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  gap: vars.space.md,
  padding: `${vars.space['2xl']} ${vars.space.xl}`,
  margin: `${vars.space.lg} auto 0 auto`,
  maxWidth: '1080px',
  width: '100%',
})

export const techStripLabel = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: vars.color.textTertiary,
})

export const techPillList = style({
  display: 'flex',
  flexWrap: 'wrap',
  alignItems: 'center',
  justifyContent: 'center',
  gap: vars.space.sm,
})

export const techPill = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  padding: '6px 14px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textSecondary,
  transition: 'all 0.15s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
    color: vars.color.textPrimary,
  },
})

export const techPillIcon = style({
  color: vars.color.primary,
  flexShrink: 0,
})

// Closed-Loop Event Feedback Architecture Section
export const loopSection = style({
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.lg,
  padding: `${vars.space['3xl']} ${vars.space['2xl']}`,
  margin: `${vars.space['3xl']} auto`,
  maxWidth: '1240px',
  boxShadow: '0 20px 50px -10px rgba(0,0,0,0.5)',
})

export const loopHeader = style({
  maxWidth: '800px',
  marginBottom: vars.space['2xl'],
})

export const feedbackBadge = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize['2xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  color: vars.color.teal,
  backgroundColor: vars.color.tealMuted,
  border: `1px solid ${vars.color.tealBorder}`,
  padding: '3px 10px',
  borderRadius: vars.radii.full,
  marginBottom: vars.space.sm,
})

export const loopGrid = style({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: vars.space.md,
  '@media': {
    '(min-width: 768px)': {
      gridTemplateColumns: 'repeat(5, 1fr)',
    },
  },
})

export const loopCard = style({
  display: 'flex',
  flexDirection: 'column',
  padding: vars.space.lg,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.borderSubtle}`,
  borderRadius: vars.radii.sm,
  gap: vars.space.xs,
  position: 'relative',
  transition: 'all 0.2s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
    transform: 'translateY(-2px)',
  },
})

export const loopCardActive = style({
  border: `1px solid ${vars.color.tealBorder}`,
  backgroundColor: vars.color.surfaceHover,
})

export const loopStepNum = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.primary,
  fontWeight: 700,
  letterSpacing: '0.04em',
})

export const loopStepTitle = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const loopStepDesc = style({
  fontSize: vars.fontSize['2xs'],
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  margin: 0,
})

// Deterministic Safety Guarantees
export const safetyGrid = style({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: vars.space.xl,
  marginTop: vars.space['2xl'],
  '@media': {
    '(min-width: 768px)': {
      gridTemplateColumns: 'repeat(3, 1fr)',
    },
  },
})

export const safetyCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
  padding: vars.space.xl,
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.md,
  transition: 'border-color 0.2s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
  },
})

export const safetyIconBox = style({
  width: '38px',
  height: '38px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.primary,
})

export const safetyTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const safetyDesc = style({
  fontSize: vars.fontSize.xs,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  margin: 0,
})

// Metrics Strip
export const metricsGrid = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(2, 1fr)',
  gap: vars.space.lg,
  padding: `${vars.space['2xl']} 0`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  '@media': {
    '(min-width: 768px)': {
      gridTemplateColumns: 'repeat(4, 1fr)',
    },
  },
})

export const metricCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
})

export const metricValue = style({
  fontFamily: fonts.mono,
  fontSize: '2.25rem',
  fontWeight: 700,
  color: vars.color.textPrimary,
  letterSpacing: '-0.03em',
})

export const metricLabel = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const metricSubtext = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

// Final CTA Section
export const finalCtaCard = style({
  position: 'relative',
  padding: `${vars.space['3xl']} ${vars.space['2xl']}`,
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.md,
  textAlign: 'center',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  overflow: 'hidden',
  margin: `${vars.space['3xl']} auto ${vars.space['3xl']} auto`,
  maxWidth: '1240px',
  boxShadow: '0 20px 40px rgba(0,0,0,0.5)',
})

export const finalCtaTitle = style({
  fontSize: '2rem',
  fontWeight: 700,
  letterSpacing: '-0.025em',
  color: vars.color.textPrimary,
  margin: `0 0 ${vars.space.sm} 0`,
  '@media': {
    '(min-width: 640px)': {
      fontSize: '2.5rem',
    },
  },
})

export const finalCtaSubtitle = style({
  fontSize: vars.fontSize.base,
  color: vars.color.textSecondary,
  maxWidth: '580px',
  margin: `0 0 ${vars.space.xl} 0`,
})

export const finalCtaButtons = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.md,
  justifyContent: 'center',
})
