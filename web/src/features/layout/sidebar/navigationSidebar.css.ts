import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const sidebarContainer = style({
  width: '244px',
  backgroundColor: vars.color.surface,
  borderRight: `1px solid ${vars.color.borderSubtle}`,
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.md} ${vars.space.sm}`,
  gap: vars.space['3xs'],
  flexShrink: 0,
  height: '100vh',
  overflowY: 'auto',
})

export const brandRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  marginBottom: vars.space.md,
})

export const brandSubtitle = style({
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.xs,
})

export const newChatLink = style({
  textDecoration: 'none',
  display: 'block',
})

export const composeButton = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: vars.space.xs,
  width: '100%',
  height: '32px',
  backgroundColor: vars.color.primary,
  border: 'none',
  borderRadius: vars.radii.sm,
  color: vars.color.textInverse,
  fontSize: vars.fontSize.sm,
  fontWeight: 500,
  cursor: 'pointer',
  marginBottom: vars.space.sm,
  transition: 'background-color 0.1s ease',
  ':hover': {
    backgroundColor: vars.color.primaryHover,
  },
})

export const searchRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  marginBottom: vars.space.md,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.sm,
  cursor: 'pointer',
  transition: 'background-color 0.1s ease, border-color 0.1s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    borderColor: vars.color.borderStrong,
  },
})

export const searchLabel = style({
  flex: 1,
})

export const kbdHint = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
  backgroundColor: vars.color.surfaceActive,
  border: `1px solid ${vars.color.borderStrong}`,
  borderRadius: vars.radii.xs,
  padding: '1px 5px',
  lineHeight: vars.lineHeight.snug,
})

export const navItem = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  height: '30px',
  color: vars.color.textSecondary,
  fontSize: vars.fontSize.sm,
  borderRadius: vars.radii.sm,
  textDecoration: 'none',
  transition: 'color 0.1s ease, background-color 0.1s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const navItemActive = style({
  backgroundColor: vars.color.surfaceHover,
  color: vars.color.textPrimary,
  fontWeight: 500,
})

export const navItemLabel = style({
  flex: 1,
})

export const sectionTitle = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textTertiary,
  padding: `${vars.space.xl} ${vars.space.sm} ${vars.space.xs} ${vars.space.sm}`,
})

export const threadItem = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.xs,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  height: '30px',
  color: vars.color.textSecondary,
  fontSize: vars.fontSize.sm,
  borderRadius: vars.radii.sm,
  cursor: 'pointer',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const threadItemLabel = style({
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const threadItemActive = style({
  backgroundColor: vars.color.surfaceElevated,
  color: vars.color.textPrimary,
  fontWeight: 500,
})

export const activeDot = style({
  width: '5px',
  height: '5px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.primary,
  flexShrink: 0,
})

export const viewAllRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.sm,
  cursor: 'pointer',
})
