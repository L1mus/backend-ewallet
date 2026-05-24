DROP TYPE IF EXISTS oauth_provider AS ENUM (
  'google',
  'facebook'
);

DROP TYPE IF EXISTS subscribe_status AS ENUM (
  'active',
  'unsubscribe'
);

DROP TYPE IF EXISTS type_transaction AS ENUM (
  'income',
  'expense'
);

DROP TYPE IF EXISTS type_activity_transaction AS ENUM (
  'transfer',
  'topup'
);

DROP TYPE IF EXISTS status_transaction AS ENUM (
  'pending',
  'success',
  'failed'
);

DROP TABLE users;