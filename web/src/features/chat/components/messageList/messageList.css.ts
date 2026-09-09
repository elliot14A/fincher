import { keyframes, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const list = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xl,
  width: '100%',
})

export const row = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const userRow = style([row, { alignItems: 'flex-end' }])

export const assistantRow = style([row, { alignItems: 'flex-start' }])

export const roleLabel = style({
  fontSize: vars.fontSize['3xs'],
  fontFamily: fonts.mono,
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
  color: vars.color.textTertiary,
})

export const userBubble = style({
  backgroundColor: vars.color.primaryMuted,
  border: `1px solid ${vars.color.primaryBorder}`,
  borderRadius: vars.radii.md,
  padding: `${vars.space.sm} ${vars.space.md}`,
  fontSize: vars.fontSize.sm,
  lineHeight: vars.lineHeight.snug,
  color: vars.color.textPrimary,
  maxWidth: '80%',
  whiteSpace: 'pre-wrap',
})

export const assistantBubble = style({
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.md,
  padding: vars.space.lg,
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.md,
})

const blink = keyframes({
  '0%, 100%': { opacity: 0.25 },
  '50%': { opacity: 1 },
})

export const thinking = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize.xs,
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
})

export const thinkingDot = style({
  width: '6px',
  height: '6px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.teal,
  animation: `${blink} 1.1s ease-in-out infinite`,
})
