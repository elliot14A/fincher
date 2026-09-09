import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const page = style({
  display: 'flex',
  flexDirection: 'column',
  height: '100%',
  minHeight: 0,
  backgroundColor: vars.color.background,
})

export const header = style({
  padding: `${vars.space.lg} ${vars.space.xl}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  flexShrink: 0,
})

export const pageTitle = style({
  fontSize: vars.fontSize.xl,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
})

export const pageSubtitle = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textSecondary,
})

export const body = style({
  flex: 1,
  minHeight: 0,
  overflowY: 'auto',
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) 380px',
  gap: vars.space.xl,
  padding: vars.space.xl,
  alignItems: 'start',
})

export const column = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.lg,
  minWidth: 0,
})

export const sectionLabel = style({
  fontSize: vars.fontSize['2xs'],
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
  color: vars.color.textTertiary,
  marginBottom: vars.space.xs,
})

export const scenarioGrid = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: vars.space.sm,
})

export const scenarioCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  padding: vars.space.md,
  borderRadius: vars.radii.md,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
  cursor: 'pointer',
  textAlign: 'left',
  transition: 'border-color 0.12s ease, background-color 0.12s ease',
  ':hover': {
    borderColor: vars.color.borderStrong,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const scenarioCardActive = style({
  borderColor: vars.color.primary,
  backgroundColor: vars.color.primaryMuted,
})

export const scenarioHead = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const scenarioIcon = style({
  color: vars.color.primary,
  flexShrink: 0,
})

export const scenarioName = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const scenarioDesc = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: 1.45,
})

export const triggerTag = style({
  alignSelf: 'flex-start',
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
  color: vars.color.textTertiary,
  padding: '1px 6px',
  borderRadius: vars.radii.full,
  border: `1px solid ${vars.color.border}`,
})

export const field = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
})

export const fieldLabel = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  fontWeight: 500,
})

export const input = style({
  width: '100%',
  padding: `${vars.space.xs} ${vars.space.sm}`,
  borderRadius: vars.radii.sm,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surface,
  color: vars.color.textPrimary,
  fontSize: vars.fontSize.sm,
  fontFamily: fonts.sans,
  outline: 'none',
  ':focus': {
    borderColor: vars.color.primary,
  },
})

export const severityBadge = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  alignSelf: 'flex-start',
  height: '34px',
  padding: `0 ${vars.space.sm}`,
  borderRadius: vars.radii.sm,
  border: `1px solid ${vars.color.border}`,
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  letterSpacing: '0.04em',
})

export const severityBadgeInfo = style({
  color: vars.color.textSecondary,
  backgroundColor: vars.color.surface,
})

export const severityBadgeWarn = style({
  color: vars.color.warning,
  backgroundColor: vars.color.warningMuted,
  borderColor: vars.color.warningBorder,
})

export const severityBadgeCritical = style({
  color: vars.color.danger,
  backgroundColor: vars.color.dangerMuted,
  borderColor: vars.color.dangerBorder,
})

export const fieldRow = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: vars.space.sm,
})

export const hintText = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const textarea = style([
  input,
  {
    minHeight: '220px',
    resize: 'vertical',
    fontFamily: fonts.mono,
    fontSize: vars.fontSize.xs,
    lineHeight: 1.5,
    whiteSpace: 'pre',
  },
])

export const previewCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
  padding: vars.space.md,
  borderRadius: vars.radii.md,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
  position: 'sticky',
  top: 0,
})

export const previewPre = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
  backgroundColor: vars.color.background,
  padding: vars.space.sm,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.borderSubtle}`,
  overflowX: 'auto',
  lineHeight: 1.5,
  margin: 0,
  maxHeight: '360px',
})

export const emitRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
})

export const resultCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
  padding: vars.space.md,
  borderRadius: vars.radii.md,
  border: `1px solid ${vars.color.successBorder}`,
  backgroundColor: vars.color.successMuted,
})

export const resultTitle = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.success,
})

export const resultMeta = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
})

export const runLinkList = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space['2xs'],
})

export const runLink = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  color: vars.color.primary,
  textDecoration: 'none',
  ':hover': {
    textDecoration: 'underline',
  },
})

export const noRunNote = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
})
