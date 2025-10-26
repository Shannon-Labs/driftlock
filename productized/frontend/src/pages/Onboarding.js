import React, { useEffect, useState } from 'react';
import {
  Container,
  Box,
  Typography,
  Step,
  StepLabel,
  Stepper,
  Button,
  Card,
  CardContent,
  CardActions,
  LinearProgress,
  Grid,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { onboardingAPI } from '../api';

const OnboardingFlow = () => {
  const [onboardingData, setOnboardingData] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchOnboardingProgress();
  }, []);

  const fetchOnboardingProgress = async () => {
    try {
      const response = await onboardingAPI.getProgress();
      setOnboardingData(response.data);
      setLoading(false);
    } catch (error) {
      console.error('Failed to fetch onboarding progress:', error);
      setLoading(false);
    }
  };

  const markStepComplete = async (stepName) => {
    try {
      await onboardingAPI.markStepComplete(stepName);
      // Refresh the onboarding data
      fetchOnboardingProgress();
    } catch (error) {
      console.error('Failed to mark step complete:', error);
    }
  };

  const skipOnboarding = async () => {
    try {
      await onboardingAPI.skipOnboarding();
      navigate('/dashboard');
    } catch (error) {
      console.error('Failed to skip onboarding:', error);
    }
  };

  if (loading) {
    return (
      <Container maxWidth="md">
        <Typography variant="h4" align="center" sx={{ mt: 4 }}>
          Loading onboarding...
        </Typography>
      </Container>
    );
  }

  if (!onboardingData) {
    return (
      <Container maxWidth="md">
        <Typography variant="h4" align="center" sx={{ mt: 4 }}>
          No onboarding data available
        </Typography>
      </Container>
    );
  }

  const { steps, progress } = onboardingData;

  return (
    <Container maxWidth="md">
      <Typography variant="h4" gutterBottom align="center">
        Welcome to DriftLock!
      </Typography>
      <Typography variant="h6" color="textSecondary" gutterBottom align="center">
        Complete these steps to get started
      </Typography>

      <Box sx={{ width: '100%', mb: 4 }}>
        <LinearProgress variant="determinate" value={progress} />
        <Typography variant="body2" align="center" sx={{ mt: 1 }}>
          {Math.round(progress)}% Complete
        </Typography>
      </Box>

      <Stepper activeStep={steps.findIndex(step => !step.Completed)} orientation="vertical">
        {steps.map((step, index) => (
          <Step key={step.name} completed={step.Completed}>
            <StepLabel>{step.title}</StepLabel>
            <Box sx={{ ml: 4, mt: 1 }}>
              <Card>
                <CardContent>
                  <Typography variant="body2" color="text.secondary">
                    {step.description}
                  </Typography>
                </CardContent>
                <CardActions>
                  {!step.Completed ? (
                    <Button 
                      size="small" 
                      variant="contained"
                      onClick={() => markStepComplete(step.name)}
                    >
                      Mark Complete
                    </Button>
                  ) : (
                    <Typography variant="body2" color="success.main">
                      Completed
                    </Typography>
                  )}
                </CardActions>
              </Card>
            </Box>
          </Step>
        ))}
      </Stepper>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 3 }}>
        <Button onClick={skipOnboarding} color="secondary">
          Skip Onboarding
        </Button>
        <Button 
          variant="contained" 
          onClick={() => navigate('/dashboard')}
          disabled={progress < 100}
        >
          Go to Dashboard
        </Button>
      </Box>
    </Container>
  );
};

export default OnboardingFlow;