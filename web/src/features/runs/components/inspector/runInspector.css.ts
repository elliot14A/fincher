import { keyframes, style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

const pulseGlow = keyframes({
  '0%, 100%': { opacity: '1', transform: 'scale(1)' },
  '50%': { opacity: '0.4', transform: 'scale(0.85)' },
})

export const inspectorPanel = style({
  width: '540px',
  backgroundColor: vars.color.surface,
  borderLeft: `1px solid ${vars.color.borderSubtle}`,
  display: 'flex',
  flexDirection: 'column',
  flexShrink: 0,
  height: '100%',
  minHeight: 0,
  overflowY: 'auto',
})

export const header = style({
  display: 'flex',
  flexDirection: 'column',
  padding: `${vars.space.md} ${vars.space.lg}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  gap: vars.space.sm,
  flexShrink: 0,
  backgroundColor: vars.color.surfaceElevated,
})

export const headerTopRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.md,
})

export const headerTopRight = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  flexShrink: 0,
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

export const titleStack = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '2px',
  minWidth: 0,
})

export const runTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const runIdRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const runIdText = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const copyBtn = style({
  background: 'none',
  border: 'none',
  padding: '2px 4px',
  color: vars.color.textTertiary,
  cursor: 'pointer',
  borderRadius: vars.radii.xs,
  display: 'inline-flex',
  alignItems: 'center',
  gap: '3px',
  fontSize: vars.fontSize['3xs'],
  transition: 'all 0.12s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
  },
})

export const statBar = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(3, 1fr)',
  gap: vars.space.xs,
  paddingTop: vars.space.xs,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
})

export const statCard = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
})

export const statLabel = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
})

export const statValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
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

export const sectionHeading = style({
  fontSize: vars.fontSize['2xs'],
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
  color: vars.color.textTertiary,
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  margin: 0,
})

export const summaryLede = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textSecondary,
  lineHeight: 1.5,
  margin: 0,
})

export const summaryStrong = style({
  color: vars.color.textPrimary,
  fontWeight: 600,
})

export const stepTitleRow = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.sm,
  minHeight: '32px',
})

export const stepTitleGroup = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
  minWidth: 0,
})

export const stepTitle = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const stepAgent = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const stepDescription = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: 1.5,
  margin: 0,
})

export const riskBanner = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
  borderRadius: vars.radii.sm,
  padding: vars.space.sm,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
})

export const riskBannerBreach = style({
  borderColor: vars.color.dangerBorder,
  backgroundColor: vars.color.dangerMuted,
})

export const riskBannerUrgent = style({
  borderColor: vars.color.warningBorder,
  backgroundColor: vars.color.warningMuted,
})

export const riskHeaderRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
})

export const riskBandLabel = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  fontWeight: 700,
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
})

export const riskMetrics = style({
  display: 'grid',
  gridTemplateColumns: 'repeat(3, 1fr)',
  gap: vars.space.xs,
})

export const riskMetric = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
})

export const riskMetricValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.base,
  fontWeight: 700,
  color: vars.color.textPrimary,
  fontVariantNumeric: 'tabular-nums',
})

export const riskMetricValueDanger = style({
  color: vars.color.danger,
})

export const riskMetricLabel = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.03em',
})

export const dataTable = style({
  display: 'flex',
  flexDirection: 'column',
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  overflow: 'hidden',
})

export const dataTableRow = style({
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1.4fr) minmax(0, 1fr) auto',
  gap: vars.space.sm,
  padding: `${vars.space.xs} ${vars.space.sm}`,
  alignItems: 'center',
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  selectors: {
    '&:first-child': {
      borderTop: 'none',
      backgroundColor: vars.color.surfaceElevated,
    },
  },
})

export const dataTableHeadCell = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  textTransform: 'uppercase',
  letterSpacing: '0.04em',
  fontWeight: 600,
})

export const dataCellPrimary = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textPrimary,
  fontWeight: 500,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const dataCellSecondary = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
})

export const dataCellRight = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
  textAlign: 'right',
  whiteSpace: 'nowrap',
})

export const actionList = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.xs,
})

export const actionRow = style({
  display: 'flex',
  alignItems: 'flex-start',
  gap: vars.space.xs,
  padding: vars.space.sm,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.successBorder}`,
  backgroundColor: vars.color.successMuted,
})

export const actionIcon = style({
  color: vars.color.success,
  flexShrink: 0,
  marginTop: '1px',
})

export const actionText = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '1px',
  minWidth: 0,
})

