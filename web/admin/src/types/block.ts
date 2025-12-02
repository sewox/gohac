// Block types for the editor
export interface Block {
  id: string
  type: 'hero' | 'text' | 'image'
  data: HeroData | TextData | ImageData
}

export interface HeroData {
  title: string
  subtitle?: string
  image_url?: string
}

export interface TextData {
  content: string
  align?: 'left' | 'center' | 'right'
}

export interface ImageData {
  url: string
  alt?: string
  caption?: string
}

