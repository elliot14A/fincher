import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export * from '#/components/ui/sidebar/sidebar.css'

export const posterThumb = style({
  width: '44px',
  height: '62px',
  borderRadius: vars.radii.xs,
  objectFit: 'cover',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  flexShrink: 0,
})

export const posterPlaceholder = style({
  width: '44px',
  height: '62px',
  borderRadius: vars.radii.xs,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  color: vars.color.textTertiary,
  flexShrink: 0,
})

export const titleInfo = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '3px',
  minWidth: 0,
})

export const titleName = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
  lineHeight: vars.lineHeight.tight,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const titleMeta = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textSecondary,
})

export const progressCard = style({
  padding: vars.space.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const progressHeader = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
})

export const progressLabel = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textSecondary,
})

export const progressPct = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const progressBar = style({
  height: '5px',
  borderRadius: vars.radii.full,
  backgroundColor: vars.color.surfaceHover,
  overflow: 'hidden',
  position: 'relative',
})

export const progressFill = style({
  height: '100%',
  backgroundColor: vars.color.primary,
  borderRadius: vars.radii.full,
  transition: 'width 0.25s ease',
})

export const synopsisBox = style({
  fontSize: vars.fontSize.xs,
  lineHeight: vars.lineHeight.normal,
  color: vars.color.textSecondary,
  padding: vars.space.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  margin: 0,
})

export const sectionHeaderRow = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
})

export const skeletonBlock = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  padding: vars.space.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
})

export const skeletonLine = style({
  height: '14px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceHover,
  opacity: 0.7,
})

export const skeletonLineShort = style({
  width: '60%',
})

export const skeletonLineMedium = style({
  width: '80%',
})

export const skeletonLineThin = style({
  width: '100%',
  height: '6px',
})
