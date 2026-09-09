import { style } from '@vanilla-extract/css'
import { vars } from '#/styles/theme.css'
import { fonts } from '#/styles/tokens'

export const page = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  height: '100%',
  width: '100%',
  overflow: 'hidden',
})

export const centerArea = style({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  flex: 1,
  padding: `${vars.space['3xl']} ${vars.space['2xl']}`,
  overflowY: 'auto',
})

export const centerColumn = style({
  display: 'flex',
  flexDirection: 'column',
  maxWidth: '640px',
  width: '100%',
  gap: vars.space['2xl'],
  margin: '0 auto',
})

export const heading = style({
  display: 'flex',
  flexDirection: 'column',
  gap: vars.space.sm,
})

export const title = style({
  fontSize: vars.fontSize['2xl'],
  fontWeight: 600,
  color: vars.color.textPrimary,
  letterSpacing: '-0.015em',
  lineHeight: vars.lineHeight.tight,
  margin: 0,
})

export const subtitle = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textTertiary,
  lineHeight: vars.lineHeight.normal,
  margin: 0,
  maxWidth: '480px',
})

export const composerShell = style({
  position: 'relative',
  width: '100%',
})

export const composer = style({
  backgroundColor: vars.color.surface,
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.sm,
  padding: `${vars.space.md} ${vars.space.lg}`,
  display: 'flex',
  alignItems: 'flex-end',
  gap: vars.space.sm,
  width: '100%',
})

export const composerTextarea = style({
  backgroundColor: 'transparent',
  border: 'none',
  color: vars.color.textPrimary,
  fontSize: vars.fontSize.sm,
  fontFamily: fonts.sans,
  lineHeight: vars.lineHeight.normal,
  outline: 'none',
  flex: 1,
  resize: 'none',
  maxHeight: '160px',
  minHeight: '20px',
  '::placeholder': {
    color: vars.color.textTertiary,
  },
})

export const composerPrompt = style({
  fontFamily: fonts.mono,
  color: vars.color.textTertiary,
  fontSize: vars.fontSize.base,
  flexShrink: 0,
})

export const sendButton = style({
  width: '28px',
  height: '28px',
  borderRadius: vars.radii.xs,
  backgroundColor: vars.color.primary,
  border: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: vars.color.textInverse,
  cursor: 'pointer',
  flexShrink: 0,
  transition: 'background-color 0.12s ease',
  ':hover': {
    backgroundColor: vars.color.primaryHover,
  },
})

export const queryListLabel = style({
  fontSize: vars.fontSize.xs,
  fontWeight: 500,
  color: vars.color.textTertiary,
  marginBottom: vars.space.xs,
})

export const queryList = style({
  display: 'flex',
  flexDirection: 'column',
  border: `1px solid ${vars.color.borderSubtle}`,
  borderRadius: vars.radii.sm,
  overflow: 'hidden',
})

export const queryRow = style({
  display: 'flex',
  alignItems: 'center',
  gap: vars.space.md,
  padding: `${vars.space.sm} ${vars.space.md}`,
  backgroundColor: 'transparent',
  border: 'none',
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  cursor: 'pointer',
  textAlign: 'left',
  transition: 'background-color 0.12s ease',
  ':last-child': {
    borderBottom: 'none',
  },
  ':hover': {
    backgroundColor: vars.color.surface,
  },
})

export const queryIcon = style({
  color: vars.color.textTertiary,
  flexShrink: 0,
})

export const queryText = style({
  fontSize: vars.fontSize.sm,
  color: vars.color.textSecondary,
  lineHeight: vars.lineHeight.snug,
  flex: 1,
})

export const queryTag = style({
  fontFamily: fonts.mono,
  fontSize: vars.fontSize['3xs'],
  color: vars.color.textTertiary,
  flexShrink: 0,
})

export const conversationArea = style({
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  overflowY: 'auto',
  padding: `${vars.space.xl} ${vars.space['2xl']}`,
})

export const conversationColumn = style({
  display: 'flex',
  flexDirection: 'column',
  maxWidth: '760px',
  width: '100%',
  margin: '0 auto',
})

export const conversationHeader = style({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: vars.space.md,
  width: '100%',
  padding: `${vars.space.md} ${vars.space['2xl']}`,
  borderBottom: `1px solid ${vars.color.borderSubtle}`,
  flexShrink: 0,
})

export const conversationTitle = style({
  fontSize: vars.fontSize.sm,
  fontWeight: 600,
  color: vars.color.textPrimary,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const newChatButton = style({
  display: 'inline-flex',
  alignItems: 'center',
  gap: vars.space.xs,
  fontFamily: fonts.mono,
  fontSize: vars.fontSize.xs,
  color: vars.color.textSecondary,
  backgroundColor: 'transparent',
  border: `1px solid ${vars.color.border}`,
  borderRadius: vars.radii.xs,
  padding: `${vars.space['2xs']} ${vars.space.sm}`,
  cursor: 'pointer',
  flexShrink: 0,
  transition: 'background-color 0.12s ease, color 0.12s ease',
  ':hover': {
    backgroundColor: vars.color.surfaceHover,
    color: vars.color.textPrimary,
  },
})

export const composerDock = style({
  padding: `${vars.space.md} ${vars.space['2xl']} ${vars.space.xl}`,
  borderTop: `1px solid ${vars.color.borderSubtle}`,
})

export const composerDockColumn = style({
  maxWidth: '760px',
  width: '100%',
  margin: '0 auto',
})

export const sendButtonDisabled = style({
  backgroundColor: vars.color.surfaceActive,
  color: vars.color.textTertiary,
  cursor: 'not-allowed',
  ':hover': {
    backgroundColor: vars.color.surfaceActive,
  },
})
