// Block types for the editor
export interface Block {
  id: string
  type: 'hero' | 'text' | 'image' | 'features' | 'pricing' | 'faq' | 'testimonial' | 'video' | 'cta'
  data: HeroData | TextData | ImageData | FeaturesData | PricingData | FAQData | TestimonialData | VideoData | CTAData
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

export interface FeaturesData {
  title?: string
  subtitle?: string
  items: FeatureItem[]
  columns?: 2 | 3 | 4
}

export interface FeatureItem {
  title: string
  description?: string
  icon?: string
}

export interface PricingData {
  title?: string
  subtitle?: string
  plans: PricingPlan[]
}

export interface PricingPlan {
  name: string
  price: string
  description?: string
  features: string[]
  button_text?: string
  button_url?: string
  highlighted?: boolean
}

export interface FAQData {
  title?: string
  items: FAQItem[]
}

export interface FAQItem {
  question: string
  answer: string
}

export interface TestimonialData {
  title?: string
  subtitle?: string
  testimonials: TestimonialItem[]
}

export interface TestimonialItem {
  quote: string
  author: string
  avatar_url?: string
  role?: string
}

export interface VideoData {
  url: string
  title?: string
  description?: string
  autoplay?: boolean
  loop?: boolean
}

export interface CTAData {
  title: string
  subtitle?: string
  button_text: string
  button_url: string
  button_style?: 'primary' | 'secondary' | 'outline'
  background?: string
}

