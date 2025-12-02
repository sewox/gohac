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
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
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
  
  logout: () => api.post('/auth/logout'),
}

export const pagesAPI = {
  list: () => api.get('/v1/pages'),
  getById: (id: string) => api.get(`/v1/pages/${id}`),
  create: (data: any) => api.post('/v1/pages', data),
  update: (id: string, data: any) => api.put(`/v1/pages/${id}`, data),
  delete: (id: string) => api.delete(`/v1/pages/${id}`),
}

