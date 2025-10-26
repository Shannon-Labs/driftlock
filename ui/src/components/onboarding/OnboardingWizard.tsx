import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, Stepper, Step, StepLabel, Button, TextField, Select, MenuItem, FormControlLabel, Box, Typography, Alert, AlertTitle, CircularProgress } from '@mui/material';
import { CheckCircle, CloudUpload, Security, Storage, Analytics, Dashboard, Settings } from '@mui/icons-material';

interface OnboardingStep {
  id: string;
  title: string;
  description: string;
  component: React.ReactNode;
  completed: boolean;
}

interface TenantFormData {
  name: string;
  domain: string;
  plan: 'trial' | 'starter' | 'pro' | 'enterprise';
  industry: string;
  companySize: string;
  contactEmail: string;
  contactName: string;
}

interface IntegrationFormData {
  type: string;
  apiKey: string;
  endpoint: string;
}

const OnboardingWizard: React.FC = () => {
  const [activeStep, setActiveStep] = useState(0);
  const [tenantData, setTenantData] = useState<TenantFormData>({
    name: '',
    domain: '',
    plan: 'trial',
    industry: '',
    companySize: '',
    contactEmail: '',
    contactName: '',
  });
  const [integrationData, setIntegrationData] = useState<IntegrationFormData>({
    type: '',
    apiKey: '',
    endpoint: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const steps: OnboardingStep[] = [
    {
      id: 'tenant',
      title: 'Create Tenant',
      description: 'Set up your organization and configure basic settings',
      component: <TenantSetupStep data={tenantData} onChange={setTenantData} />,
      completed: tenantData.name !== '' && tenantData.domain !== '',
    },
    {
      id: 'integration',
      title: 'Configure Integration',
      description: 'Connect your data sources and configure event collection',
      component: <IntegrationSetupStep data={integrationData} onChange={setIntegrationData} />,
      completed: integrationData.type !== '' && integrationData.apiKey !== '',
    },
    {
      id: 'dashboard',
      title: 'Explore Dashboard',
      description: 'View your anomaly detection results and analytics',
      component: <DashboardPreviewStep />,
      completed: true, // Always accessible once tenant is created
    },
    {
      id: 'complete',
      title: 'Complete Setup',
      description: 'Review your configuration and start using Driftlock',
      component: <CompletionStep tenantData={tenantData} integrationData={integrationData} />,
      completed: false, // Only completed when user clicks finish
    },
  ];

  const handleNext = () => {
    if (activeStep < steps.length - 1) {
      setActiveStep(activeStep + 1);
    }
  };

  const handleBack = () => {
    if (activeStep > 0) {
      setActiveStep(activeStep - 1);
    }
  };

  const handleSubmit = async () => {
    setIsSubmitting(true);
    try {
      // In a real implementation, this would submit data to the backend
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Move to completion step
      setActiveStep(steps.length - 1);
    } catch (error) {
      console.error('Onboarding failed:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const getStepIcon = (stepId: string) => {
    switch (stepId) {
      case 'tenant':
        return <Dashboard />;
      case 'integration':
        return <CloudUpload />;
      case 'dashboard':
        return <Dashboard />;
      case 'complete':
        return <CheckCircle />;
      default:
        return <Settings />;
    }
  };

  return (
    <Box sx={{ width: '100%', maxWidth: 800, mx: 'auto', p: 3 }}>
      <Card>
        <CardHeader
          title={
            <Box display="flex" alignItems="center">
              <Typography variant="h5" component="h2">
                Driftlock Onboarding
              </Typography>
              <Box sx={{ ml: 2 }}>
                <Typography variant="body2" color="text.secondary">
                  Get started with anomaly detection in minutes
                </Typography>
              </Box>
            </Box>
          }
        />
        <CardContent>
          <Stepper activeStep={activeStep} orientation="vertical">
            {steps.map((step, index) => (
              <Step key={step.id} completed={step.completed}>
                <StepLabel
                  icon={getStepIcon(step.id)}
                  optional={index > 1}
                >
                  {step.title}
                </StepLabel>
                <StepContent>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    {step.description}
                  </Typography>
                  {step.component}
                </StepContent>
              </Step>
            ))}
          </Stepper>
          
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 3 }}>
            <Button
              disabled={activeStep === 0}
              onClick={handleBack}
              sx={{ mr: 1 }}
            >
              Back
            </Button>
            <Button
              variant="contained"
              disabled={activeStep === steps.length - 1 || isSubmitting}
              onClick={activeStep === steps.length - 1 ? handleSubmit : handleNext}
              startIcon={activeStep === steps.length - 1 ? <CheckCircle /> : undefined}
            >
              {activeStep === steps.length - 1 ? 'Complete Setup' : 'Next'}
            </Button>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};

// Tenant Setup Step Component
const TenantSetupStep: React.FC<{ data: TenantFormData; onChange: (data: TenantFormData) => void }> = ({ data, onChange }) => {
  const handleChange = (field: keyof TenantFormData) => (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      ...data,
      [field]: event.target.value,
    });
  };

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Organization Information
      </Typography>
      <TextField
        fullWidth
        label="Organization Name"
        value={data.name}
        onChange={(e) => handleChange('name', e)}
        margin="normal"
        required
      />
      <TextField
        fullWidth
        label="Domain"
        value={data.domain}
        onChange={(e) => handleChange('domain', e)}
        margin="normal"
        required
        placeholder="your-company.driftlock.com"
      />
      
      <FormControlLabel component="legend">Plan</FormControlLabel>
      <Select
        fullWidth
        value={data.plan}
        onChange={(e) => handleChange('plan', e.target.value as any)}
        margin="normal"
      >
        <MenuItem value="trial">Trial - 14 days, 100 anomalies/day</MenuItem>
        <MenuItem value="starter">Starter - 30 days, 500 anomalies/day</MenuItem>
        <MenuItem value="pro">Pro - 90 days, 2000 anomalies/day</MenuItem>
        <MenuItem value="enterprise">Enterprise - Unlimited</MenuItem>
      </Select>
      
      <TextField
        fullWidth
        label="Industry"
        value={data.industry}
        onChange={(e) => handleChange('industry', e)}
        margin="normal"
        select
      >
        <MenuItem value="finance">Finance</MenuItem>
        <MenuItem value="healthcare">Healthcare</MenuItem>
        <MenuItem value="retail">Retail</MenuItem>
        <MenuItem value="technology">Technology</MenuItem>
        <MenuItem value="manufacturing">Manufacturing</MenuItem>
        <MenuItem value="other">Other</MenuItem>
      </TextField>
      
      <TextField
        fullWidth
        label="Company Size"
        value={data.companySize}
        onChange={(e) => handleChange('companySize', e)}
        margin="normal"
        select
      >
        <MenuItem value="small">Small (1-50 employees)</MenuItem>
        <MenuItem value="medium">Medium (51-500 employees)</MenuItem>
        <MenuItem value="large">Large (500+ employees)</MenuItem>
      </TextField>
      
      <Typography variant="h6" gutterBottom>
        Contact Information
      </Typography>
      <TextField
        fullWidth
        label="Contact Name"
        value={data.contactName}
        onChange={(e) => handleChange('contactName', e)}
        margin="normal"
        required
      />
      <TextField
        fullWidth
        label="Contact Email"
        value={data.contactEmail}
        onChange={(e) => handleChange('contactEmail', e)}
        margin="normal"
        type="email"
        required
      />
    </Box>
  );
};

