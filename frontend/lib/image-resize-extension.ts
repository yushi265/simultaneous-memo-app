import { Node, mergeAttributes } from '@tiptap/core'
import { ReactNodeViewRenderer } from '@tiptap/react'
import { ResizableImage } from '@/components/ResizableImage'

export interface ImageOptions {
  HTMLAttributes: Record<string, any>
  inline: boolean
  allowBase64: boolean
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    resizableImage: {
      /**
       * Add an image
       */
      setImage: (options: {
        src: string
        alt?: string
        title?: string
        width?: number
        height?: number
        'data-image-id'?: string
        'data-width'?: string
        'data-height'?: string
      }) => ReturnType
    }
  }
}

export const ResizableImageExtension = Node.create<ImageOptions>({
  name: 'resizableImage',

  group: 'block',

  content: '',

  draggable: true,

  isolating: true,

  addOptions() {
    return {
      HTMLAttributes: {
        class: 'max-w-full h-auto rounded-lg',
      },
      inline: false,
      allowBase64: false,
    }
  },

  addAttributes() {
    return {
      src: {
        default: null,
        parseHTML: element => element.getAttribute('src'),
        renderHTML: attributes => {
          if (!attributes.src) {
            return {}
          }
          return { src: attributes.src }
        },
      },
      alt: {
        default: null,
        parseHTML: element => element.getAttribute('alt'),
        renderHTML: attributes => {
          if (!attributes.alt) {
            return {}
          }
          return { alt: attributes.alt }
        },
      },
      title: {
        default: null,
        parseHTML: element => element.getAttribute('title'),
        renderHTML: attributes => {
          if (!attributes.title) {
            return {}
          }
          return { title: attributes.title }
        },
      },
      width: {
        default: null,
        parseHTML: element => {
          const width = element.getAttribute('width')
          return width ? parseInt(width, 10) : null
        },
        renderHTML: attributes => {
          if (!attributes.width) {
            return {}
          }
          return { width: attributes.width }
        },
      },
      height: {
        default: null,
        parseHTML: element => {
          const height = element.getAttribute('height')
          return height ? parseInt(height, 10) : null
        },
        renderHTML: attributes => {
          if (!attributes.height) {
            return {}
          }
          return { height: attributes.height }
        },
      },
      'data-image-id': {
        default: null,
        parseHTML: element => element.getAttribute('data-image-id'),
        renderHTML: attributes => {
          if (!attributes['data-image-id']) {
            return {}
          }
          return { 'data-image-id': attributes['data-image-id'] }
        },
      },
      'data-width': {
        default: null,
        parseHTML: element => element.getAttribute('data-width'),
        renderHTML: attributes => {
          if (!attributes['data-width']) {
            return {}
          }
          return { 'data-width': attributes['data-width'] }
        },
      },
      'data-height': {
        default: null,
        parseHTML: element => element.getAttribute('data-height'),
        renderHTML: attributes => {
          if (!attributes['data-height']) {
            return {}
          }
          return { 'data-height': attributes['data-height'] }
        },
      },
    }
  },

  parseHTML() {
    return [
      {
        tag: 'img[src]',
      },
    ]
  },

  renderHTML({ HTMLAttributes }) {
    return ['img', mergeAttributes(this.options.HTMLAttributes, HTMLAttributes)]
  },

  addNodeView() {
    return ReactNodeViewRenderer(ResizableImage)
  },

  addCommands() {
    return {
      setImage:
        options =>
        ({ commands }) => {
          return commands.insertContent({
            type: this.name,
            attrs: options,
          })
        },
    }
  },
})