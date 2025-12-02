// API helper functions for Astro theme

const API_BASE_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:3131';

export interface GlobalSettings {
  site_name: string;
  logo: string;
  favicon: string;
  contact_email: string;
  header_menu_id?: string;
  footer_menu_id?: string;
}

export interface MenuItem {
  label: string;
  url: string;
  target?: string;
  children?: MenuItem[];
}

export interface Menu {
  id: string;
  name: string;
  description?: string;
  items: MenuItem[];
  created_at: string;
  updated_at: string;
}

export async function getGlobalSettings(): Promise<GlobalSettings> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/public/settings`);
    if (!response.ok) {
      throw new Error(`Failed to fetch settings: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching global settings:', error);
    // Return default settings on error
    return {
      site_name: 'Gohac CMS',
      logo: '',
      favicon: '',
      contact_email: '',
    };
  }
}

export async function getMenuById(id: string): Promise<Menu> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/public/menus/${id}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch menu: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error(`Error fetching menu ${id}:`, error);
    // Return empty menu on error
    return {
      id: '',
      name: '',
      items: [],
      created_at: '',
      updated_at: '',
    };
  }
}

