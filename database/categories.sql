INSERT INTO categories (name) VALUES
('General'),
('Announcements'),
('Tech'),
('Music'),
('Games'),
('Humour'),
('Personal')
ON CONFLICT (name) DO NOTHING;