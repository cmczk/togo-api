CREATE TABLE IF NOT EXISTS todos(
      id SERIAL PRIMARY KEY
    , title VARCHAR(250) NOT NULL
    , description TEXT
    , completed BOOLEAN DEFAULT FALSE
    , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    todos (title, description, completed)
VALUES
    ('Задача №1', 'Описание №1', FALSE),
    ('Задача №2', 'Описание №2', FALSE),
    ('Задача №3', 'Описание №3', FALSE)
