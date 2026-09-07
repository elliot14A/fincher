import { style } from '@vanilla-extract/css'
import { type RecipeVariants, recipe } from '@vanilla-extract/recipes'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const sidebarPanel = recipe({
  base: {
    backgroundColor: vars.color.surface,
    borderLeft: `1px solid ${vars.color.borderSubtle}`,
    display: 'flex',
    flexDirection: 'column',
    flexShrink: 0,
    height: '100%',
    minHeight: 0,
    overflowY: 'auto',
  },
  variants: {
    size: {
      compact: { width: '380px' },
      standard: { width: '400px' },
      wide: { width: '540px' },
    },
  },
  defaultVariants: {
    size: 'standard',
  },
})

export type SidebarPanelVariants = RecipeVariants<typeof sidebarPanel>

export const header = style({
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: vars.space.md,
  padding: `${vars.space.sm} ${vars.space.lg}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  flexShrink: 0,
})

export const headerMain = style({
  display: 'flex',
  gap: vars.space.md,
  alignItems: 'center',
  minWidth: 0,
})

export const entityAvatar = style({
  width: '44px',
  height: '44px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.primary,
  flexShrink: 0,
  objectFit: 'cover',
})

export const titleStack = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '3px',
  minWidth: 0,
})

export const sidebarTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
  lineHeight: vars.lineHeight.tight,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const sidebarMeta = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textSecondary,
})

export const closeBtn = style({
  background: 'transparent',
  border: 'none',
  padding: vars.space.xs,
  color: vars.color.textTertiary,
  cursor: 'pointer',
  borderRadius: vars.radii.xs,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  transition: 'color 0.1s ease, background-color 0.1s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const scrollArea = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  minHeight: 0,
  overflowY: 'auto',
  padding: vars.space.lg,
  gap: vars.space.lg,
})

export const section = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const sectionTitle = style({
  fontSize: vars.fontSize['2xs'],
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
  color: vars.color.textTertiary,
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  margin: 0,
})

export const statGrid = recipe({
  base: {
    display: 'grid',
    gap: vars.space.xs,
  },
  variants: {
    columns: {
      2: { gridTemplateColumns: 'repeat(2, 1fr)' },
      3: { gridTemplateColumns: 'repeat(3, 1fr)' },
      4: { gridTemplateColumns: 'repeat(4, 1fr)' },
    },
  },
  defaultVariants: {
    columns: 3,
  },
})

export type StatGridVariants = RecipeVariants<typeof statGrid>

export const statCard = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  gap: '2px',
})

export const statValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const statValueDanger = style({
  color: vars.color.danger,
})

export const statValueSuccess = style({
  color: vars.color.success,
})

export const statValueWarning = style({
  color: vars.color.warning,
})

export const statLabel = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
})

export const badgeGroup = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.xs,
})

export const capabilityTag = style({
  fontSize: vars.fontSize['2xs'],
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textSecondary,
  fontWeight: 500,
})

export const marketTag = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textPrimary,
  fontWeight: 600,
})

export const itemList = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
})

export const itemRow = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  width: '100%',
  textAlign: 'left',
  background: 'none',
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  font: 'inherit',
  color: 'inherit',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  backgroundColor: vars.color.surfaceElevated,
  fontSize: vars.fontSize.xs,
  cursor: 'pointer',
  transition: 'all 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    borderColor: vars.color.borderStrong,
  },
})

export const itemLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  minWidth: 0,
})

export const itemCode = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const itemDetail = style({
  color: vars.color.textTertiary,
  fontSize: vars.fontSize['2xs'],
})

export const itemRight = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  flexShrink: 0,
})

export const itemChevron = style({
  color: vars.color.textTertiary,
  opacity: 0.6,
  transition: 'transform 0.15s ease, opacity 0.15s ease, color 0.15s ease',
  selectors: {
    [`${itemRow}:hover &`]: {
      opacity: 1,
      transform: 'translateX(2px)',
      color: vars.color.primary,
    },
  },
})

export const metadataCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  padding: vars.space.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
})

export const metadataRow = style({
  display: 'flex',
  justifyContent: 'space-between',
  fontSize: vars.fontSize.xs,
})

export const metadataKey = style({
  color: vars.color.textTertiary,
})

export const metadataVal = style({
  color: vars.color.textSecondary,
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['2xs'],
})

export const emptyNotice = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  fontStyle: 'italic',
  padding: vars.space.xs,
})

export const footer = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `${vars.space.sm} ${vars.space.lg}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surface,
  flexShrink: 0,
})