// Integration Setup Step Component
const IntegrationSetupStep: React.FC<{ data: IntegrationFormData; onChange: (data: IntegrationFormData) => void }> = ({ data, onChange }) => {
  const handleChange = (field: keyof IntegrationFormData) => (event: React.ChangeEvent<HTMLInputElement | { value: string }>) => {
    onChange({
      ...data,
      [field]: event.target.value,
    });
  };

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Data Source Integration
      </Typography>
      <FormControlLabel component="legend">Integration Type</FormControlLabel>
      <Select
        fullWidth
        value={data.type}
        onChange={(e) => handleChange('type', e.target.value as any)}
        margin="normal"
      >
        <MenuItem value="otel">OpenTelemetry Collector</MenuItem>
        <MenuItem value="prometheus">Prometheus Remote Write</MenuItem>
        <MenuItem value="kafka">Kafka Producer</MenuItem>
        <MenuItem value="webhook">Webhook</MenuItem>
        <MenuItem value="syslog">Syslog</MenuItem>
      </Select>
      
      <TextField
        fullWidth
        label="API Key"
        value={data.apiKey}
        onChange={(e) => handleChange('apiKey', e)}
        margin="normal"
        type="password"
        helperText="Your secure API key for data ingestion"
      />
      
      <TextField
        fullWidth
        label="Endpoint URL"
        value={data.endpoint}
        onChange={(e) => handleChange('endpoint', e)}
        margin="normal"
        placeholder="https://your-collector.example.com:4317"
        required
      />
      
      <Alert severity="info" sx={{ mt: 2 }}>
        <AlertTitle>Integration Guide</AlertTitle>
        Need help setting up your integration? Check out our 
        <a href="#" target="_blank" rel="noopener noreferrer">integration documentation</a> for detailed guides.
      </Alert>
    </Box>
  );
};

