-- Fix security warnings from linter

-- Fix search_path for existing functions
CREATE OR REPLACE FUNCTION public.deduct_user_credit(p_user_id uuid)
RETURNS void
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  UPDATE user_credits 
  SET balance = balance - 1, total_used = total_used + 1, updated_at = now()
  WHERE user_id = p_user_id AND balance > 0;
END;
$$;

CREATE OR REPLACE FUNCTION public.add_user_credits(p_user_id uuid, p_amount integer)
RETURNS void
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  INSERT INTO user_credits (user_id, balance, total_purchased, updated_at)
  VALUES (p_user_id, p_amount, p_amount, now())
  ON CONFLICT (user_id) 
  DO UPDATE SET 
    balance = user_credits.balance + p_amount,
    total_purchased = user_credits.total_purchased + p_amount,
    updated_at = now();
END;
$$;

CREATE OR REPLACE FUNCTION public.check_user_credits(p_user_id uuid)
RETURNS boolean
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  RETURN EXISTS (
    SELECT 1 FROM user_credits 
    WHERE user_id = p_user_id AND balance > 0
  );
END;
$$;

CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS trigger
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  INSERT INTO public.profiles (id, email, name)
  VALUES (new.id, new.email, new.raw_user_meta_data->>'full_name');

  -- Create default subscription
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'subscriptions' AND column_name = 'quota') THEN
    INSERT INTO public.subscriptions (user_id, tier, quota)
    VALUES (new.id, 'free', 10000);
  ELSIF EXISTS (SELECT 1 FROM information_schema.columns
                WHERE table_name = 'subscriptions' AND column_name = 'monthly_quota') THEN
    INSERT INTO public.subscriptions (user_id, tier, monthly_quota)
    VALUES (new.id, 'free', 10000);
  END IF;

  RETURN new;
END;
$$;

CREATE OR REPLACE FUNCTION public.generate_api_key_for_subscription()
RETURNS trigger
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM api_keys
    WHERE user_id = NEW.user_id
    AND revoked = false
  ) THEN
    INSERT INTO api_keys (user_id, key_secret, key_prefix, tier, name)
    VALUES (
      NEW.user_id,
      'sk_' || CASE WHEN NEW.tier = 'enterprise' THEN 'enterprise' ELSE 'test' END || '_' || encode(gen_random_bytes(32), 'hex'),
      'sk_' || CASE WHEN NEW.tier = 'enterprise' THEN 'enterprise' ELSE 'test' END || '_' || substr(encode(gen_random_bytes(6), 'hex'), 1, 12),
      NEW.tier,
      'Default API Key'
    );
  END IF;
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS trigger
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;