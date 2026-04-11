CREATE TABLE IF NOT EXISTS snippet_tags(
  snippet_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  PRIMARY KEY (snippet_id, tag_id),
  FOREIGN KEY (snippet_id) REFERENCES snippets(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
