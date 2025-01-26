import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8001', // Your server URL
});

// Add token to headers for requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
});

export default api;
