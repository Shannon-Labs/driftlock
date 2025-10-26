import axios from 'axios';

// Create an axios instance
const API = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
});

// Add a request interceptor to include the auth token
API.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor to handle token expiration
API.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token might be expired, clear it and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Authentication APIs
export const authAPI = {
  login: (email, password) => API.post('/auth/login', { email, password }),
  register: (name, email, password) => API.post('/auth/register', { name, email, password }),
  refreshToken: (refreshToken) => API.post('/auth/refresh', { refresh_token: refreshToken }),
  getProfile: () => API.get('/user'),
  updateProfile: (userData) => API.put('/user', userData),
};

// Anomaly APIs
export const anomalyAPI = {
  getAnomalies: (params) => API.get('/anomalies', { params }),
  getAnomaly: (id) => API.get(`/anomalies/${id}`),
  resolveAnomaly: (id, resolution) => API.put(`/anomalies/${id}/resolve`, { resolution }),
  deleteAnomaly: (id) => API.delete(`/anomalies/${id}`),
};

// Event APIs
export const eventAPI = {
  getEvents: (params) => API.get('/events', { params }),
  ingestEvent: (eventData) => API.post('/events/ingest', eventData),
};

// Dashboard APIs
export const dashboardAPI = {
  getStats: () => API.get('/dashboard/stats'),
  getRecentAnomalies: () => API.get('/dashboard/recent'),
};

// Billing APIs
export const billingAPI = {
  getPlans: () => API.get('/billing/plans'),
  createCheckoutSession: (planId) => API.post('/billing/checkout', { plan_id: planId }),
  getCustomerPortal: () => API.get('/billing/portal'),
  getSubscription: () => API.get('/billing/subscription'),
  cancelSubscription: () => API.delete('/billing/subscription'),
  getUsage: (params) => API.get('/billing/usage', { params }),
  recordUsage: (featureName, quantity) => API.post('/billing/usage', { feature_name: featureName, quantity }),
};

// Email APIs
export const emailAPI = {
  sendTestEmail: (to, subject, body) => API.post('/email/test', { to, subject, body }),
  sendWelcomeEmail: () => API.post('/email/welcome'),
  sendAnomalyAlert: (anomalyTitle, anomalyDescription, to) => 
    API.post('/email/anomaly-alert', { anomaly_title: anomalyTitle, anomaly_description: anomalyDescription, to }),
};

// Onboarding APIs
export const onboardingAPI = {
  getProgress: () => API.get('/onboarding/progress'),
  markStepComplete: (stepName) => API.post('/onboarding/step/complete', { step_name: stepName }),
  skipOnboarding: () => API.post('/onboarding/skip'),
  getResources: () => API.get('/onboarding/resources'),
};

export default API;