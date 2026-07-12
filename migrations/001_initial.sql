CREATE TABLE decks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL
);

CREATE TABLE cards (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  deck_id INTEGER NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
  front TEXT NOT NULL,
  back TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  times_reviewed INTEGER NOT NULL DEFAULT 0,
  last_grade TEXT CHECK (last_grade IS NULL OR last_grade IN ('again', 'hard', 'easy')),
  last_reviewed_at TEXT,
  sessions_since_review INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE queue (
  deck_id INTEGER NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
  position INTEGER NOT NULL,
  card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
  PRIMARY KEY (deck_id, position),
  UNIQUE (deck_id, card_id)
);

CREATE INDEX idx_cards_deck_id ON cards(deck_id);
CREATE INDEX idx_queue_card_id ON queue(card_id);
