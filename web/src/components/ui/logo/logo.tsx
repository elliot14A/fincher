import type { JSX } from 'preact'
import { logoRecipe } from './logo.css'

export type LogoProps = JSX.HTMLAttributes<HTMLSpanElement> & {
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

export function Logo({ size = 'md', class: userClass, ...props }: LogoProps) {
  const recipeClass = logoRecipe({ size })
  const combinedClass = userClass ? `${recipeClass} ${userClass}` : recipeClass

  return (
    <span class={combinedClass} {...props}>
      Fincher
    </span>
  )
}
