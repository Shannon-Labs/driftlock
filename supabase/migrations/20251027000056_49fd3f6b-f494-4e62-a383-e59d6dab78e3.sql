-- Enhanced user setup: Create organization, membership, subscription, and API key on signup

-- Drop existing trigger and function to recreate
DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;
DROP FUNCTION IF EXISTS public.handle_new_user();

-- Create enhanced user setup function
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
  new_org_id UUID;
  org_slug TEXT;
  api_key_secret TEXT;
  api_key_prefix TEXT;
BEGIN
  -- Create profile
  INSERT INTO public.profiles (id, email, name)
  VALUES (
    NEW.id, 
    NEW.email, 
    COALESCE(NEW.raw_user_meta_data->>'name', split_part(NEW.email, '@', 1))
  );

  -- Generate unique org slug from email
  org_slug := lower(regexp_replace(split_part(NEW.email, '@', 1), '[^a-zA-Z0-9]', '', 'g')) || '_' || substr(NEW.id::text, 1, 8);

  -- Create default organization
  INSERT INTO public.organizations (id, name, slug, settings)
  VALUES (
    gen_random_uuid(),
    COALESCE(NEW.raw_user_meta_data->>'name', split_part(NEW.email, '@', 1)) || '''s Organization',
    org_slug,
    '{}'::jsonb
  )
  RETURNING id INTO new_org_id;

  -- Add user as organization owner
  INSERT INTO public.organization_members (organization_id, user_id, role)
  VALUES (new_org_id, NEW.id, 'owner');

  -- Create default subscription (free tier)
  INSERT INTO public.subscriptions (
    user_id, 
    organization_id, 
    tier, 
    plan,
    status,
    current_period_start,
    current_period_end,
    included_calls,
    overage_rate_per_call
  )
  VALUES (
    NEW.id,
    new_org_id,
    'free',
    'free',
    'active',
    NOW(),
    NOW() + INTERVAL '30 days',
    10000,
    0
  );

  -- Create org settings with defaults
  INSERT INTO public.org_settings (
    organization_id,
    anomaly_sensitivity,
    usage_alert_thresholds,
    dunning_behavior
  )
  VALUES (
    new_org_id,
    0.5,
    '{"70": true, "90": true, "100": true}'::jsonb,
    'soft_cap'
  );

  -- Generate API key
  api_key_secret := 'sk_test_' || encode(gen_random_bytes(32), 'hex');
  api_key_prefix := 'sk_test_' || substr(encode(gen_random_bytes(6), 'hex'), 1, 12);

  INSERT INTO public.api_keys (
    user_id,
    organization_id,
    name,
    key_prefix,
    key_hash,
    tier,
    is_active,
    permissions,
    created_by
  )
  VALUES (
    NEW.id,
    new_org_id,
    'Default API Key',
    api_key_prefix,
    encode(digest(api_key_secret, 'sha256'), 'hex'),
    'free',
    true,
    '{"stream": true, "monitor": true}'::jsonb,
    NEW.id
  );

  RETURN NEW;
END;
$$;

-- Recreate trigger
CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW
  EXECUTE FUNCTION public.handle_new_user();