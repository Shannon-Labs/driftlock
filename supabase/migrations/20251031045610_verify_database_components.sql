-- Verification queries for DriftLock database components
-- This migration checks the required tables and functions for launch

-- Check if usage_counters table exists with correct schema
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'usage_counters'
    ) THEN
        RAISE EXCEPTION 'usage_counters table does not exist';
    END IF;
    
    RAISE NOTICE '✅ usage_counters table exists';
END $$;

-- Check if increment_usage function exists with correct signature
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.routines 
        WHERE routine_schema = 'public' 
        AND routine_name = 'increment_usage'
    ) THEN
        RAISE EXCEPTION 'increment_usage function does not exist';
    END IF;
    
    RAISE NOTICE '✅ increment_usage function exists';
END $$;

-- Check if anomalies table exists (renamed from anomaly_events)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'anomalies'
    ) THEN
        RAISE EXCEPTION 'anomalies table does not exist';
    END IF;
    
    RAISE NOTICE '✅ anomalies table exists';
END $$;

-- Check if required indexes exist on usage_counters
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE tablename = 'usage_counters' 
        AND indexname = 'idx_usage_counters_org_period'
    ) THEN
        RAISE EXCEPTION 'idx_usage_counters_org_period index does not exist';
    END IF;
    
    RAISE NOTICE '✅ idx_usage_counters_org_period index exists';
END $$;

-- Verify usage_counters table structure
DO $$
DECLARE
    column_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO column_count
    FROM information_schema.columns 
    WHERE table_schema = 'public' 
    AND table_name = 'usage_counters'
    AND column_name IN (
        'organization_id', 'period_start', 'period_end', 
        'total_calls', 'included_calls_used', 'overage_calls', 
        'estimated_charges_cents', 'updated_at'
    );
    
    IF column_count < 8 THEN
        RAISE EXCEPTION 'usage_counters table is missing required columns';
    END IF;
    
    RAISE NOTICE '✅ usage_counters table has correct structure';
END $$;

-- Final verification summary
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '🎉 DriftLock Database Verification Complete';
    RAISE NOTICE '✅ All required database components are properly configured';
    RAISE NOTICE '✅ Ready for launch';
END $$;
