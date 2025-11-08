import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

const ResetPassword = () => {
  const [email, setEmail] = useState("");
  const [processing, setProcessing] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setProcessing(true);

    try {
      toast.success("OSS deployments use API keys—set a new DEFAULT_API_KEY in .env to rotate credentials.");
    } finally {
      setProcessing(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-hero flex items-center justify-center p-4">
      <Card className="w-full max-w-lg glass-card">
        <CardHeader className="space-y-2 text-center">
          <CardTitle className="text-3xl">Reset access</CardTitle>
          <CardDescription>
            Lost track of your dashboard key? Generate a new `DEFAULT_API_KEY` (and optionally rotate `DRIFTLOCK_DEV_API_KEY`). We’ll send housekeeping tips to the email below.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="reset-email">Email</Label>
              <Input
                id="reset-email"
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="ops@example.com"
                required
                disabled={processing}
              />
            </div>
            <Button type="submit" className="w-full bg-gradient-primary" disabled={processing}>
              {processing ? "Sending guidance..." : "Send instructions"}
            </Button>
            <p className="text-xs text-muted-foreground text-center">
              Need help rotating keys in production? Review README.md → Authentication or contact your security admin.
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default ResetPassword;
