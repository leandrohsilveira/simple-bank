-- Active: 1721854015531@@localhost@5432@postgres
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'app_test' -- ← change this to your DB
  AND pid <> pg_backend_pid();

drop database app_test;