'use client'

import { NodeViewWrapper, NodeViewProps } from '@tiptap/react'
import { useState, useRef, useCallback } from 'react'

export function ResizableImage({ node, updateAttributes, selected }: NodeViewProps) {
  const { src, alt, title, width, height } = node.attrs as {
    src?: string
    alt?: string
    title?: string
    width?: number
    height?: number
    'data-image-id'?: string
    'data-width'?: string
    'data-height'?: string
  }
  const [isResizing, setIsResizing] = useState(false)
  const [resizeStart, setResizeStart] = useState({ x: 0, y: 0, width: 0, height: 0 })
  const imageRef = useRef<HTMLImageElement>(null)

  // Get original dimensions from data attributes or use current dimensions
  const attrs = node.attrs as any
  const originalWidth = attrs['data-width'] ? parseInt(attrs['data-width']) : width || 0
  const originalHeight = attrs['data-height'] ? parseInt(attrs['data-height']) : height || 0
  const aspectRatio = originalWidth && originalHeight ? originalWidth / originalHeight : 1

  const handleMouseDown = useCallback((event: React.MouseEvent) => {
    event.preventDefault()
    setIsResizing(true)
    
    const currentWidth = imageRef.current?.offsetWidth || width || originalWidth || 300
    const currentHeight = imageRef.current?.offsetHeight || height || originalHeight || 200
    
    const startData = {
      x: event.clientX,
      y: event.clientY,
      width: currentWidth,
      height: currentHeight,
    }
    setResizeStart(startData)

    const handleMouseMove = (moveEvent: MouseEvent) => {
      const deltaX = moveEvent.clientX - startData.x
      const newWidth = Math.max(100, startData.width + deltaX) // Minimum width of 100px
      const newHeight = newWidth / aspectRatio

      updateAttributes({
        width: Math.round(newWidth),
        height: Math.round(newHeight),
      })
    }

    const handleMouseUp = () => {
      setIsResizing(false)
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }, [aspectRatio, height, originalHeight, originalWidth, updateAttributes, width])

  const handleImageLoad = () => {
    // If no dimensions are set, use the natural dimensions
    if (!width && !height && imageRef.current) {
      const natural = imageRef.current
      updateAttributes({
        width: natural.naturalWidth,
        height: natural.naturalHeight,
        'data-width': natural.naturalWidth.toString(),
        'data-height': natural.naturalHeight.toString(),
      })
    }
  }

  const displayWidth = width || originalWidth || 'auto'
  const displayHeight = height || originalHeight || 'auto'

  if (!src) {
    return (
      <NodeViewWrapper className="inline-block p-4 border border-red-300 bg-red-50 rounded">
        <div className="text-red-600">画像のURLが設定されていません</div>
      </NodeViewWrapper>
    )
  }

  return (
    <NodeViewWrapper
      className={`relative inline-block ${selected ? 'ring-2 ring-blue-500 ring-opacity-50' : ''}`}
      style={{
        width: displayWidth,
        height: displayHeight,
      }}
    >
      <img
        ref={imageRef}
        src={src}
        alt={alt || ''}
        title={title || ''}
        className="max-w-full h-auto rounded-lg block"
        style={{
          width: displayWidth,
          height: displayHeight,
          cursor: selected ? 'pointer' : 'default',
        }}
        onLoad={handleImageLoad}
        draggable={false}
      />
      
      {/* Resize handle - only show when selected */}
      {selected && (
        <div
          className="absolute bottom-0 right-0 w-4 h-4 bg-blue-500 cursor-se-resize opacity-75 hover:opacity-100 rounded-tl"
          style={{
            transform: 'translate(50%, 50%)',
          }}
          onMouseDown={handleMouseDown}
        >
          <div className="w-full h-full flex items-center justify-center">
            <div className="w-2 h-2 border-r border-b border-white"></div>
          </div>
        </div>
      )}

      {/* Loading overlay during resize */}
      {isResizing && (
        <div className="absolute inset-0 bg-blue-100 bg-opacity-30 flex items-center justify-center rounded-lg">
          <div className="bg-white px-2 py-1 rounded shadow text-sm">
            {width} × {height}
          </div>
        </div>
      )}
    </NodeViewWrapper>
  )
}