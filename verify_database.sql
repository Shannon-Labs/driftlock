-- Final verification queries for DriftLock database

-- Check tables
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
  AND table_name IN ('anomalies', 'usage_counters')
ORDER BY table_name;

-- Check function
SELECT routine_name, routine_type
FROM information_schema.routines 
WHERE routine_schema = 'public' 
  AND routine_name = 'increment_usage';

-- Check indexes on usage_counters
SELECT indexname, tablename 
FROM pg_indexes 
WHERE tablename = 'usage_counters';

-- Check usage_counters table structure
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_schema = 'public' 
  AND table_name = 'usage_counters'
ORDER BY ordinal_position;
