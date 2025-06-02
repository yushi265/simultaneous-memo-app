/**
 * Get the full image URL from a relative path
 */
export function getFullImageUrl(url: string): string {
  if (url.startsWith('http')) {
    return url
  }
  
  const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  return `${baseUrl}${url}`
}

/**
 * Check if an image URL is accessible
 */
export function checkImageAccessibility(url: string): Promise<boolean> {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => resolve(true)
    img.onerror = () => resolve(false)
    img.src = url
  })
}

/**
 * Preload an image
 */
export function preloadImage(url: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve(img)
    img.onerror = reject
    img.src = url
  })
}