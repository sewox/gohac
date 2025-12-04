import axios from 'axios'

// Create axios instance with default config
const api = axios.create({
  baseURL: '/api',
  withCredentials: true, // Important for cookies
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor: Handle 401 Unauthorized
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Only redirect if not already on login page to avoid redirect loops
      if (window.location.pathname !== '/admin/login') {
        window.location.href = '/admin/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api

// API endpoints
export const authAPI = {
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
  
  me: () => api.get('/auth/me'),
  
  updateProfile: (data: { name?: string; password?: string }) =>
    api.put('/auth/profile', data),
  
  logout: () => api.post('/auth/logout'),
}

export const pagesAPI = {
  list: () => api.get('/v1/pages'),
  getById: (id: string) => api.get(`/v1/pages/${id}`),
  create: (data: any) => api.post('/v1/pages', data),
  update: (id: string, data: any) => api.put(`/v1/pages/${id}`, data),
  delete: (id: string) => api.delete(`/v1/pages/${id}`),
}

export const settingsAPI = {
  get: () => api.get('/public/settings'),
  update: (data: any) => api.put('/v1/settings', data),
}

export const menusAPI = {
  list: () => api.get('/v1/menus'),
  get: (id: string) => api.get(`/v1/menus/${id}`),
  create: (data: any) => api.post('/v1/menus', data),
  update: (id: string, data: any) => api.put(`/v1/menus/${id}`, data),
  delete: (id: string) => api.delete(`/v1/menus/${id}`),
  getPublic: (id: string) => api.get(`/public/menus/${id}`),
}

export const usersAPI = {
  list: () => api.get('/v1/users'),
  get: (id: string) => api.get(`/v1/users/${id}`),
  create: (data: any) => api.post('/v1/users', data),
  update: (id: string, data: any) => api.put(`/v1/users/${id}`, data),
  delete: (id: string) => api.delete(`/v1/users/${id}`),
}

export const mediaAPI = {
  list: () => api.get('/v1/media'),
  get: (filename: string) => api.get(`/v1/media/${filename}`),
}

export const postsAPI = {
  list: (status?: string) => {
    const params = status ? { status } : {}
    return api.get('/v1/posts', { params })
  },
  getById: (id: string) => api.get(`/v1/posts/${id}`),
  create: (data: any) => api.post('/v1/posts', data),
  update: (id: string, data: any) => api.put(`/v1/posts/${id}`, data),
  delete: (id: string) => api.delete(`/v1/posts/${id}`),
}

export const categoriesAPI = {
  list: () => api.get('/v1/categories'),
  getById: (id: string) => api.get(`/v1/categories/${id}`),
  create: (data: any) => api.post('/v1/categories', data),
  update: (id: string, data: any) => api.put(`/v1/categories/${id}`, data),
  delete: (id: string) => api.delete(`/v1/categories/${id}`),
}

export const dashboardAPI = {
  getStats: () => api.get('/v1/dashboard/stats'),
}

