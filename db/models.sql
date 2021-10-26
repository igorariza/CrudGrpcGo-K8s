CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL, 
    password text NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp DEFAULT now(),
    CONSTRAINT pk_notes PRIMARY KEY(id)
);

