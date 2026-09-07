import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export * from '#/components/ui/sidebar/sidebar.css'

export const territoryAvatar = style({
  width: '42px',
  height: '42px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.surfaceElevated,
  border: `1px solid ${vars.color.border}`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.sm,
  fontWeight: 700,
  color: vars.color.primary,
  flexShrink: 0,
})

export const holdNotice = style({
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

export const holdNoticeIcon = style({
  flexShrink: 0,
  marginTop: '2px',
})

export const parentTitleCard = style({
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

export const parentTitleLeft = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.sm,
  minWidth: 0,
})

export const parentTitleInfo = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '2px',
  minWidth: 0,
})

export const parentTitleName = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const parentTitleMeta = style({
  fontSize: vars.fontSize['2xs'],
  color: vars.color.textTertiary,
})

export const parentTitleChevron = style({
  color: vars.color.textTertiary,
  transition: 'color 0.15s ease, transform 0.15s ease',
  selectors: {
    [`${parentTitleCard}:hover &`]: {
      color: vars.color.primary,
      transform: 'translateX(2px)',
    },
  },
})
