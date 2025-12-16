-- トリガーの削除
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

-- 関数の削除
DROP FUNCTION IF EXISTS update_updated_at_column();

-- インデックスの削除
DROP INDEX IF EXISTS idx_events_end_date;
DROP INDEX IF EXISTS idx_events_start_date;

-- テーブルの削除
DROP TABLE IF EXISTS events;
