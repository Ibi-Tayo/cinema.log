-- +goose Up
-- +goose StatementBegin

-- Create Films table
CREATE TABLE films (
    film_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    genre TEXT,
    director TEXT,
    poster_url TEXT,
    release_year INTEGER
);

-- Add primary key constraint
ALTER TABLE films 
ADD CONSTRAINT pk_films PRIMARY KEY (film_id);

-- Create Reviews table
CREATE TABLE reviews (
    review_id UUID NOT NULL,
    content VARCHAR(3000),
    date TIMESTAMP(6) NOT NULL,
    rating REAL NOT NULL,
    film_id UUID NOT NULL,
    user_id UUID NOT NULL
);

-- Add primary key constraint
ALTER TABLE reviews 
ADD CONSTRAINT pk_reviews PRIMARY KEY (review_id);

-- Create indexes
CREATE INDEX ix_reviews_film_id ON reviews (film_id);
CREATE INDEX ix_reviews_user_id ON reviews (user_id);

-- Add foreign key constraints
ALTER TABLE reviews 
ADD CONSTRAINT fk_reviews_films_film_id 
FOREIGN KEY (film_id) REFERENCES films (film_id) ON DELETE CASCADE;

ALTER TABLE reviews 
ADD CONSTRAINT fk_reviews_users_user_id 
FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

-- Create UserFilmRatings table
CREATE TABLE user_film_ratings (
    user_film_rating_id UUID NOT NULL,
    user_id UUID NOT NULL,
    film_id UUID NOT NULL,
    elo_rating DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    number_of_comparisons INTEGER NOT NULL,
    last_updated TIMESTAMP(6) NOT NULL,
    initial_rating REAL NOT NULL,
    k_constant_value DOUBLE PRECISION NOT NULL DEFAULT 0.0
);

-- Add primary key constraint
ALTER TABLE user_film_ratings 
ADD CONSTRAINT pk_user_film_ratings PRIMARY KEY (user_film_rating_id);

-- Create indexes
CREATE INDEX ix_user_film_ratings_film_id ON user_film_ratings (film_id);
CREATE INDEX ix_user_film_ratings_user_id ON user_film_ratings (user_id);

-- Add foreign key constraints
ALTER TABLE user_film_ratings 
ADD CONSTRAINT fk_user_film_ratings_films_film_id 
FOREIGN KEY (film_id) REFERENCES films (film_id) ON DELETE CASCADE;

ALTER TABLE user_film_ratings 
ADD CONSTRAINT fk_user_film_ratings_users_user_id 
FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

-- Create ComparisonHistories table
CREATE TABLE comparison_histories (
    comparison_history_id UUID NOT NULL,
    user_id UUID NOT NULL,
    film_a_film_id UUID NOT NULL,
    film_b_film_id UUID NOT NULL,
    winning_film_film_id UUID,
    comparison_date TIMESTAMP(6) NOT NULL,
    was_equal BOOLEAN NOT NULL
);

-- Add primary key constraint
ALTER TABLE comparison_histories 
ADD CONSTRAINT pk_comparison_histories PRIMARY KEY (comparison_history_id);

-- Create indexes
CREATE INDEX ix_comparison_histories_film_a_film_id ON comparison_histories (film_a_film_id);
CREATE INDEX ix_comparison_histories_film_b_film_id ON comparison_histories (film_b_film_id);
CREATE INDEX ix_comparison_histories_user_id ON comparison_histories (user_id);
CREATE INDEX ix_comparison_histories_winning_film_film_id ON comparison_histories (winning_film_film_id);

-- Add foreign key constraints
ALTER TABLE comparison_histories 
ADD CONSTRAINT fk_comparison_histories_films_film_a_film_id 
FOREIGN KEY (film_a_film_id) REFERENCES films (film_id);

ALTER TABLE comparison_histories 
ADD CONSTRAINT fk_comparison_histories_films_film_b_film_id 
FOREIGN KEY (film_b_film_id) REFERENCES films (film_id);

ALTER TABLE comparison_histories 
ADD CONSTRAINT fk_comparison_histories_films_winning_film_film_id 
FOREIGN KEY (winning_film_film_id) REFERENCES films (film_id);

ALTER TABLE comparison_histories 
ADD CONSTRAINT fk_comparison_histories_users_user_id 
FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS films CASCADE;
DROP TABLE IF EXISTS user_film_ratings CASCADE;
DROP TABLE IF EXISTS comparison_histories CASCADE;
-- +goose StatementEnd
