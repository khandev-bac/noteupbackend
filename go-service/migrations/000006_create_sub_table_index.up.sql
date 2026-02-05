CREATE INDEX subscriptions_user_id_idx
ON subscriptions (user_id);

-- Status filtering
CREATE INDEX subscriptions_status_idx
ON subscriptions (status);

-- Plan filtering
CREATE INDEX subscriptions_plan_idx
ON subscriptions (plan);

-- Active subscriptions lookup
CREATE INDEX subscriptions_active_idx
ON subscriptions (user_id)
WHERE status = 'active' AND is_deleted = FALSE;
