CREATE TABLE IF NOT EXISTS books (
    id SERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author VARCHAR(100) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    published_date TIMESTAMP,
    edition VARCHAR(100) NOT NULL,
    description TEXT,

    CONSTRAINT unique_book_author_title_edition UNIQUE (author, title, edition)
);

CREATE INDEX IF NOT EXISTS book_title ON books (title);

CREATE INDEX IF NOT EXISTS book_author ON books (author);

CREATE INDEX IF NOT EXISTS book_genre ON books (genre);

CREATE INDEX IF NOT EXISTS idx_published_date ON books (published_date) WHERE published_date IS NOT NULL;

CREATE TABLE IF NOT EXISTS collections (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS books_collection (
    collection_id INT NOT NULL REFERENCES collections(id),
    book_id INT NOT NULL REFERENCES books(id),
    PRIMARY KEY(book_id, collection_id)
);

CREATE INDEX IF NOT EXISTS idx_collection_book ON books_collection (collection_id, book_id);
