import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

const Signup = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({ name: "", email: "", organization: "" });
  const [submitting, setSubmitting] = useState(false);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setSubmitting(true);

    try {
      // In OSS mode we simply guide the user to generate an API key locally.
      toast.success("You're all set! Generate a DEFAULT_API_KEY in your .env to finish sign-up.");
      navigate("/login");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-hero flex items-center justify-center p-4">
      <Card className="w-full max-w-lg glass-card">
        <CardHeader className="space-y-2 text-center">
          <CardTitle className="text-3xl">Create your Driftlock workspace</CardTitle>
          <CardDescription>
            Driftlock OSS is self-hosted—share a few details so we can point you to the right docs and key generation flow.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" name="name" value={form.name} onChange={handleChange} placeholder="Ada Lovelace" required disabled={submitting} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" name="email" value={form.email} onChange={handleChange} placeholder="ada@example.com" required disabled={submitting} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="organization">Organization</Label>
              <Input
                id="organization"
                name="organization"
                value={form.organization}
                onChange={handleChange}
                placeholder="Shannon Labs"
                required
                disabled={submitting}
              />
            </div>
            <Button type="submit" className="w-full bg-gradient-primary" disabled={submitting}>
              {submitting ? "Preparing instructions..." : "Continue"}
            </Button>
            <p className="text-xs text-muted-foreground text-center">
              Already have everything configured? <button type="button" className="underline" onClick={() => navigate("/login")}>Log in</button>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default Signup;