// Dashboard Preview Step Component
const DashboardPreviewStep: React.FC = () => {
  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Dashboard Preview
      </Typography>
      <Alert severity="success" sx={{ mb: 2 }}>
        <AlertTitle>Ready to Explore</AlertTitle>
        Your tenant is set up! You can now access the Driftlock dashboard to view anomaly detection results.
      </Alert>
      
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
        <Button
          variant="contained"
          size="large"
          startIcon={<Dashboard />}
          href="/dashboard"
        >
          Go to Dashboard
        </Button>
      </Box>
    </Box>
  );
};

// Completion Step Component
const CompletionStep: React.FC<{ tenantData: TenantFormData; integrationData: IntegrationFormData }> = ({ tenantData, integrationData }) => {
  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Setup Complete!
      </Typography>
      
      <Card sx={{ mb: 2 }}>
        <CardHeader title="Tenant Configuration" />
        <CardContent>
          <Typography variant="body2" gutterBottom>
            <strong>Organization:</strong> {tenantData.name}
          </Typography>
          <Typography variant="body2" gutterBottom>
            <strong>Domain:</strong> {tenantData.domain}
          </Typography>
          <Typography variant="body2" gutterBottom>
            <strong>Plan:</strong> {tenantData.plan}
          </Typography>
          <Typography variant="body2" gutterBottom>
            <strong>Industry:</strong> {tenantData.industry}
          </Typography>
          <Typography variant="body2" gutterBottom>
            <strong>Company Size:</strong> {tenantData.companySize}
          </Typography>
        </CardContent>
      </Card>
      
      <Card sx={{ mb: 2 }}>
        <CardHeader title="Integration Configuration" />
        <CardContent>
          <Typography variant="body2" gutterBottom>
            <strong>Type:</strong> {integrationData.type}
          </Typography>
          <Typography variant="body2" gutterBottom>
            <strong>Endpoint:</strong> {integrationData.endpoint}
          </Typography>
        </CardContent>
      </Card>
      
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
        <Button
          variant="contained"
          size="large"
          color="success"
          startIcon={<CheckCircle />}
          href="/dashboard"
        >
          Start Using Driftlock
        </Button>
      </Box>
    </Box>
  );
};

export default OnboardingWizard;
