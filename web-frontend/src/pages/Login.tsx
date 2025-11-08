import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { toast } from 'sonner';
import { Loader2 } from 'lucide-react';

export default function Login() {
  const [apiKey, setApiKey] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { setApiKey: setAuthApiKey } = useAuth();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      // Validate API key format (basic check)
      if (!apiKey || apiKey.length < 10) {
        throw new Error('Please enter a valid API key');
      }

      // Test the API key by making a simple request
      const response = await fetch('http://localhost:8080/v1/version', {
        headers: {
          'Authorization': `Bearer ${apiKey}`,
        },
      });

      if (!response.ok && response.status === 401) {
        throw new Error('Invalid API key');
      }

      // Set the API key in auth context
      setAuthApiKey(apiKey);
      toast.success('Successfully logged in!');
      navigate('/dashboard');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to login';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-hero flex items-center justify-center p-4">
      <Card className="w-full max-w-md glass-card">
        <CardHeader className="space-y-1">
          <div className="flex justify-center mb-4">
            <div className="w-12 h-12 rounded-lg bg-gradient-primary"></div>
          </div>
          <CardTitle className="text-2xl text-center">Welcome back</CardTitle>
          <CardDescription className="text-center">
            Sign in with your API key
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleLogin} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="apiKey">API Key</Label>
              <Input
                id="apiKey"
                type="password"
                placeholder="Enter your API key"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                required
                disabled={loading}
              />
              <p className="text-xs text-muted-foreground">
                Get your API key from the DriftLock API server configuration
              </p>
            </div>
            <Button
              type="submit"
              className="w-full bg-gradient-primary hover:opacity-90"
              disabled={loading}
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Signing in...
                </>
              ) : (
                'Sign in'
              )}
            </Button>
          </form>

          <div className="mt-6 text-center space-y-2">
            <div className="text-sm text-muted-foreground">
              For OSS deployments, set{' '}
              <code className="text-xs bg-muted px-1 py-0.5 rounded">DEFAULT_API_KEY</code>
              {' '}in your API server (it will mirror into{' '}
              <code className="text-xs bg-muted px-1 py-0.5 rounded">DRIFTLOCK_DEV_API_KEY</code>{' '}
              when unset) and reuse that key to log into the dashboard.
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
