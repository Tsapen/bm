CREATE TABLE IF NOT EXISTS books (
    id SERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author VARCHAR(100) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    published_date timestamp,
    edition VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS collections (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS books_collection (
    collection_id INT NOT NULL REFERENCES collections(id),
    book_id INT NOT NULL REFERENCES books(id),
    PRIMARY KEY(book_id, collection_id)
);