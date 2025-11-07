-- Simple verification script to check current database state
-- This script will check what tables exist and their names

-- Check for tables with either name
SELECT
    table_name,
    table_type
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_name IN ('anomalies', 'anomaly_events')
ORDER BY table_name;

-- Check for required tables
SELECT
    table_name,
    table_type
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_name IN ('organizations', 'usage_counters', 'subscriptions', 'billing_customers')
ORDER BY table_name;

-- Check if increment_usage function exists
SELECT routine_name, routine_type
FROM information_schema.routines
WHERE routine_schema = 'public'
  AND routine_name = 'increment_usage';

-- Get basic counts (if tables exist)
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '=== TABLE COUNTS ===';

    -- Check anomalies vs anomaly_events
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'anomalies') THEN
        EXECUTE 'SELECT COUNT(*) as anomalies_count FROM anomalies' INTO result;
        RAISE NOTICE 'anomalies table exists with % records', result;
    ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'anomaly_events') THEN
        EXECUTE 'SELECT COUNT(*) as anomaly_events_count FROM anomaly_events' INTO result;
        RAISE NOTICE 'anomaly_events table exists with % records', result;
    ELSE
        RAISE NOTICE 'Neither anomalies nor anomaly_events table found!';
    END IF;

    -- Check other tables
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'organizations') THEN
        EXECUTE 'SELECT COUNT(*) as org_count FROM organizations' INTO result;
        RAISE NOTICE 'organizations table exists with % records', result;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'usage_counters') THEN
        EXECUTE 'SELECT COUNT(*) as usage_count FROM usage_counters' INTO result;
        RAISE NOTICE 'usage_counters table exists with % records', result;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'Error checking table counts: %', SQLERRM;
END $$;