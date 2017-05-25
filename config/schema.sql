CREATE TABLE IF NOT EXISTS scores (
  id      serial      PRIMARY KEY,
  game    varchar(32) NOT NULL,
  player  varchar(32) NOT NULL,
  score   integer     NOT NULL DEFAULT 0,
  created timestamptz NOT NULL DEFAULT now(),
  UNIQUE (game, player)
);
CREATE INDEX player_i ON scores (player);
CREATE INDEX score_i ON scores (score);
CREATE INDEX created_i ON scores (created);