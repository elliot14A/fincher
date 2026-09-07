import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export {
  closeBtn,
  emptyNotice,
  footer,
  header,
  headerMain,
  metadataCard,
  metadataKey,
  metadataRow,
  metadataVal,
  scrollArea,
  section,
  sectionTitle,
  sidebarPanel,
  statCard,
  statGrid,
  statLabel,
  statValue,
} from '#/components/ui/sidebar'

export const componentAvatar = style({
  width: '42px',
  height: '42px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.primary,
  flexShrink: 0,
})

export const packageInfo = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '3px',
  minWidth: 0,
})

export const packageIdClass = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
  margin: 0,
  lineHeight: vars.lineHeight.tight,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
})

export const packageMeta = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textSecondary,
})

export const headerActions = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space['2xs'],
  flexShrink: 0,
})

export const backBtn = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space['3xs'],
  background: 'transparent',
  border: `1px solid ${vars.color.border}`,
  padding: `${vars.space['3xs']} ${vars.space.xs}`,
  color: vars.color.textSecondary,
  cursor: 'pointer',
  borderRadius: vars.radii.xs,
  fontSize: vars.fontSize['2xs'],
  fontWeight: 500,
  transition: 'all 0.12s ease',
  ':hover': {
    color: vars.color.textPrimary,
    backgroundColor: vars.color.surfaceHover,
    borderColor: vars.color.borderStrong,
  },
})

export const failedNotice = style({
  display: 'flex',
  alignItems: 'flex-start',
  gap: vars.space.sm,
  padding: vars.space.sm,
  backgroundColor: vars.color.dangerMuted,
  border: `1px solid ${vars.color.dangerBorder}`,
  borderRadius: vars.radii.xs,
  color: vars.color.danger,
  fontSize: vars.fontSize.xs,
  lineHeight: vars.lineHeight.normal,
})

export const failedNoticeIcon = style({
  flexShrink: 0,
  marginTop: '2px',
})

export const entityCard = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  width: '100%',
  textAlign: 'left',
  font: 'inherit',
  padding: vars.space.sm,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  cursor: 'pointer',
  transition: 'all 0.15s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    borderColor: vars.color.borderStrong,
  },
})

export const entityCardLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  minWidth: 0,
})

export const entityCardInfo = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '2px',
  minWidth: 0,
})

export const entityCardName = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const entityCardMeta = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

export const entityCardChevron = style({
  color: vars.color.textTertiary,
  transition: 'color 0.15s ease, transform 0.15s ease',
  selectors: {
    [`${entityCard}:hover &`]: {
      color: vars.color.primary,
      transform: 'translateX(2px)',
    },
  },
})
