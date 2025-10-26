import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';

import { AuthProvider } from './src/contexts/AuthContext';
import { AnomalyProvider } from './src/contexts/AnomalyContext';

import Layout from './src/components/Layout';
import Dashboard from './src/pages/Dashboard';
import Anomalies from './src/pages/Anomalies';
import Events from './src/pages/Events';
import Billing from './src/pages/Billing';
import Profile from './src/pages/Profile';
import Login from './src/pages/Login';
import Register from './src/pages/Register';
import AnomalyDetail from './src/pages/AnomalyDetail';
import PrivateRoute from './src/components/PrivateRoute';

const theme = createTheme({
  palette: {
    primary: {
      main: '#4f46e5',
    },
    secondary: {
      main: '#f59e0b',
    },
    background: {
      default: '#f9fafb',
    },
  },
  typography: {
    fontFamily: [
      'Inter',
      'Arial',
      'sans-serif'
    ].join(','),
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <AnomalyProvider>
          <Router>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/" element={
                <PrivateRoute>
                  <Layout />
                </PrivateRoute>
              }>
                <Route index element={<Dashboard />} />
                <Route path="dashboard" element={<Dashboard />} />
                <Route path="anomalies" element={<Anomalies />} />
                <Route path="anomalies/:id" element={<AnomalyDetail />} />
                <Route path="events" element={<Events />} />
                <Route path="billing" element={<Billing />} />
                <Route path="profile" element={<Profile />} />
              </Route>
            </Routes>
          </Router>
        </AnomalyProvider>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;