export const actionType = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const actionReason = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: 1.45,
})

export const chipRow = style({
  display: 'flex',
  flexWrap: 'wrap',
  gap: vars.space.xs,
})

export const chip = style({
  display: 'inline-flex',
  alignItems: 'baseline',
  gap: vars.space['2xs'],
  padding: '3px 8px',
  borderRadius: vars.radii.full,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surfaceElevated,
})

export const chipValue = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 700,
  color: vars.color.textPrimary,
})

export const chipLabel = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
})

export const sqlBlock = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.teal,
  backgroundColor: vars.color.surfaceElevated,
  padding: vars.space.xs,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.border}`,
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-word',
  lineHeight: 1.4,
})

export const judgeLoop = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
})

export const judgeCard = style({
  display: 'flex',
  flexDirection: 'column',
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  padding: vars.space.md,
  gap: vars.space.xs,
})

export const judgeCardApproved = style({
  borderColor: vars.color.successBorder,
})

export const judgeCardRejected = style({
  borderColor: vars.color.dangerBorder,
})

export const judgeHeader = style({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  gap: vars.space.sm,
})

export const judgeHeaderLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  minWidth: 0,
})

export const judgeName = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
})

export const attemptPill = style({
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  fontFamily: fonts.mono,
  padding: '1px 6px',
  borderRadius: vars.radii.full,
  border: `1px solid ${vars.color.border}`,
  backgroundColor: vars.color.surface,
  flexShrink: 0,
})

export const judgeRationale = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  lineHeight: 1.5,
  margin: 0,
  borderLeft: `2px solid ${vars.color.border}`,
  paddingLeft: vars.space.sm,
})

export const verdictBanner = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  padding: vars.space.sm,
  borderRadius: vars.radii.sm,
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
})

export const verdictBannerApproved = style({
  border: `1px solid ${vars.color.successBorder}`,
  backgroundColor: vars.color.successMuted,
  color: vars.color.success,
})

export const verdictBannerRejected = style({
  border: `1px solid ${vars.color.dangerBorder}`,
  backgroundColor: vars.color.dangerMuted,
  color: vars.color.danger,
})

export const rawToggle = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  background: 'none',
  border: 'none',
  padding: 0,
  cursor: 'pointer',
  color: vars.color.textTertiary,
  fontSize: vars.fontSize['2xs'],
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
  ':hover': {
    color: vars.color.textSecondary,
  },
})

export const rawJsonPre = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textSecondary,
  backgroundColor: vars.color.background,
  padding: vars.space.sm,
  borderRadius: vars.radii.xs,
  border: `1px solid ${vars.color.borderSubtle}`,
  overflowX: 'auto',
  maxHeight: '280px',
  lineHeight: 1.45,
  marginTop: vars.space.xs,
})

export const pulseDot = style({
  width: '6px',
  height: '6px',
  borderRadius: '50%',
  backgroundColor: vars.color.warning,
  display: 'inline-block',
  marginRight: '6px',
  animation: `${pulseGlow} 1.5s infinite ease-in-out`,
})

export const emptyState = style({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: `${vars.space['3xl']} ${vars.space['2xl']}`,
  textAlign: 'center',
  color: vars.color.textTertiary,
  gap: vars.space.sm,
  minHeight: '380px',
  flex: 1,
})

export const emptyTitle = style({
  fontSize: vars.fontSize.base,
  fontWeight: 500,
  color: vars.color.textSecondary,
})

export const emptyText = style({
  fontSize: vars.fontSize.xs,
  color: vars.color.textTertiary,
  maxWidth: '320px',
})

export const footer = style({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  gap: vars.space.sm,
  padding: `${vars.space.sm} ${vars.space.lg}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
  backgroundColor: vars.color.surfaceElevated,
  flexShrink: 0,
})
