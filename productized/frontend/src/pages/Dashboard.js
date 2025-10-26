import React, { useEffect } from 'react';
import {
  Container,
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  CardHeader,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import { dashboardAPI } from '../api';
import { useAnomaly } from '../contexts/AnomalyContext';

const Dashboard = () => {
  const { user } = useAuth();
  const { fetchAnomalies } = useAnomaly();
  const [dashboardStats, setDashboardStats] = React.useState({});
  const [recentAnomalies, setRecentAnomalies] = React.useState([]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const statsResponse = await dashboardAPI.getStats();
        const recentResponse = await dashboardAPI.getRecentAnomalies();
        
        setDashboardStats(statsResponse.data || {});
        setRecentAnomalies(recentResponse.data || []);
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error);
      }
    };

    fetchDashboardData();
    fetchAnomalies();
  }, [fetchAnomalies]);

  return (
    <Container maxWidth="xl">
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      <Typography variant="h6" gutterBottom>
        Welcome back, {user?.name || user?.email}!
      </Typography>

      <Grid container spacing={3} sx={{ mt: 1 }}>
        {/* Stats Cards */}
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Anomalies
              </Typography>
              <Typography variant="h5">
                {dashboardStats.totalAnomalies || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Active Anomalies
              </Typography>
              <Typography variant="h5">
                {dashboardStats.activeAnomalies || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Events Processed
              </Typography>
              <Typography variant="h5">
                {dashboardStats.eventsProcessed || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Anomaly Severity
              </Typography>
              <Typography variant="h5">
                {dashboardStats.avgSeverity || 'N/A'}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Anomalies */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" gutterBottom>
              Recent Anomalies
            </Typography>
            {recentAnomalies.length > 0 ? (
              <Box>
                {recentAnomalies.map((anomaly) => (
                  <Box key={anomaly.id} sx={{ borderBottom: '1px solid #eee', py: 1 }}>
                    <Typography variant="subtitle2">{anomaly.title}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {anomaly.description.substring(0, 100)}...
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {new Date(anomaly.timestamp).toLocaleString()}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography variant="body2" color="text.secondary">
                No recent anomalies detected.
              </Typography>
            )}
          </Paper>
        </Grid>

        {/* System Health */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" gutterBottom>
              System Health
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Typography variant="body2">API Response Time: <strong>45ms</strong></Typography>
              <Typography variant="body2">System Load: <strong>Medium</strong></Typography>
              <Typography variant="body2">Data Processed: <strong>1.2TB</strong></Typography>
              <Typography variant="body2">Active Connections: <strong>1,240</strong></Typography>
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Container>
  );
};

export default Dashboard;