export const getInitials = (name: string): string => {
  if (!name)
    return ''

  const parts = name.trim().split(/\s+/)

  if (parts.length === 1)
    return parts[0][0].toUpperCase()

  if (parts.length === 2)
    return (parts[0][0] + parts[1][0]).toUpperCase()

  // for 3+ words: first + last
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}
