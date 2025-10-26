import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, Typography, Box, LinearProgress, Chip, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material';
import { TrendingUp, TrendingDown, CheckCircle, Warning, Schedule } from '@mui/icons-material';

interface HealthMetric {
  id: string;
  name: string;
  value: number;
  target: number;
  unit: string;
  trend: 'up' | 'down' | 'stable';
  status: 'good' | 'warning' | 'critical';
  description: string;
}

interface CustomerHealth {
  score: number;
  status: 'excellent' | 'good' | 'fair' | 'poor';
  metrics: HealthMetric[];
  lastUpdated: Date;
}

const CustomerHealth: React.FC = () => {
  const [healthData, setHealthData] = useState<CustomerHealth | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // In a real implementation, this would fetch data from the API
    const fetchHealthData = async () => {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Mock health data
      const mockHealthData: CustomerHealth = {
        score: 78,
        status: 'good',
        lastUpdated: new Date(),
        metrics: [
          {
            id: 'anomaly-detection-rate',
            name: 'Anomaly Detection Rate',
            value: 85,
            target: 90,
            unit: '%',
            trend: 'up',
            status: 'good',
            description: 'Successfully detecting 85% of anomalies with 90% target'
          },
          {
            id: 'false-positive-rate',
            name: 'False Positive Rate',
            value: 5,
            target: 10,
            unit: '%',
            trend: 'down',
            status: 'good',
            description: 'False positive rate reduced from 8% to 5% this month'
          },
          {
            id: 'data-ingestion',
            name: 'Data Ingestion',
            value: 950,
            target: 1000,
            unit: 'events/sec',
            trend: 'stable',
            status: 'good',
            description: 'Processing 950 events per second with 1000 target'
          },
          {
            id: 'api-response-time',
            name: 'API Response Time',
            value: 120,
            target: 100,
            unit: 'ms',
            trend: 'up',
            status: 'warning',
            description: 'Average response time is 120ms with 100ms target'
          },
          {
            id: 'system-uptime',
            name: 'System Uptime',
            value: 99.8,
            target: 99.9,
            unit: '%',
            trend: 'stable',
            status: 'good',
            description: 'System uptime of 99.8% with 99.9% target'
          },
          {
            id: 'storage-utilization',
            name: 'Storage Utilization',
            value: 65,
            target: 80,
            unit: '%',
            trend: 'down',
            status: 'good',
            description: 'Storage utilization at 65% with 80% target'
          }
        ]
      };
      
      setHealthData(mockHealthData);
      setLoading(false);
    };

    fetchHealthData();
  }, []);

  const getHealthColor = (status: string) => {
    switch (status) {
      case 'excellent':
        return '#4caf50';
      case 'good':
        return '#8bc34a';
      case 'fair':
        return '#ff9800';
      case 'poor':
        return '#f44336';
      default:
        return '#9e9e9e';
    }
  };

  const getHealthIcon = (status: string) => {
    switch (status) {
      case 'excellent':
      case 'good':
        return <CheckCircle sx={{ color: getHealthColor(status) }} />;
      case 'fair':
        return <Warning sx={{ color: getHealthColor(status) }} />;
      case 'poor':
        return <Warning sx={{ color: getHealthColor(status) }} />;
      default:
        return <Schedule sx={{ color: getHealthColor(status) }} />;
    }
  };

  const getTrendIcon = (trend: string) => {
    return trend === 'up' ? 
      <TrendingUp sx={{ color: '#4caf50' }} /> : 
      <TrendingDown sx={{ color: '#f44336' }} />;
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
        <Typography variant="h5">Loading customer health data...</Typography>
      </Box>
    );
  }

  if (!healthData) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
        <Typography variant="h5">No health data available</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Customer Health Score
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent sx={{ textAlign: 'center', py: 4 }}>
              <Typography variant="h3" component="div" sx={{ color: getHealthColor(healthData.status) }}>
                {healthData.score}
              </Typography>
              <Typography variant="h6" color="textSecondary">
                Health Score
              </Typography>
              <Box sx={{ mt: 2, display: 'flex', justifyContent: 'center' }}>
                {getHealthIcon(healthData.status)}
                <Typography variant="body2" sx={{ ml: 1, textTransform: 'capitalize' }}>
                  {healthData.status}
                </Typography>
              </Box>
              <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
                Last updated: {healthData.lastUpdated.toLocaleDateString()}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={8}>
          <Card>
            <CardHeader title="Health Metrics" />
            <CardContent>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Metric</TableCell>
                      <TableCell>Current</TableCell>
                      <TableCell>Target</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Trend</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {healthData.metrics.map((metric) => (
                      <TableRow key={metric.id}>
                        <TableCell component="th" scope="row">
                          {metric.name}
                        </TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center' }}>
                            <Typography variant="body2">
                              {metric.value} {metric.unit}
                            </Typography>
                            <Box sx={{ ml: 1, display: 'flex', alignItems: 'center' }}>
                              {getTrendIcon(metric.trend)}
                              <Typography variant="caption" sx={{ ml: 0.5 }}>
                                {metric.trend}
                              </Typography>
                            </Box>
                          </Box>
                        </TableCell>
                        <TableCell>
                          {metric.target} {metric.unit}
                        </TableCell>
                        <TableCell>
                          <Chip 
                            label={metric.status} 
                            size="small" 
                            color={
                              metric.status === 'good' ? 'success' : 
                              metric.status === 'warning' ? 'warning' : 'error'
                            }
                          />
                        </TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center' }}>
                            {getTrendIcon(metric.trend)}
                          </Box>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      <Box sx={{ mt: 3 }}>
        <Paper sx={{ p: 2, bgcolor: 'background.default' }}>
          <Typography variant="h6" gutterBottom>
            Health Score Details
          </Typography>
          <Typography variant="body2" paragraph>
            Your customer health score is calculated based on several key metrics:
          </Typography>
          <Typography variant="body2" paragraph>
            <strong>Anomaly Detection Rate:</strong> How effectively we're detecting anomalies in your data
          </Typography>
          <Typography variant="body2" paragraph>
            <strong>False Positive Rate:</strong> The percentage of normal events incorrectly flagged as anomalies
          </Typography>
          <Typography variant="body2" paragraph>
            <strong>Data Ingestion:</strong> The rate at which we're processing your event data
          </Typography>
          <Typography variant="body2" paragraph>
            <strong>API Response Time:</strong> The average time it takes for our API to respond to requests
          </Typography>
          <Typography variant="body2" paragraph>
            <strong>System Uptime:</strong> The percentage of time our systems are operational
          </Typography>
          <Typography variant="body2" paragraph>
            <strong>Storage Utilization:</strong> How much of your allocated storage you're using
          </Typography>
        </Paper>
      </Box>
    </Box>
  );
};

export default CustomerHealth;
