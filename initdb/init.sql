CREATE TABLE categories (
  id serial not null PRIMARY key,
  name VARCHAR(100) not null
);

create table expenses (
  id serial not null PRIMARY KEY,
  date DATE not NULL,
  description VARCHAR(255),
  amount BIGINT not null,
  created_at TIMESTAMP DEFAULT NOW(),
  user_id VARCHAR(36) NOT NULL,
  category_id int not NULL,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);

INSERT INTO categories (name)
VALUES ('food'),
  ('drink'),
  ('sweets'),
  ('home');


INSERT INTO expenses (date, description, amount, user_id, category_id)
VALUES ('2024-01-17', 'Handmade', 8701, 'd3aaa3f0-b9aa-4011-9549-c6d7c406c336', 4),
	('2024-02-09', 'Licensed', 2316, 'ac4fbc41-86ee-4c44-bbbf-9a4012055090', 3),
    ('2024-01-15', 'Fantastic', 4539, '28dbcb06-64e8-4edc-beb2-0f3d14cf2b5f', 2),
    ('2024-01-03', 'Awesome', 3432, '4feef6e1-3145-46fb-81bf-818137f383ba', 3),
    ('2024-02-21', 'Bespoke', 9730, '988d8221-c36e-402c-bb6a-5fd5837e7fd3', 4),
    ('2024-02-24', 'Generic', 1774, '040d5fcb-2878-42e8-8205-29add0367a12', 4),
    ('2024-02-24', 'Generic', 71020, '040d5fcb-2878-42e8-8205-29add0367a12', 3),
    ('2024-02-24', 'Generic', 68407, '040d5fcb-2878-42e8-8205-29add0367a12', 2)
    
   